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

type InitDstData struct {
	// InstallDstEnv bool                      `json:"installDstEnv"`
	// ClusterToken  string                    `json:"clusterToken"`
	DstConfig *dstConfigUtils.DstConfig `json:"dstConfig"`
	UserInfo  *vo.UserInfo              `json:"userInfo"`
}

func InstallSteamCmd() error {
	cmd := exec.Command("./static/install.sh")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("执行./static/install.sh脚本失败：", err)
		return err
	}
	fmt.Println("./static/install.sh脚本输出：", string(output))
	return err
}

func InstallSteamCmdAndDst() map[string]string {
	//安装 steam cmd 和 dst
	log.Println("installing steamcmd")
	err := InstallSteamCmd()
	if err != nil {
		log.Panicln("安装失败")
	}
	steamcmdPath := path.Join(constant.HOME_PATH, "steamcmd")
	dstPath := path.Join(constant.HOME_PATH, "dontstarve_dedicated_server")

	return map[string]string{"steamcmdPath": steamcmdPath, "dstPath": dstPath}
}

func InitDstEnv(initDst *InitDstData, ctx *gin.Context) {

	InitUserInfo(initDst.UserInfo)
	InitDstConfig(initDst.DstConfig)
	InitBaseLevel(initDst.DstConfig, initDst.UserInfo.Username, "", false)

	log.Println("创建完成")
}

func InitDstConfig(dstConfig *dstConfigUtils.DstConfig) {

	if dstConfig.Backup == "" {
		dstConfig.Backup = filepath.Join(constant.HOME_PATH, ".klei", "DoNotStarveTogether")
	}
	if dstConfig.Backup == "" {
		dstConfig.Mod_download_path = filepath.Join(constant.HOME_PATH, ".klei", "DoNotStarveTogether", "mod_download")
		createDirIfNotExsists(dstConfig.Mod_download_path)
	}
	dstConfigUtils.SaveDstConfig(dstConfig)
}

func InitBaseLevel(dstConfig *dstConfigUtils.DstConfig, username, token string, exsitesNotInit bool) {
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

	createDirIfNotExsists(baseLevelPath)
	createDirIfNotExsists(dstConfig.Backup)
	createDirIfNotExsists(dstConfig.Mod_download_path)

	log.Println(baseLevelPath)

	InitClusterIni(baseLevelPath, username)
	InitClusterToken(baseLevelPath, token)
	InitBaseMaster(baseLevelPath)
	InitBaseCaves(baseLevelPath)
}

func InitClusterIni(basePath string, username string) {
	cluster_ini_path := filepath.Join(basePath, "cluster.ini")
	createFileIfNotExsists(cluster_ini_path)
	clusterIni := cluster.NewCluster()
	clusterName := ""
	if username != "" {
		clusterName = username + "的世界"
	} else {
		clusterName = "我的饥荒服务世界"
	}
	clusterIni.ClusterName = clusterName
	clusterIni.MaxPlayers = 8
	fileUtils.WriterTXT(cluster_ini_path, pareseTemplate(cluster_template, clusterIni))
}

func InitClusterToken(basePath string, token string) {
	cluster_token_path := filepath.Join(basePath, "cluster_token.txt")
	createFileIfNotExsists(cluster_token_path)
	fileUtils.WriterTXT(token, cluster_token_path)
}

func InitBaseMaster(basePath string) {

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

	createDirIfNotExsists(filepath.Join(basePath, "Master"))

	createFileIfNotExsists(l_path)
	createFileIfNotExsists(m_path)
	createFileIfNotExsists(s_path)

	fileUtils.WriterTXT(l_path, leveldataoverride)
	fileUtils.WriterTXT(m_path, modoverrides)
	fileUtils.WriterTXT(s_path, server_ini)
}

func InitBaseCaves(basePath string) {

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

	createDirIfNotExsists(filepath.Join(basePath, "Caves"))

	createFileIfNotExsists(l_path)
	createFileIfNotExsists(m_path)
	createFileIfNotExsists(s_path)

	fileUtils.WriterTXT(l_path, leveldataoverride)
	fileUtils.WriterTXT(m_path, modoverrides)
	fileUtils.WriterTXT(s_path, server_ini)
}
