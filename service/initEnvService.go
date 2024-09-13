package service

import (
	"dst-admin-go/config/database"
	"dst-admin-go/constant"
	"dst-admin-go/model"
	"dst-admin-go/utils/dstConfigUtils"
	"dst-admin-go/utils/dstUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/utils/levelConfigUtils"
	"dst-admin-go/utils/systemUtils"
	"dst-admin-go/vo"
	"dst-admin-go/vo/level"
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/gin-gonic/gin"
)

const (
	cluster_template  = "./static/template/cluster2.ini"
	master_template   = "./static/Master"
	caves_template    = "./static/Caves"
	porkland_template = "./static/Porkland"
)

var gameLevel2Service GameLevel2Service
var clusterManger ClusterManager

type InitService struct {
	GameConfigService
	LoginService
}

type InitDstData struct {
	DstConfig *dstConfigUtils.DstConfig `json:"dstConfig"`
	UserInfo  *vo.UserInfo              `json:"userInfo"`
}

func (i *InitService) InitDstEnv(initDst *InitDstData, ctx *gin.Context) {

	i.InitUserInfo(initDst.UserInfo)
	// TODO 写入默认配置
	log.Println("初始化用户完成")
}

func (i *InitService) InitDstConfig(dstConfig *dstConfigUtils.DstConfig) {

	if dstConfig.Backup == "" {
		dstConfig.Backup = filepath.Join(constant.HOME_PATH, ".klei", "DoNotStarveTogether")
	}
	if dstConfig.Mod_download_path == "" {
		dstConfig.Mod_download_path = filepath.Join(constant.HOME_PATH, ".klei", "DoNotStarveTogether", "mod_download")
		fileUtils.CreateDirIfNotExists(dstConfig.Mod_download_path)
	}
	dstConfigUtils.SaveDstConfig(dstConfig)
}

func (i *InitService) InitBaseLevel(dstConfig *dstConfigUtils.DstConfig, username, token string, exsitesNotInit bool) {
	clusterName := dstConfig.Cluster
	klei_path := ""
	if runtime.GOOS == "windows" {
		klei_path = filepath.Join(constant.HOME_PATH, "Documents", "klei", "DoNotStarveTogether")
	} else {
		klei_path = filepath.Join(constant.HOME_PATH, ".klei", "DoNotStarveTogether")
	}
	baseLevelPath := filepath.Join(klei_path, clusterName)

	if exsitesNotInit {
		if fileUtils.Exists(baseLevelPath) {
			return
		}
	}

	fileUtils.CreateDirIfNotExists(baseLevelPath)
	fileUtils.CreateDirIfNotExists(dstConfig.Backup)
	fileUtils.CreateDirIfNotExists(dstConfig.Mod_download_path)

	log.Println(baseLevelPath)

	i.InitClusterIni(baseLevelPath, username)
	i.InitClusterToken(baseLevelPath, token)
	i.InitBaseMaster(baseLevelPath)
	i.InitBaseCaves(baseLevelPath)
}

func (i *InitService) InitClusterIni(basePath string, username string) {
	cluster_ini_path := filepath.Join(basePath, "cluster.ini")
	fileUtils.CreateFileIfNotExists(cluster_ini_path)
	clusterIni := level.NewClusterIni()
	clusterName := ""
	if username != "" {
		clusterName = username + "的世界"
	} else {
		clusterName = "我的饥荒服务世界"
	}
	clusterIni.ClusterName = clusterName
	clusterIni.GameMode = "survival"
	clusterIni.MaxPlayers = 6
	fileUtils.WriterTXT(cluster_ini_path, dstUtils.ParseTemplate(cluster_template, clusterIni))
}

func (i *InitService) InitClusterToken(basePath string, token string) {
	cluster_token_path := filepath.Join(basePath, "cluster_token.txt")
	fileUtils.CreateFileIfNotExists(cluster_token_path)
	fileUtils.WriterTXT(cluster_token_path, token)
}

func (i *InitService) InitBaseMaster(basePath string) {

	leveldataoverride, err := fileUtils.ReadFile(filepath.Join(master_template, "leveldataoverride.lua"))
	if err != nil {
		panic("read ./static/Master/leveldataoverride.lua file error: " + err.Error())
	}
	modoverrides, err := fileUtils.ReadFile(filepath.Join(master_template, "modoverrides.lua"))
	if err != nil {
		panic("read ./static/Master/modoverrides.lua file error: " + err.Error())
	}
	server_ini, err := fileUtils.ReadFile(filepath.Join(master_template, "server.ini"))
	if err != nil {
		panic("read /static/Master/server.ini file error: " + err.Error())
	}

	l_path := filepath.Join(basePath, "Master", "leveldataoverride.lua")
	m_path := filepath.Join(basePath, "Master", "modoverrides.lua")
	s_path := filepath.Join(basePath, "Master", "server.ini")

	fileUtils.CreateDirIfNotExists(filepath.Join(basePath, "Master"))

	fileUtils.CreateFileIfNotExists(l_path)
	fileUtils.CreateFileIfNotExists(m_path)
	fileUtils.CreateFileIfNotExists(s_path)

	fileUtils.WriterTXT(l_path, leveldataoverride)
	fileUtils.WriterTXT(m_path, modoverrides)
	fileUtils.WriterTXT(s_path, server_ini)
}

func (i *InitService) InitBaseCaves(basePath string) {

	leveldataoverride, err := fileUtils.ReadFile(filepath.Join(caves_template, "leveldataoverride.lua"))
	if err != nil {
		panic("read ./static/Caves/leveldataoverride.lua file error: " + err.Error())
	}
	modoverrides, err := fileUtils.ReadFile(filepath.Join(caves_template, "modoverrides.lua"))
	if err != nil {
		panic("read ./static/Caves/modoverrides.lua file error: " + err.Error())
	}
	server_ini, err := fileUtils.ReadFile(filepath.Join(caves_template, "server.ini"))
	if err != nil {
		panic("read /static/Caves/server.ini file error: " + err.Error())
	}

	l_path := filepath.Join(basePath, "Caves", "leveldataoverride.lua")
	m_path := filepath.Join(basePath, "Caves", "modoverrides.lua")
	s_path := filepath.Join(basePath, "Caves", "server.ini")

	fileUtils.CreateDirIfNotExists(filepath.Join(basePath, "Caves"))

	fileUtils.CreateFileIfNotExists(l_path)
	fileUtils.CreateFileIfNotExists(m_path)
	fileUtils.CreateFileIfNotExists(s_path)

	fileUtils.WriterTXT(l_path, leveldataoverride)
	fileUtils.WriterTXT(m_path, modoverrides)
	fileUtils.WriterTXT(s_path, server_ini)
}

func (i *InitService) InitCluster(cluster *model.Cluster, token string) {

	kleiPath := filepath.Join(constant.HOME_PATH, ".klei", "DoNotStarveTogether")
	baseLevelPath := filepath.Join(kleiPath, cluster.ClusterName)
	if fileUtils.Exists(baseLevelPath) {
		return
	}

	fileUtils.CreateDirIfNotExists(baseLevelPath)
	fileUtils.CreateDirIfNotExists(cluster.Backup)
	fileUtils.CreateDirIfNotExists(cluster.ModDownloadPath)

	// info := i.GetUserInfo()
	i.InitClusterIni(baseLevelPath, cluster.Name)
	i.InitClusterToken(baseLevelPath, token)
	i.InitBaseMaster(baseLevelPath)
	i.InitBaseCaves(baseLevelPath)
}

func (i *InitService) InitCluster2(cluster *model.Cluster, token string) {

	kleiPath := filepath.Join(constant.HOME_PATH, ".klei", "DoNotStarveTogether")
	baseLevelPath := filepath.Join(kleiPath, cluster.ClusterName)
	if fileUtils.Exists(baseLevelPath) {
		return
	}

	fileUtils.CreateDirIfNotExists(baseLevelPath)
	fileUtils.CreateDirIfNotExists(cluster.Backup)
	fileUtils.CreateDirIfNotExists(cluster.ModDownloadPath)

	// info := i.GetUserInfo()
	i.InitClusterIni(baseLevelPath, cluster.Name)
	i.InitClusterToken(baseLevelPath, token)

	if cluster.LevelType == "porkland" {
		ports := i.getNextPorts()
		// 初始化标准的森林和洞穴
		i.InitBaseNewLevel(filepath.Join(baseLevelPath, "Master"), ports[0], true, "porkland", "Porkland", 1)
		levelConfig := levelConfigUtils.LevelConfig{
			LevelList: []levelConfigUtils.Item{
				{
					Name: "猪镇",
					File: "Master",
				},
			},
		}
		levelConfigUtils.SaveLevelConfig(cluster.ClusterName, &levelConfig)
	} else {
		ports := i.getNextPorts()
		// 初始化标准的森林和洞穴
		i.InitBaseNewLevel(filepath.Join(baseLevelPath, "Master"), ports[0], true, "forest", "Master", 1)
		i.InitBaseNewLevel(filepath.Join(baseLevelPath, "Caves"), ports[0], true, "caves", "Caves", 2)

		levelConfig := levelConfigUtils.LevelConfig{
			LevelList: []levelConfigUtils.Item{
				{
					Name: "森林",
					File: "Master",
				},
				{
					Name: "洞穴",
					File: "Caves",
				},
			},
		}
		levelConfigUtils.SaveLevelConfig(cluster.ClusterName, &levelConfig)
	}

}

func (i *InitService) getNextPorts() (portList []int) {

	// 找到本机可用的udp 10998-11018
	ports, err := systemUtils.FindFreeUDPPorts(10998, 11038)
	//log.Println(err)
	//if err != nil {
	//	portList = append(portList, 10998, 10999, 10997)
	//	return portList
	//}

	// 过滤目前存档已经使用了的端口
	db := database.DB
	clusters := make([]model.Cluster, 0)
	if err = db.Find(&clusters).Error; err != nil {
		fmt.Println(err.Error())
	}
	var userPorts []uint
	for i1 := range clusters {
		clusterName := clusters[i1].ClusterName
		levelList := gameLevel2Service.GetLevelList(clusterName)
		for i2 := range levelList {
			userPorts = append(userPorts, levelList[i2].ServerIni.ServerPort)
		}
	}
	log.Println("userPorts", userPorts, "ports", ports)
	// 将切片 b 转换为集合（map）
	bSet := make(map[uint]struct{})
	for _, v := range userPorts {
		bSet[v] = struct{}{}
	}

	// 过滤切片 a
	var filterPorts []int
	for _, v := range ports {
		if _, found := bSet[uint(v)]; !found {
			filterPorts = append(filterPorts, v)
		}
	}
	log.Println("filterPorts", filterPorts)
	if len(filterPorts) > 1 {
		if len(filterPorts) > 2 {
			portList = append(portList, filterPorts[0], filterPorts[1], filterPorts[2])
		} else {
			portList = append(portList, filterPorts[0], filterPorts[1])
		}
	} else {
		portList = append(portList, 10998, 10999, 10997)
	}
	log.Println("portList", portList)
	return portList

}

func (i *InitService) InitBaseNewLevel(basePath string, port int, isMaster bool, levelType string, levelName string, id int) {

	template := ""
	if levelType == "porkland" {
		template = porkland_template
	} else if levelType == "caves" {
		template = caves_template
	} else {
		template = master_template
	}

	leveldataoverride, err := fileUtils.ReadFile(filepath.Join(template, "leveldataoverride.lua"))
	if err != nil {
		panic(err.Error())
	}
	modoverrides, err := fileUtils.ReadFile(filepath.Join(template, "modoverrides.lua"))
	if err != nil {
		panic(err.Error())
	}

	serverIni := level.ServerIni{
		ServerPort:         uint(port),
		IsMaster:           isMaster,
		Name:               levelName,
		Id:                 uint(id),
		EncodeUserPath:     true,
		AuthenticationPort: 8766,
		MasterServerPort:   27016,
	}

	server_ini := dstUtils.ParseTemplate(ServerIniTemplate, serverIni)

	l_path := filepath.Join(basePath, "leveldataoverride.lua")
	m_path := filepath.Join(basePath, "modoverrides.lua")
	s_path := filepath.Join(basePath, "server.ini")

	fileUtils.CreateDirIfNotExists(filepath.Join(basePath, "Master"))

	fileUtils.CreateFileIfNotExists(l_path)
	fileUtils.CreateFileIfNotExists(m_path)
	fileUtils.CreateFileIfNotExists(s_path)

	fileUtils.WriterTXT(l_path, leveldataoverride)
	fileUtils.WriterTXT(m_path, modoverrides)
	fileUtils.WriterTXT(s_path, server_ini)

}
