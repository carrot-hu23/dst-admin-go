package level

import (
	"dst-admin-go/internal/pkg/utils/dstUtils"
	"dst-admin-go/internal/pkg/utils/fileUtils"
	"dst-admin-go/internal/service/archive"
	"dst-admin-go/internal/service/dstConfig"
	"dst-admin-go/internal/service/game"
	"dst-admin-go/internal/service/gameConfig"
	"dst-admin-go/internal/service/levelConfig"
	"log"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-ini/ini"
)

// LevelService 关卡服务结构体
type LevelService struct {
	gameProcess      game.Process
	dstConfig        dstConfig.Config
	resolver         *archive.PathResolver
	levelConfigUtils *levelConfig.LevelConfigUtils
}

// NewLevelService 创建关卡服务实例
func NewLevelService(gameProcess game.Process, dstConfig dstConfig.Config, resolver *archive.PathResolver, levelConfigUtils *levelConfig.LevelConfigUtils) *LevelService {
	return &LevelService{
		gameProcess:      gameProcess,
		dstConfig:        dstConfig,
		resolver:         resolver,
		levelConfigUtils: levelConfigUtils,
	}
}

// GetLevelList 获取关卡列表
func (l *LevelService) GetLevelList(clusterName string) []levelConfig.LevelInfo {
	config, err := l.levelConfigUtils.GetLevelConfig(clusterName)
	if err != nil {
		return []levelConfig.LevelInfo{}
	}
	var levels []levelConfig.LevelInfo
	if len(config.LevelList) == 0 {
		masterLevelPath := filepath.Join(l.resolver.ClusterPath(clusterName), "Master")
		if !fileUtils.Exists(masterLevelPath) {
			master := levelConfig.LevelInfo{
				IsMaster:          true,
				LevelName:         "森林",
				Uuid:              "Master",
				Leveldataoverride: "return {}",
				Modoverrides:      "return {}",
				ServerIni:         levelConfig.NewMasterServerIni(),
			}
			l.initLevel(filepath.Join(l.resolver.ClusterPath(clusterName), "Master"), &master)
			levels = append([]levelConfig.LevelInfo{}, master)
			config.LevelList = append(config.LevelList, levelConfig.Item{
				Name: "森林",
				File: "Master",
			})
			err = l.levelConfigUtils.SaveLevelConfig(clusterName, config)
			if err != nil {
				log.Println(err)
			}
			return levels
		} else {
			// 读取现有的 Master 世界配置
			master := l.GetLevel(clusterName, "Master")
			levels = append([]levelConfig.LevelInfo{}, master)
			config.LevelList = append(config.LevelList, levelConfig.Item{
				Name: "森林",
				File: "Master",
			})
			cavesLevelPath := filepath.Join(l.resolver.ClusterPath(clusterName), "Caves")
			if fileUtils.Exists(cavesLevelPath) {
				config.LevelList = append(config.LevelList, levelConfig.Item{
					Name: "洞穴",
					File: "Caves",
				})
				caves := l.GetLevel(clusterName, "Caves")
				levels = append(levels, caves)
			}
			err = l.levelConfigUtils.SaveLevelConfig(clusterName, config)
			if err != nil {
				log.Println(err)
			}
			return levels
		}
	}
	for i := range config.LevelList {
		level1 := levelConfig.LevelInfo{}
		level1.LevelName = config.LevelList[i].Name
		level1.Uuid = config.LevelList[i].File
		level1.RunVersion = config.LevelList[i].RunVersion
		world := l.GetLevel(clusterName, config.LevelList[i].File)
		level1.Leveldataoverride = world.Leveldataoverride
		level1.Modoverrides = world.Modoverrides
		level1.ServerIni = world.ServerIni
		levels = append(levels, level1)
	}
	return levels
}

// GetLevel 获取单个关卡配置
func (l *LevelService) GetLevel(clusterName string, levelName string) levelConfig.LevelInfo {
	levelFolderPath := filepath.Join(l.resolver.ClusterPath(clusterName), levelName)
	config, _ := l.levelConfigUtils.GetLevelConfig(clusterName)
	name := ""
	for _, item := range config.LevelList {
		if item.File == levelName {
			name = item.Name
		}
	}

	// 读取 leveldataoverride.lua
	lPath := filepath.Join(levelFolderPath, "leveldataoverride.lua")
	leveldataoverride, err := fileUtils.ReadFile(lPath)
	if err != nil {
		leveldataoverride = "return {}"
	}

	// 读取 modoverrides.lua
	mPath := filepath.Join(levelFolderPath, "modoverrides.lua")
	modoverrides, err := fileUtils.ReadFile(mPath)
	if err != nil {
		modoverrides = "return {}"
	}

	// 读取 server.ini
	sPath := filepath.Join(levelFolderPath, "server.ini")
	serverIni := l.GetServerIni(sPath, levelName == "Master")

	return levelConfig.LevelInfo{
		IsMaster:          levelName == "Master",
		LevelName:         name,
		Uuid:              levelName,
		Leveldataoverride: leveldataoverride,
		Modoverrides:      modoverrides,
		ServerIni:         serverIni,
	}
}

func (l *LevelService) GetServerIni(filepath string, isMaster bool) levelConfig.ServerIni {
	fileUtils.CreateFileIfNotExists(filepath)
	var serverPortDefault uint = 10998
	idDefault := 10010

	if isMaster {
		serverPortDefault = 10999
		idDefault = 10000
	}

	serverIni := levelConfig.NewCavesServerIni()
	// 加载 INI 文件
	cfg, err := ini.Load(filepath)
	if err != nil {
		return serverIni
	}

	// [NETWORK]
	NETWORK := cfg.Section("NETWORK")

	serverIni.ServerPort = NETWORK.Key("server_port").MustUint(serverPortDefault)

	// [SHARD]
	SHARD := cfg.Section("SHARD")

	serverIni.IsMaster = SHARD.Key("is_master").MustBool(isMaster)
	serverIni.Name = SHARD.Key("name").String()
	serverIni.Id = SHARD.Key("id").MustUint(uint(idDefault))

	// [ACCOUNT]
	ACCOUNT := cfg.Section("ACCOUNT")
	serverIni.EncodeUserPath = ACCOUNT.Key("encode_user_path").MustBool(true)

	// [STEAM]
	STEAM := cfg.Section("STEAM")

	serverIni.AuthenticationPort = STEAM.Key("authentication_port").MustUint(8766)
	serverIni.MasterServerPort = STEAM.Key("master_server_port").MustUint(27016)

	return serverIni
}

// UpdateLevels 更新多个关卡配置
func (l *LevelService) UpdateLevels(clusterName string, levels []levelConfig.LevelInfo) error {
	for i := range levels {
		config, _ := l.dstConfig.GetDstConfig(clusterName)
		dstUtils.DedicatedServerModsSetup(config, levels[i].Modoverrides)
		err := l.UpdateLevel(clusterName, &levels[i])
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateLevel 更新单个关卡配置
func (l *LevelService) UpdateLevel(clusterName string, level *levelConfig.LevelInfo) error {
	levelFolderPath := filepath.Join(l.resolver.ClusterPath(clusterName), level.Uuid)
	fileUtils.CreateDirIfNotExists(levelFolderPath)
	l.initLevel(levelFolderPath, level)

	// 记录level.json 文件
	levelConfig, err := l.levelConfigUtils.GetLevelConfig(clusterName)
	if err != nil {
		return err
	}
	for i := range levelConfig.LevelList {
		if level.Uuid == levelConfig.LevelList[i].File {
			levelConfig.LevelList[i].Name = level.LevelName
			// 更新世界配置
			l.initLevel(levelFolderPath, level)
			break
		}
	}
	err = l.levelConfigUtils.SaveLevelConfig(clusterName, levelConfig)
	return err
}

// CreateLevel 创建新关卡
func (l *LevelService) CreateLevel(clusterName string, level *levelConfig.LevelInfo) error {
	uuid := ""
	if level.Uuid == "" {
		uuid = l.generateUUID()
	} else {
		uuid = level.Uuid
	}
	levelFolderPath := filepath.Join(l.resolver.ClusterPath(clusterName), uuid)
	fileUtils.CreateDirIfNotExists(levelFolderPath)
	l.initLevel(levelFolderPath, level)

	// 记录level.json 文件
	config, err := l.levelConfigUtils.GetLevelConfig(clusterName)
	if err != nil {
		return err
	}
	config.LevelList = append(config.LevelList, levelConfig.Item{Name: level.LevelName, File: uuid})
	err = l.levelConfigUtils.SaveLevelConfig(clusterName, config)
	if err != nil {
		err := fileUtils.DeleteFile(filepath.Join(l.resolver.ClusterPath(clusterName), uuid))
		return err
	}
	level.Uuid = uuid
	return nil
}

// DeleteLevel 删除关卡
func (l *LevelService) DeleteLevel(clusterName string, levelName string) error {
	// 停止关卡服务
	l.gameProcess.Stop(clusterName, levelName)

	// 删除关卡目录
	err := fileUtils.DeleteDir(filepath.Join(l.resolver.ClusterPath(clusterName), levelName))
	if err != nil {
		return err
	}

	// 删除 json 文件中的记录
	config, err := l.levelConfigUtils.GetLevelConfig(clusterName)
	if err != nil {
		log.Panicln("删除文件失败")
	}
	newLevelsConfig := levelConfig.LevelConfig{}
	for i := range config.LevelList {
		if config.LevelList[i].File != levelName {
			newLevelsConfig.LevelList = append(newLevelsConfig.LevelList, config.LevelList[i])
		}
	}
	err = l.levelConfigUtils.SaveLevelConfig(clusterName, &newLevelsConfig)

	// TODO 同时删除定时任务和自动维护

	return err
}

// initLevel 初始化关卡文件
func (l *LevelService) initLevel(levelFolderPath string, level *levelConfig.LevelInfo) {

	lPath := filepath.Join(levelFolderPath, "leveldataoverride.lua")
	mPath := filepath.Join(levelFolderPath, "modoverrides.lua")
	sPath := filepath.Join(levelFolderPath, "server.ini")

	fileUtils.CreateFileIfNotExists(lPath)
	fileUtils.CreateFileIfNotExists(mPath)
	fileUtils.CreateFileIfNotExists(sPath)

	fileUtils.WriterTXT(lPath, level.Leveldataoverride)
	fileUtils.WriterTXT(mPath, level.Modoverrides)
	serverBuf := l.ParseTemplate(level.ServerIni)
	fileUtils.WriterTXT(sPath, serverBuf)
}

// ParseTemplate 解析服务器配置模板
func (l *LevelService) ParseTemplate(serverIni levelConfig.ServerIni) string {
	// 使用 gameConfig 中的 ParseTemplate 方法
	return dstUtils.ParseTemplate(gameConfig.ServerIniTemplate, serverIni)
}

// generateUUID 生成 UUID
func (l *LevelService) generateUUID() string {
	// 简化实现，实际应该使用标准的 UUID 生成库
	return "level_" + strconv.FormatInt(time.Now().UnixNano(), 10)
}
