package service

import (
	"dst-admin-go/utils/dstUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/utils/levelConfigUtils"
	"dst-admin-go/vo/level"
	"log"
	"path/filepath"
)

var gameServe GameService
var homeServe HomeService

type GameLevel2Service struct{}

func (g *GameLevel2Service) GetLevelList(clusterName string) []level.World {
	config, err := levelConfigUtils.GetLevelConfig(clusterName)
	if err != err {
		return []level.World{}
	}
	var levels []level.World
	if len(config.LevelList) == 0 {
		masterLevelPath := filepath.Join(dstUtils.GetClusterBasePath(clusterName), "Master")
		if !fileUtils.Exists(masterLevelPath) {
			master := level.World{
				IsMaster:          true,
				LevelName:         "森林",
				Uuid:              "Master",
				Leveldataoverride: "return {}",
				Modoverrides:      "return {}",
				ServerIni:         level.NewMasterServerIni(),
			}
			initLevel(filepath.Join(dstUtils.GetClusterBasePath(clusterName), "Master"), &master)
			levels = append([]level.World{}, master)
			config.LevelList = append(config.LevelList, levelConfigUtils.Item{
				Name: "森林",
				File: "Master",
			})
			err = levelConfigUtils.SaveLevelConfig(clusterName, config)
			if err != nil {
				log.Println(err)
			}
			return levels
		} else {
			masterLevel := homeServe.GetLevel(clusterName, "Master")
			master := level.World{
				IsMaster:          true,
				LevelName:         "森林",
				Uuid:              "Master",
				Leveldataoverride: masterLevel.Leveldataoverride,
				Modoverrides:      masterLevel.Modoverrides,
				ServerIni:         masterLevel.ServerIni,
			}
			levels = append([]level.World{}, master)
			config.LevelList = append(config.LevelList, levelConfigUtils.Item{
				Name: "森林",
				File: "Master",
			})
			cavesLevelPath := filepath.Join(dstUtils.GetClusterBasePath(clusterName), "Caves")
			if fileUtils.Exists(cavesLevelPath) {
				config.LevelList = append(config.LevelList, levelConfigUtils.Item{
					Name: "洞穴",
					File: "Caves",
				})
				cavesLevel := homeServe.GetLevel(clusterName, "Caves")
				caves := level.World{
					IsMaster:          true,
					LevelName:         "洞穴",
					Uuid:              "Caves",
					Leveldataoverride: cavesLevel.Leveldataoverride,
					Modoverrides:      cavesLevel.Modoverrides,
					ServerIni:         cavesLevel.ServerIni,
				}
				levels = append(levels, caves)
			}
			err = levelConfigUtils.SaveLevelConfig(clusterName, config)
			if err != nil {
				log.Println(err)
			}
			return levels
		}
	}
	for i := range config.LevelList {
		level1 := level.World{}
		level1.LevelName = config.LevelList[i].Name
		level1.Uuid = config.LevelList[i].File
		world := homeServe.GetLevel(clusterName, config.LevelList[i].File)
		level1.Leveldataoverride = world.Leveldataoverride
		level1.Modoverrides = world.Modoverrides
		level1.ServerIni = world.ServerIni
		levels = append(levels, level1)
	}
	return levels
}

func (g *GameLevel2Service) UpdateLevels(clusterName string, levels []level.World) error {
	// TODO 保留之前的数据

	for i := range levels {
		dstUtils.DedicatedServerModsSetup2(clusterName, levels[i].Modoverrides)
		err := g.UpdateLevel(clusterName, &levels[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *GameLevel2Service) UpdateLevel(clusterName string, level *level.World) error {
	// cluster := clusterUtils.GetCluster(clusterName)
	levelFolderPath := filepath.Join(dstUtils.GetClusterBasePath(clusterName), level.Uuid)
	fileUtils.CreateDirIfNotExists(levelFolderPath)
	initLevel(levelFolderPath, level)

	// 记录level.json 文件
	levelConfig, err := levelConfigUtils.GetLevelConfig(clusterName)
	if err != nil {
		return err
	}
	for i := range levelConfig.LevelList {
		if level.Uuid == levelConfig.LevelList[i].File {
			levelConfig.LevelList[i].Name = level.LevelName
			// 更新世界配置
			initLevel(levelFolderPath, level)
			break
		}
	}
	err = levelConfigUtils.SaveLevelConfig(clusterName, levelConfig)
	return err
}

func (g *GameLevel2Service) CreateLevel(clusterName string, level *level.World) error {
	uuid := ""
	if level.Uuid == "" {
		uuid = generateUUID()
	} else {
		uuid = level.Uuid
	}
	// cluster := clusterUtils.GetCluster(clusterName)
	levelFolderPath := filepath.Join(dstUtils.GetClusterBasePath(clusterName), uuid)
	fileUtils.CreateDirIfNotExists(levelFolderPath)
	initLevel(levelFolderPath, level)

	// 记录level.json 文件
	levelConfig, err := levelConfigUtils.GetLevelConfig(clusterName)
	if err != nil {
		return err
	}
	levelConfig.LevelList = append(levelConfig.LevelList, levelConfigUtils.Item{Name: level.LevelName, File: uuid})
	err = levelConfigUtils.SaveLevelConfig(clusterName, levelConfig)
	if err != nil {
		err := fileUtils.DeleteFile(filepath.Join(dstUtils.GetClusterBasePath(clusterName), uuid))
		return err
	}
	level.Uuid = uuid
	return nil
}

func (g *GameLevel2Service) DeleteLevel(clusterName string, levelName string) error {
	gameServe.shutdownLevel(clusterName, levelName)
	if clusterName == "" {
		log.Panicln("cluster is not allow null")
	}
	err := fileUtils.DeleteDir(filepath.Join(dstUtils.GetClusterBasePath(clusterName), levelName))
	if err != nil {
		return err
	}
	// 删除 json 文件
	config, err := levelConfigUtils.GetLevelConfig(clusterName)
	if err != nil {
		log.Panicln("删除文件失败")
	}
	newLevelsConfig := levelConfigUtils.LevelConfig{}
	for i := range config.LevelList {
		if config.LevelList[i].File != levelName {
			newLevelsConfig.LevelList = append(newLevelsConfig.LevelList, config.LevelList[i])
		}
	}
	err = levelConfigUtils.SaveLevelConfig(clusterName, &newLevelsConfig)

	// TODO 同时删除定时任务和自动维护

	return err
}

func initLevel(levelFolderPath string, level *level.World) {

	lPath := filepath.Join(levelFolderPath, "leveldataoverride.lua")
	mPath := filepath.Join(levelFolderPath, "modoverrides.lua")
	sPath := filepath.Join(levelFolderPath, "server.ini")

	fileUtils.CreateFileIfNotExists(lPath)
	fileUtils.CreateFileIfNotExists(mPath)
	fileUtils.CreateFileIfNotExists(sPath)

	fileUtils.WriterTXT(lPath, level.Leveldataoverride)
	fileUtils.WriterTXT(mPath, level.Modoverrides)
	serverBuf := dstUtils.ParseTemplate(ServerIniTemplate, level.ServerIni)
	fileUtils.WriterTXT(sPath, serverBuf)
}
