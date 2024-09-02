package levelConfigUtils

import (
	"dst-admin-go/utils/dstUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo/level"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Item struct {
	Name string `json:"name"`
	File string `json:"file"`
}

type LevelConfig struct {
	LevelList []Item `json:"levelList"`
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
	serverBuf := dstUtils.ParseTemplate("./static/template/server.ini", level.ServerIni)
	fileUtils.WriterTXT(sPath, serverBuf)
}

func GetLevelConfig(clusterName string) (*LevelConfig, error) {
	clusterBasePath := dstUtils.GetClusterBasePath(clusterName)
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
			config.LevelList = append(config.LevelList, Item{
				Name: "森林",
				File: "Master",
			})
			err = SaveLevelConfig(clusterName, &config)
			if err != nil {
				log.Println(err)
			}
		} else {
			config.LevelList = append(config.LevelList, Item{
				Name: "森林",
				File: "Master",
			})
			cavesLevelPath := filepath.Join(dstUtils.GetClusterBasePath(clusterName), "Caves")
			if fileUtils.Exists(cavesLevelPath) {
				config.LevelList = append(config.LevelList, Item{
					Name: "洞穴",
					File: "Caves",
				})
			}
			err = SaveLevelConfig(clusterName, &config)
			if err != nil {
				log.Println(err)
			}
		}
	}

	return &config, nil
}

func SaveLevelConfig(clusterName string, levelConfig *LevelConfig) error {
	clusterBasePath := dstUtils.GetClusterBasePath(clusterName)
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

	// 创建 JSON 编码器
	//encoder := json.NewEncoder(file)
	//
	//// 将结构体编码为 JSON 并写入文件
	//err = encoder.Encode(levelConfig)
	//if err != nil {
	//	fmt.Println("编码 level.json 失败:", err)
	//	return err
	//}
	//return nil
}
