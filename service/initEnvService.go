package service

import (
	"dst-admin-go/constant"
	"dst-admin-go/utils/dstConfigUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo"
	"dst-admin-go/vo/cluster"
	"fmt"
	"log"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"

	"github.com/gin-gonic/gin"
)

const (
	cluster_template = "./static/template/cluster2.ini"
	master_template  = "./static/Master"
	caves_template   = "./static/Caves"

	// cluster_template = "C:\\Users\\xm\\Desktop\\dst-admin-go\\static\\template\\cluster2.ini"
	// master_template  = "C:\\Users\\xm\\Desktop\\dst-admin-go\\static\\Master"
	// caves_template   = "C:\\Users\\xm\\Desktop\\dst-admin-go\\static\\Caves"
)

type InitService struct {
	GameConfigService
	LoginService
}

type InitDstData struct {
	// InstallDstEnv bool                      `json:"installDstEnv"`
	// ClusterToken  string                    `json:"clusterToken"`
	DstConfig *dstConfigUtils.DstConfig `json:"dstConfig"`
	UserInfo  *vo.UserInfo              `json:"userInfo"`
}

func (i *InitService) InstallSteamCmd() error {
	cmd := exec.Command("./static/install.sh")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("执行./static/install.sh脚本失败：", err)
		return err
	}
	fmt.Println("./static/install.sh脚本输出：", string(output))
	return err
}

func (i *InitService) InstallSteamCmdAndDst() map[string]string {
	//安装 steam cmd 和 dst
	log.Println("installing steamcmd")
	err := i.InstallSteamCmd()
	if err != nil {
		log.Panicln("安装失败")
	}
	steamcmdPath := path.Join(constant.HOME_PATH, "steamcmd")
	dstPath := path.Join(constant.HOME_PATH, "dontstarve_dedicated_server")

	return map[string]string{"steamcmdPath": steamcmdPath, "dstPath": dstPath}
}

func (i *InitService) InitDstEnv(initDst *InitDstData, ctx *gin.Context) {

	i.InitUserInfo(initDst.UserInfo)
	i.InitDstConfig(initDst.DstConfig)
	i.InitBaseLevel(initDst.DstConfig, initDst.UserInfo.Username, "", false)

	log.Println("创建完成")
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
	clusterIni := cluster.NewCluster()
	clusterName := ""
	if username != "" {
		clusterName = username + "的世界"
	} else {
		clusterName = "我的饥荒服务世界"
	}
	clusterIni.ClusterName = clusterName
	clusterIni.MaxPlayers = 8
	fileUtils.WriterTXT(cluster_ini_path, i.ParseTemplate(cluster_template, clusterIni))
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

	fileUtils.CreateDirIfNotExists(l_path)
	fileUtils.CreateDirIfNotExists(m_path)
	fileUtils.CreateDirIfNotExists(s_path)

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
