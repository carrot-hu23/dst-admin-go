package levelConfigUtils

import (
	"dst-admin-go/constant/dst"
	"dst-admin-go/utils/fileUtils"
	"encoding/json"
	"fmt"
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

func GetLevelConfig(clusterName string) (*LevelConfig, error) {
	clusterBasePath := dst.GetClusterBasePath(clusterName)
	jsonPath := filepath.Join(clusterBasePath, "level.json")
	// fileUtils.CreateFileIfNotExists(jsonPath)
	if !fileUtils.Exists(jsonPath) {
		fileUtils.CreateFile(jsonPath)
		fileUtils.WriterTXT(jsonPath, "{}")
	}
	// 打开JSON文件
	file, err := os.Open(jsonPath)
	if err != nil {
		fmt.Println("无法打开level.json文件:", err)
		return nil, err
	}
	defer file.Close()

	// 解码JSON数据
	var config LevelConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("无法解析level.json文件:", err)
		return nil, err
	}
	return &config, nil
}

func SaveLevelConfig(clusterName string, levelConfig *LevelConfig) error {
	clusterBasePath := dst.GetClusterBasePath(clusterName)
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
