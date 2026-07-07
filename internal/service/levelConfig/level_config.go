package levelConfig

import (
	"dst-admin-go/internal/pkg/utils/dstUtils"
	"dst-admin-go/internal/pkg/utils/fileUtils"
	"dst-admin-go/internal/service/archive"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Item struct {
	Name       string `json:"name"`
	File       string `json:"file"`
	RunVersion int64  `json:"runVersion"`
	Version    int64  `json:"Version"`
}

type LevelConfig struct {
	LevelList []Item `json:"levelList"`
}

// LevelInfo 世界配置结构体
type LevelInfo struct {
	IsMaster          bool      `json:"isMaster"`
	LevelName         string    `json:"levelName"`
	Uuid              string    `json:"uuid"`
	RunVersion        int64     `json:"runVersion"`
	Leveldataoverride string    `json:"leveldataoverride"`
	Modoverrides      string    `json:"modoverrides"`
	ServerIni         ServerIni `json:"server_ini"`
}

type ServerIni struct {

	// [NETWORK]
	ServerPort uint `json:"server_port"`

	// [SHARD]
	IsMaster bool   `json:"is_master"`
	Name     string `json:"name"`
	Id       uint   `json:"id"`

	// [ACCOUNT]
	EncodeUserPath bool `json:"encode_user_path"`

	// [STEAM]
	AuthenticationPort uint `json:"authentication_port"`
	MasterServerPort   uint `json:"master_server_port"`
}

func NewMasterServerIni() ServerIni {
	return ServerIni{
		ServerPort:     10999,
		IsMaster:       true,
		Name:           "Master",
		Id:             10000,
		EncodeUserPath: true,
	}
}

func NewCavesServerIni() ServerIni {
	return ServerIni{
		ServerPort:         10998,
		IsMaster:           false,
		Name:               "Caves",
		Id:                 10010,
		EncodeUserPath:     true,
		AuthenticationPort: 8766,
		MasterServerPort:   27016,
	}
}

type LevelConfigUtils struct {
	archive *archive.PathResolver
}

func NewLevelConfigUtils(archive *archive.PathResolver) *LevelConfigUtils {
	return &LevelConfigUtils{
		archive: archive,
	}
}

func (p *LevelConfigUtils) initLevel(levelFolderPath string, level *LevelInfo) {

	lPath := filepath.Join(levelFolderPath, "leveldataoverride.lua")
	mPath := filepath.Join(levelFolderPath, "modoverrides.lua")
	sPath := filepath.Join(levelFolderPath, "server.ini")

	fileUtils.CreateFileIfNotExists(lPath)
	fileUtils.CreateFileIfNotExists(mPath)
	fileUtils.CreateFileIfNotExists(sPath)

	fileUtils.WriterTXT(lPath, level.Leveldataoverride)
	fileUtils.WriterTXT(mPath, level.Modoverrides)
	serverBuf := dstUtils.ParseTemplate("./static/template/server.ini", level.ServerIni)
	fileUtils.WriterTXT(sPath, serverBuf)
}

func (p *LevelConfigUtils) GetLevelConfig(clusterName string) (*LevelConfig, error) {
	clusterBasePath := p.archive.ClusterPath(clusterName)
	jsonPath := filepath.Join(clusterBasePath, "level.json")
	fileUtils.CreateDirIfNotExists(clusterBasePath)
	// fileUtils.CreateFileIfNotExists(jsonPath)
	if !fileUtils.Exists(jsonPath) {
		fileUtils.CreateFile(jsonPath)
		fileUtils.WriterTXT(jsonPath, "{}")
	}
	// 打开JSON文件
	file, err := os.Open(jsonPath)
	if err != nil {
		log.Println("无法打开level.json文件:", err)
		return nil, err
	}
	defer file.Close()

	// 解码JSON数据
	var config LevelConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Println("无法解析level.json文件:", err)
		return nil, err
	}

	if len(config.LevelList) == 0 {
		masterLevelPath := filepath.Join(clusterBasePath, "Master")
		if !fileUtils.Exists(masterLevelPath) {
			master := LevelInfo{
				IsMaster:          true,
				LevelName:         "森林",
				Uuid:              "Master",
				Leveldataoverride: "return {}",
				Modoverrides:      "return {}",
				ServerIni:         NewMasterServerIni(),
			}
			p.initLevel(filepath.Join(clusterBasePath, "Master"), &master)
			config.LevelList = append(config.LevelList, Item{
				Name: "森林",
				File: "Master",
			})
			err = p.SaveLevelConfig(clusterName, &config)
			if err != nil {
				log.Println(err)
			}
		} else {
			config.LevelList = append(config.LevelList, Item{
				Name: "森林",
				File: "Master",
			})
			cavesLevelPath := filepath.Join(clusterBasePath, "Caves")
			if fileUtils.Exists(cavesLevelPath) {
				config.LevelList = append(config.LevelList, Item{
					Name: "洞穴",
					File: "Caves",
				})
			}
			err = p.SaveLevelConfig(clusterName, &config)
			if err != nil {
				log.Println(err)
			}
		}
	}

	return &config, nil
}

func (p *LevelConfigUtils) SaveLevelConfig(clusterName string, levelConfig *LevelConfig) error {
	clusterBasePath := p.archive.ClusterPath(clusterName)
	jsonPath := filepath.Join(clusterBasePath, "level.json")
	fileUtils.CreateFileIfNotExists(jsonPath)
	// 打开JSON文件
	file, err := os.Open(jsonPath)
	if err != nil {
		log.Println("无法打开level.json文件:", err)
		return err
	}
	defer file.Close()

	bytes, err := json.Marshal(levelConfig)
	if err != nil {
		log.Println("json 解析错误")
	}
	fileUtils.WriterTXT(jsonPath, string(bytes))
	return err
}
