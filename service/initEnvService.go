package service

import (
	"dst-admin-go/constant"
	"dst-admin-go/model"
	"dst-admin-go/utils/dstConfigUtils"
	"dst-admin-go/utils/dstUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo"
	"dst-admin-go/vo/level"
	"log"
	"path/filepath"
	"runtime"

	"github.com/gin-gonic/gin"
)

const (
	cluster_template = "./static/template/cluster2.ini"
	master_template  = "./static/Master"
	caves_template   = "./static/Caves"
)

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
	clusterIni.MaxPlayers = 8
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

	info := i.GetUserInfo()
	i.InitClusterIni(baseLevelPath, info["displayName"].(string))
	i.InitClusterToken(baseLevelPath, token)
	i.InitBaseMaster(baseLevelPath)
	i.InitBaseCaves(baseLevelPath)
}
