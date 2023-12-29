package service

import (
	"dst-admin-go/constant/consts"
	"dst-admin-go/utils/dstUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo/level"
	"log"
	"path/filepath"
	"strings"
)

type GameLevelService struct {
	h HomeService
}

type LevelStatus struct {
	Name   string `json:"name"`
	Status bool   `json:"status"`
}

type LevelConfig struct {
	Name        string `json:"levelName"`
	LevelType   string `json:"levelType"`
	Description string `json:"description"`
}

const (
	LevelConfigName = "level_config"
	SplitFlag       = ":"
)

func (g *GameLevelService) recodeLevel(clusterName string, levelConfig LevelConfig) {
	levelConfigPath := g.getLevelConfigPath(clusterName)
	data := levelConfig.Name + SplitFlag + levelConfig.LevelType + SplitFlag + levelConfig.Description
	content, err := fileUtils.ReadLnFile(levelConfigPath)
	if err != nil {
		log.Panicln("读取 level_config 文件失败", err)
	}
	content = append(content, data)
	err = fileUtils.WriterLnFile(levelConfigPath, content)
	if err != nil {
		log.Panicln("写入 level_config 文件失败", err)
	}
}

func (g *GameLevelService) parseLevelConfig(lines []string) []LevelConfig {
	// 以:分割 Master1:master:主世界1
	var levelConfigList []LevelConfig
	for _, line := range lines {
		split := strings.Split(line, ":")
		if len(split) == 3 {
			levelConfig := LevelConfig{
				Name:        split[0],
				LevelType:   split[1],
				Description: split[2],
			}
			levelConfigList = append(levelConfigList, levelConfig)
		}
	}
	return levelConfigList
}

func (g *GameLevelService) createLevel(clusterName string, levelConfig LevelConfig, templatePath string) {
	levelPath := filepath.Join(dstUtils.GetClusterBasePath(clusterName), levelConfig.Name)
	fileUtils.CreateDirIfNotExists(levelPath)

	// 初始化 level 和 mod server.ini 文件

	leveldataoverride, err := fileUtils.ReadFile(filepath.Join(templatePath, "leveldataoverride.lua"))
	if err != nil {
		panic("read ./static/Master/leveldataoverride.lua file error: " + err.Error())
	}
	modoverrides, err := fileUtils.ReadFile(filepath.Join(templatePath, "modoverrides.lua"))
	if err != nil {
		panic("read ./static/Master/modoverrides.lua file error: " + err.Error())
	}
	server_ini, err := fileUtils.ReadFile(filepath.Join(templatePath, "server.ini"))
	if err != nil {
		panic("read /static/Master/server.ini file error: " + err.Error())
	}

	l_path := filepath.Join(levelPath, "leveldataoverride.lua")
	m_path := filepath.Join(levelPath, "modoverrides.lua")
	s_path := filepath.Join(levelPath, "server.ini")

	fileUtils.CreateFileIfNotExists(l_path)
	fileUtils.CreateFileIfNotExists(m_path)
	fileUtils.CreateFileIfNotExists(s_path)

	fileUtils.WriterTXT(l_path, leveldataoverride)
	fileUtils.WriterTXT(m_path, modoverrides)
	fileUtils.WriterTXT(s_path, server_ini)

	// TODO 写入到 level 文件里面
	g.recodeLevel(clusterName, levelConfig)
}

func (g *GameLevelService) GetLevelsConfig(clusterName string) []LevelConfig {
	levelConfigPath := g.getLevelConfigPath(clusterName)
	log.Println(filepath.Join(clusterName, LevelConfigName))
	err := fileUtils.CreateFileIfNotExists(levelConfigPath)
	if err != nil {
		log.Panicln("创建 level_config 失败", err)
	}
	content, err := fileUtils.ReadLnFile(levelConfigPath)
	if err != nil {
		log.Panicln("读取 level_config 失败", err)
	}
	return g.parseLevelConfig(content)
}

func (g *GameLevelService) CreateNewLevel(clusterName string, levelConfig LevelConfig) {

	if levelConfig.LevelType == consts.MasterLevelType {
		g.createLevel(clusterName, levelConfig, "./static/Master")
	}

	if levelConfig.LevelType == consts.CaveLevelType {
		g.createLevel(clusterName, levelConfig, "./static/Caves")
	}
}

func (g *GameLevelService) DeleteLevelConfig(clusterName, levelName string) {
	levelConfigPath := g.getLevelConfigPath(clusterName)
	content, err := fileUtils.ReadLnFile(levelConfigPath)
	if err != nil {
		log.Panicln("读取 level_config 文件失败", err)
	}
	var newContent []string
	for _, line := range content {
		split := strings.Split(line, SplitFlag)
		if len(split) == 3 {
			if split[0] != levelName {
				newContent = append(newContent, line)
			}
		}
	}
	err = fileUtils.WriterLnFile(levelConfigPath, newContent)
	if err != nil {
		log.Panicln("写入 level_config 文件失败", err)
	}

	levelPath := filepath.Join(dstUtils.GetClusterBasePath(clusterName), levelName)
	err = fileUtils.DeleteDir(levelPath)
	if err != nil {
		log.Panicln("删除level失败", levelName, err)
	}
}

func (g *GameLevelService) getLevelConfigPath(clusterName string) string {
	levelConfigPath := filepath.Join(dstUtils.GetClusterBasePath(clusterName), LevelConfigName)
	return levelConfigPath
}

func (g *GameLevelService) GetLeveldataoverride(clusterName, levelName string) string {
	leveldataoverridePath := dstUtils.GetLevelLeveldataoverridePath(clusterName, levelName)
	log.Println("leveldataoverridePath: ", leveldataoverridePath)
	if !fileUtils.Exists(leveldataoverridePath) {
		return "return {}"
	}
	leveldataoverride, err := fileUtils.ReadFile(leveldataoverridePath)
	if err != nil {
		return "return {}"
	}
	return leveldataoverride
}

func (g *GameLevelService) GetModoverrides(clusterName, levelName string) string {
	modoverridesPath := dstUtils.GetLevelModoverridesPath(clusterName, levelName)
	log.Println("modoverridesPath: ", modoverridesPath)
	if !fileUtils.Exists(modoverridesPath) {
		return "return {}"
	}
	modoverrides, err := fileUtils.ReadFile(modoverridesPath)
	if err != nil {
		return "return {}"
	}
	return modoverrides
}

func (g *GameLevelService) GetServerIni(clusterName, levelName string) *level.ServerIni {
	levelServerIniPath := dstUtils.GetLevelServerIniPath(clusterName, levelName)
	log.Println("levelServerIniPath: ", levelServerIniPath)
	if !fileUtils.Exists(levelServerIniPath) {
		return &level.ServerIni{}
	}

	return g.h.GetServerIni(levelServerIniPath, false)
}

func (g *GameLevelService) SaveLeveldataoverride(clusterName, levelName string, leveldataoverride string) {
	leveldataoverridePath := dstUtils.GetLevelLeveldataoverridePath(clusterName, levelName)
	err := fileUtils.CreateFileIfNotExists(leveldataoverridePath)
	if err != nil {
		log.Panicln("创建失败 leveldataoverride ", err)
	}
	err = fileUtils.WriterTXT(leveldataoverridePath, leveldataoverride)
	if err != nil {
		log.Panicln("写入失败 leveldataoverride", err)
	}
}

func (g *GameLevelService) SaveModoverrides(clusterName, levelName string, modoverrides string) {
	modoverridesPath := dstUtils.GetLevelModoverridesPath(clusterName, levelName)
	err := fileUtils.CreateFileIfNotExists(modoverridesPath)
	if err != nil {
		log.Panicln("创建失败 modoverrides ", err)
	}
	err = fileUtils.WriterTXT(modoverridesPath, modoverrides)
	if err != nil {
		log.Panicln("写入失败 modoverrides", err)
	}
}

func (g *GameLevelService) SaveServerIni(clusterName, levelName string, serverIni *level.ServerIni) {
	levelServerIniPath := dstUtils.GetLevelServerIniPath(clusterName, levelName)
	err := fileUtils.CreateFileIfNotExists(levelServerIniPath)
	if err != nil {
		log.Panicln("创建失败 serverIni ", err)
	}
	err = fileUtils.WriterTXT(levelServerIniPath, dstUtils.ParseTemplate(consts.ServerIniTemplate, serverIni))
	if err != nil {
		log.Panicln("写入失败 serverIni ", err)
	}
}
