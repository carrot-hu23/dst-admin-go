package service

import (
	"dst-admin-go/constant"
	"dst-admin-go/utils/dstConfigUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo"
	"log"
	"path"
	"strconv"
	"strings"
)

const START_NEW_GAME uint8 = 1
const SAVE_RESTART uint8 = 2

var cluster_init_template = "./static/template/cluster2.ini"
var master_server_init_template = "./static/template/master_server.ini"
var caves_server_init_template = "./static/template/caves_server.ini"

type GameConfigService struct {
	ClusterService
}

func (c *GameConfigService) GetConfig() vo.GameConfigVO {
	gameConfig := vo.NewGameConfigVO()

	gameConfig.Token = c.getClusterToken()
	c.GetClusterIni(gameConfig)
	gameConfig.MasterMapData = c.getMasteLeveldataoverride()
	gameConfig.CavesMapData = c.getCavesLeveldataoverride()
	gameConfig.ModData = c.getModoverrides()

	return *gameConfig
}

func (c *GameConfigService) getClusterToken() string {
	fileUtils.CreateFileIfNotExists(constant.GET_CLUSTER_TOKEN_PATH())
	token, err := fileUtils.ReadFile(constant.GET_CLUSTER_TOKEN_PATH())
	if err != nil {
		panic("read cluster_token.txt file error: " + err.Error())
	}

	return token
}

func (c *GameConfigService) GetClusterIni(gameconfig *vo.GameConfigVO) {
	if !fileUtils.Exists(constant.GET_CLUSTER_INI_PATH()) {
		fileUtils.CreateFileIfNotExists(constant.GET_CLUSTER_INI_PATH())
		return
	}
	clusterIni, err := fileUtils.ReadLnFile(constant.GET_CLUSTER_INI_PATH())
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

func (c *GameConfigService) getMasteLeveldataoverride() string {

	filepath := constant.GET_MASTER_LEVELDATAOVERRIDE_PATH()
	if !fileUtils.Exists(filepath) {
		fileUtils.CreateFileIfNotExists(filepath)
		return "return {}"
	}

	level, err := fileUtils.ReadFile(constant.GET_MASTER_LEVELDATAOVERRIDE_PATH())
	if err != nil {
		panic("read Master/leveldataoverride.lua file error: " + err.Error())
	}
	return level
}

func (c *GameConfigService) getCavesLeveldataoverride() string {

	filepath := constant.GET_CAVES_LEVELDATAOVERRIDE_PATH()
	if !fileUtils.Exists(filepath) {
		fileUtils.CreateFileIfNotExists(filepath)
		return "return {}"
	}

	level, err := fileUtils.ReadFile(constant.GET_CAVES_LEVELDATAOVERRIDE_PATH())
	if err != nil {
		panic("read Caves/leveldataoverride.lua file error: " + err.Error())
	}
	return level
}

func (c *GameConfigService) getModoverrides() string {

	filepath := constant.GET_MASTER_MOD_PATH()
	if !fileUtils.Exists(filepath) {
		fileUtils.CreateFileIfNotExists(filepath)
		return "return {}"
	}

	level, err := fileUtils.ReadFile(constant.GET_MASTER_MOD_PATH())
	if err != nil {
		panic("read Master/modoverrides.lua file error: " + err.Error())
	}
	return level
}

func (c *GameConfigService) SaveConfig(gameConfigVo vo.GameConfigVO) {

	//创建存档目录
	c.createMyDediServerDir()
	//创建房间配置
	c.createClusterIni(gameConfigVo)
	//创建token配置
	c.createClusterToken(strings.TrimSpace(gameConfigVo.Token))
	//创建地面和洞穴的ini配置文件
	// createMasterServerIni()
	// createCavesServerIni()
	//创建地面世界设置
	c.createMasteLeveldataoverride(gameConfigVo.MasterMapData)
	//创建洞穴世界设置
	c.createCavesLeveldataoverride(gameConfigVo.CavesMapData)
	//创建mod设置
	c.createModoverrides(gameConfigVo.ModData)
	//TODO dedicated_server_mods_setup

}

func (c *GameConfigService) createMyDediServerDir() {
	dstConfig := dstConfigUtils.GetDstConfig()
	basePath := constant.GET_DST_USER_GAME_CONFG_PATH()
	myDediServerPath := path.Join(basePath, dstConfig.Cluster)
	log.Println("生成 myDediServer 目录：" + myDediServerPath)
	fileUtils.CreateDir(myDediServerPath)
}

func (c *GameConfigService) createClusterIni(gameConfigVo vo.GameConfigVO) {

	log.Println("生成游戏配置文件 cluster.ini文件: ", constant.GET_CLUSTER_INI_PATH())
	oldCluster := c.ReadClusterIniFile()

	oldCluster.ClusterName = gameConfigVo.ClusterName
	oldCluster.ClusterDescription = gameConfigVo.ClusterDescription
	oldCluster.GameMode = gameConfigVo.GameMode
	oldCluster.MaxPlayers = uint(gameConfigVo.MaxPlayers)
	oldCluster.Pvp = gameConfigVo.Pvp
	oldCluster.VoteEnabled = gameConfigVo.VoteEnabled
	oldCluster.PauseWhenNobody = gameConfigVo.PauseWhenNobody
	oldCluster.ClusterPassword = gameConfigVo.ClusterPassword

	clusterIni := c.ParseTemplate(cluster_init_template, oldCluster)
	fileUtils.WriterTXT(constant.GET_CLUSTER_INI_PATH(), clusterIni)
}

func (c *GameConfigService) createClusterToken(token string) {
	log.Println("生成cluster_token.txt 文件 ", constant.GET_CLUSTER_TOKEN_PATH())
	fileUtils.WriterTXT(constant.GET_CLUSTER_TOKEN_PATH(), token)
}

func (c *GameConfigService) createMasteLeveldataoverride(mapConfig string) {

	log.Println("生成master_leveldataoverride.txt 文件 ", constant.GET_MASTER_LEVELDATAOVERRIDE_PATH())
	if mapConfig != "" {
		fileUtils.WriterTXT(constant.GET_MASTER_LEVELDATAOVERRIDE_PATH(), mapConfig)
	} else {
		//置空
		fileUtils.WriterTXT(constant.GET_MASTER_LEVELDATAOVERRIDE_PATH(), "")
	}
}
func (c *GameConfigService) createCavesLeveldataoverride(mapConfig string) {

	log.Println("生成caves_leveldataoverride.lua 文件 ", constant.GET_CAVES_LEVELDATAOVERRIDE_PATH())
	if mapConfig != "" {
		fileUtils.WriterTXT(constant.GET_CAVES_LEVELDATAOVERRIDE_PATH(), mapConfig)
	} else {
		//置空
		fileUtils.WriterTXT(constant.GET_CAVES_LEVELDATAOVERRIDE_PATH(), "")
	}
}
func (c *GameConfigService) createModoverrides(modConfig string) {

	log.Println("生成master_modoverrides.lua 文件 ", constant.GET_MASTER_MOD_PATH())
	log.Println("生成caves_modoverrides.lua 文件 ", constant.GET_CAVES_MOD_PATH())
	if modConfig != "" {
		fileUtils.WriterTXT(constant.GET_MASTER_MOD_PATH(), modConfig)
		fileUtils.WriterTXT(constant.GET_CAVES_MOD_PATH(), modConfig)

		var serverModSetup = ""
		//TODO 添加mod setup
		workshopIds := WorkshopIds(modConfig)
		for _, workshopId := range workshopIds {
			serverModSetup += "ServerModSetup(\"" + workshopId + "\")\n"
		}
		fileUtils.WriterTXT(constant.GET_DST_MOD_SETUP_PATH(), serverModSetup)
	} else {
		//置空
		fileUtils.WriterTXT(constant.GET_MASTER_MOD_PATH(), "")
		fileUtils.WriterTXT(constant.GET_CAVES_MOD_PATH(), "")
	}
}

func (c *GameConfigService) UpdateDedicatedServerModsSetup(modConfig string) {
	if modConfig != "" {
		var serverModSetup = ""
		workshopIds := WorkshopIds(modConfig)
		for _, workshopId := range workshopIds {
			serverModSetup += "ServerModSetup(\"" + workshopId + "\")\n"
		}
		fileUtils.WriterTXT(constant.GET_DST_MOD_SETUP_PATH(), serverModSetup)
	}

}
