package service

import (
	"dst-admin-go/constant"
	"dst-admin-go/constant/dst"
	"dst-admin-go/utils/dstConfigUtils"
	"dst-admin-go/utils/dstUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo"
	"log"
	"path"
	"strconv"
	"strings"
)

var cluster_init_template = "./static/template/cluster2.ini"
var master_server_init_template = "./static/template/master_server.ini"
var caves_server_init_template = "./static/template/caves_server.ini"

type GameConfigService struct {
	w HomeService
}

func (c *GameConfigService) GetConfig(clusterName string) vo.GameConfigVO {

	gameConfig := vo.NewGameConfigVO()
	gameConfig.Token = c.getClusterToken(clusterName)
	c.GetClusterIni(clusterName, gameConfig)
	gameConfig.MasterMapData = c.getMasterLeveldataoverride(clusterName)
	gameConfig.CavesMapData = c.getCavesLeveldataoverride(clusterName)
	gameConfig.ModData = c.getModoverrides(clusterName)

	return *gameConfig
}

func (c *GameConfigService) getClusterToken(clusterName string) string {
	clusterToken := dst.GetClusterTokenPath(clusterName)
	token, err := fileUtils.ReadFile(clusterToken)
	if err != nil {
		panic("read cluster_token.txt file error: " + err.Error())
	}

	return token
}

func (c *GameConfigService) GetClusterIni(clusterName string, gameconfig *vo.GameConfigVO) {

	clusterIniPath := dst.GetClusterIniPath(clusterName)
	clusterIni, err := fileUtils.ReadLnFile(clusterIniPath)
	if err != nil {
		panic("read cluster.ini file error: " + err.Error())
	}
	for _, value := range clusterIni {
		if value == "" {
			continue
		}
		if strings.Contains(value, "game_mod") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				gameconfig.GameMode = s
			}
		}
		if strings.Contains(value, "max_players") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				n, err := strconv.ParseUint(s, 10, 8)
				if err == nil {
					gameconfig.MaxPlayers = uint8(n)
				}
			}
		}
		if strings.Contains(value, "pvp") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				b, err := strconv.ParseBool(s)
				if err == nil {
					gameconfig.Pvp = b
				}
			}
		}
		if strings.Contains(value, "pause_when_empty") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				b, err := strconv.ParseBool(s)
				if err == nil {
					gameconfig.PauseWhenNobody = b
				}
			}
		}
		if strings.Contains(value, "cluster_intention") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				gameconfig.ClusterIntention = s
			}
		}
		if strings.Contains(value, "cluster_password") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				gameconfig.ClusterPassword = s
			}
		}
		if strings.Contains(value, "cluster_description") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				gameconfig.ClusterDescription = s
			}
		}
		if strings.Contains(value, "cluster_name") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				gameconfig.ClusterName = s
			}
		}

	}
}

func (c *GameConfigService) getMasterLeveldataoverride(clusterName string) string {

	leveldataoverridePath := dst.GetMasterLeveldataoverridePath(clusterName)

	level, err := fileUtils.ReadFile(leveldataoverridePath)
	if err != nil {
		panic("read Master/leveldataoverride.lua file error: " + err.Error())
	}
	return level
}

func (c *GameConfigService) getCavesLeveldataoverride(clusterName string) string {

	leveldataoverridePath := dst.GetCavesLeveldataoverridePath(clusterName)
	level, err := fileUtils.ReadFile(leveldataoverridePath)
	if err != nil {
		panic("read Caves/leveldataoverride.lua file error: " + err.Error())
	}
	return level
}

func (c *GameConfigService) getModoverrides(clusterName string) string {

	modoverridesPath := dst.GetMasterModoverridesPath(clusterName)
	modoverrides, err := fileUtils.ReadFile(modoverridesPath)
	if err != nil {
		panic("read Master/modoverrides.lua file error: " + err.Error())
	}
	return modoverrides
}

func (c *GameConfigService) SaveConfig(clusterName string, gameConfigVo vo.GameConfigVO) {

	//创建房间配置
	c.createClusterIni(clusterName, gameConfigVo)
	//创建token配置
	c.createClusterToken(clusterName, strings.TrimSpace(gameConfigVo.Token))
	//创建地面和洞穴的ini配置文件
	// createMasterServerIni()
	// createCavesServerIni()
	//创建地面世界设置
	c.createMasteLeveldataoverride(clusterName, gameConfigVo.MasterMapData)
	//创建洞穴世界设置
	c.createCavesLeveldataoverride(clusterName, gameConfigVo.CavesMapData)
	//创建mod设置
	c.createModoverrides(clusterName, gameConfigVo.ModData)
	//TODO dedicated_server_mods_setup

}

func (c *GameConfigService) createMyDediServerDir() {
	dstConfig := dstConfigUtils.GetDstConfig()
	basePath := constant.GET_DST_USER_GAME_CONFG_PATH()
	myDediServerPath := path.Join(basePath, dstConfig.Cluster)
	log.Println("生成 myDediServer 目录：" + myDediServerPath)
	fileUtils.CreateDir(myDediServerPath)
}

func (c *GameConfigService) createClusterIni(clusterName string, gameConfigVo vo.GameConfigVO) {
	clusterIniPath := dst.GetClusterIniPath(clusterName)
	log.Println("生成游戏配置文件 cluster.ini文件: ", clusterIniPath)
	oldCluster := c.w.GetClusterIni(clusterName)

	oldCluster.ClusterName = gameConfigVo.ClusterName
	oldCluster.ClusterDescription = gameConfigVo.ClusterDescription
	oldCluster.GameMode = gameConfigVo.GameMode
	oldCluster.MaxPlayers = uint(gameConfigVo.MaxPlayers)
	oldCluster.Pvp = gameConfigVo.Pvp
	oldCluster.VoteEnabled = gameConfigVo.VoteEnabled
	oldCluster.PauseWhenNobody = gameConfigVo.PauseWhenNobody
	oldCluster.ClusterPassword = gameConfigVo.ClusterPassword

	clusterIni := dstUtils.ParseTemplate(cluster_init_template, oldCluster)
	fileUtils.WriterTXT(clusterIniPath, clusterIni)
}

func (c *GameConfigService) createClusterToken(clusterName, token string) {
	fileUtils.WriterTXT(dst.GetClusterTokenPath(clusterName), token)
}

func (c *GameConfigService) createMasteLeveldataoverride(clusterName string, mapConfig string) {
	leveldataoverridePath := dst.GetMasterLeveldataoverridePath(clusterName)
	log.Println("生成master_leveldataoverride.txt 文件 ", leveldataoverridePath)
	if mapConfig != "" {
		fileUtils.WriterTXT(leveldataoverridePath, mapConfig)
	} else {
		//置空
		fileUtils.WriterTXT(leveldataoverridePath, "")
	}
}
func (c *GameConfigService) createCavesLeveldataoverride(clusterName string, mapConfig string) {
	leveldataoverridePath := dst.GetCavesLeveldataoverridePath(clusterName)
	log.Println("生成caves_leveldataoverride.lua 文件 ", leveldataoverridePath)
	if mapConfig != "" {
		fileUtils.WriterTXT(leveldataoverridePath, mapConfig)
	} else {
		//置空
		fileUtils.WriterTXT(leveldataoverridePath, "")
	}
}
func (c *GameConfigService) createModoverrides(clusterName string, modConfig string) {

	masterModoverridesPath := dst.GetMasterModoverridesPath(clusterName)
	cavesModoverridesPath := dst.GetCavesModoverridesPath(clusterName)

	if modConfig != "" {
		fileUtils.WriterTXT(masterModoverridesPath, modConfig)
		fileUtils.WriterTXT(cavesModoverridesPath, modConfig)

		var serverModSetup = ""
		//TODO 添加mod setup
		workshopIds := dstUtils.WorkshopIds(modConfig)
		for _, workshopId := range workshopIds {
			serverModSetup += "ServerModSetup(\"" + workshopId + "\")\n"
		}
		fileUtils.WriterTXT(dst.GetModSetup(clusterName), serverModSetup)
	} else {
		//置空
		fileUtils.WriterTXT(masterModoverridesPath, "")
		fileUtils.WriterTXT(cavesModoverridesPath, "")
	}
}

func (c *GameConfigService) UpdateDedicatedServerModsSetup(clusterName, modConfig string) {
	if modConfig != "" {
		var serverModSetup = ""
		workshopIds := dstUtils.WorkshopIds(modConfig)
		for _, workshopId := range workshopIds {
			serverModSetup += "ServerModSetup(\"" + workshopId + "\")\n"
		}
		fileUtils.WriterTXT(dst.GetModSetup(clusterName), serverModSetup)
	}

}
