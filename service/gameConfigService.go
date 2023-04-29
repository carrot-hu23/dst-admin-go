package service

import (
	"bytes"
	"dst-admin-go/constant"
	optype "dst-admin-go/constant/opType"
	"dst-admin-go/utils/dstConfigUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo"
	"fmt"
	"html/template"
	"log"
	"strconv"
	"strings"
)

const START_NEW_GAME uint8 = 1
const SAVE_RESTART uint8 = 2

var cluster_init_template = "./static/template/cluster.ini"
var master_server_init_template = "./static/template/master_server.ini"
var caves_server_init_template = "./static/template/caves_server.ini"

func Dst_user_game_confg_path() string {
	cluster := dstConfigUtils.GetDstConfig().Cluster
	var path = constant.HOME_PATH + "/.klei/DoNotStarveTogether/" + cluster + "/"
	return path
}

func GetClusterTokenPath() string {
	return Dst_user_game_confg_path() + constant.DST_USER_CLUSTER_TOKEN
}

func GetClusterIniPath() string {
	return Dst_user_game_confg_path() + constant.DST_USER_CLUSTER_INI_NAME
}

func GetMasterDirPath() string {
	return Dst_user_game_confg_path() + constant.DST_MASTER
}

func GetMasterDirServerIniPath() string {
	return GetMasterDirPath() + constant.SINGLE_SLASH + constant.DST_USER_SERVER_INI_NAME
}

func GetCavesDirPath() string {
	return Dst_user_game_confg_path() + constant.DST_CAVES
}

func GetCavesDirServerIniPath() string {
	return GetCavesDirPath() + constant.SINGLE_SLASH + constant.DST_USER_SERVER_INI_NAME
}

func GetMasteLeveldataoverridePath() string {
	return Dst_user_game_confg_path() + "/" + constant.DST_MASTER + "/leveldataoverride.lua"
}

func GetCavesLeveldataoverridePath() string {
	return Dst_user_game_confg_path() + "/" + constant.DST_CAVES + "/leveldataoverride.lua"
}

func GetMasterModPath() string {
	return Dst_user_game_confg_path() + "/" + constant.DST_MASTER + "/modoverrides.lua"
}

func GetCavesModPath() string {
	return Dst_user_game_confg_path() + "/" + constant.DST_CAVES + "/modoverrides.lua"
}

func GetConfig() vo.GameConfigVO {
	gameConfig := vo.NewGameConfigVO()

	gameConfig.Token = getClusterToken()
	GetClusterIni(gameConfig)
	gameConfig.MasterMapData = getMasteLeveldataoverride()
	gameConfig.CavesMapData = getCavesLeveldataoverride()
	gameConfig.ModData = getModoverrides()

	return *gameConfig
}

func getClusterToken() string {
	token, err := fileUtils.ReadFile(constant.GET_CLUSTER_TOKEN_PATH())
	if err != nil {
		panic("read cluster_token.txt file error: " + err.Error())
	}

	return token
}

func GetClusterIni(gameconfig *vo.GameConfigVO) {
	cluster_ini, err := fileUtils.ReadLnFile(constant.GET_CLUSTER_INI_PATH())
	if err != nil {
		panic("read cluster.ini file error: " + err.Error())
	}
	for _, value := range cluster_ini {
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
					gameconfig.Pvp = b
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

func getMasteLeveldataoverride() string {
	level, err := fileUtils.ReadFile(constant.GET_MASTER_LEVELDATAOVERRIDE_PATH())
	if err != nil {
		panic("read Master/leveldataoverride.lua file error: " + err.Error())
	}
	return level
}

func getCavesLeveldataoverride() string {
	level, err := fileUtils.ReadFile(constant.GET_CAVES_LEVELDATAOVERRIDE_PATH())
	if err != nil {
		panic("read Caves/leveldataoverride.lua file error: " + err.Error())
	}
	return level
}

func getModoverrides() string {
	level, err := fileUtils.ReadFile(constant.GET_MASTER_MOD_PATH())
	if err != nil {
		panic("read Master/modoverrides.lua file error: " + err.Error())
	}
	return level
}

func SaveConfig(gameConfigVo vo.GameConfigVO) {

	//创建存档目录
	createMyDediServerDir()
	//创建房间配置
	createClusterIni(gameConfigVo)
	//创建token配置
	createClusterToken(strings.TrimSpace(gameConfigVo.Token))
	//创建地面和洞穴的ini配置文件
	createMasterServerIni()
	createCavesServerIni()
	//创建地面世界设置
	createMasteLeveldataoverride(gameConfigVo.MasterMapData)
	//创建洞穴世界设置
	createCavesLeveldataoverride(gameConfigVo.CavesMapData)
	//创建mod设置
	createModoverrides(gameConfigVo.ModData)
	//TODO dedicated_server_mods_setup
	otype := gameConfigVo.Otype
	if otype == START_NEW_GAME {
		DeleteGameRecord()
		StartGame(optype.START_GAME)
	} else if otype == SAVE_RESTART {
		StartGame(optype.START_GAME)
	}
}

func createMyDediServerDir() {
	myDediServerPath := constant.HOME_PATH + constant.DST_USER_GAME_CONFG_PATH
	log.Println("生成 myDediServer 目录：" + myDediServerPath)
	fileUtils.CreateDir(myDediServerPath)
}

func createClusterIni(gameConfigVo vo.GameConfigVO) {

	log.Println("生成游戏配置文件 cluster.ini文件: ", constant.GET_CLUSTER_INI_PATH())

	// cluster_ini := ""
	// cluster_ini += "[GAMEPLAY]\n"
	// cluster_ini += "game_mode = " + gameConfigVo.GameMode + "\n"
	// cluster_ini += "max_players = " + strconv.Itoa(int(gameConfigVo.MaxPlayers)) + "\n"
	// cluster_ini += "pvp = " + strconv.FormatBool(gameConfigVo.Pvp) + "\n"
	// cluster_ini += "pause_when_empty = " + strconv.FormatBool(gameConfigVo.PauseNobody) + "\n"
	// cluster_ini += "\n"
	// cluster_ini += "\n"
	// cluster_ini += "[NETWORK]\n"
	// cluster_ini += "lan_only_cluster = false\n"
	// cluster_ini += "cluster_intention = " + gameConfigVo.ClusterIntention + "\n"
	// password := gameConfigVo.ClusterPassword
	// if password != "" {
	// 	password = strings.TrimSpace(password)
	// 	cluster_ini += "cluster_password = " + password + "\n"
	// } else {
	// 	cluster_ini += "cluster_password = \n"
	// }
	// cluster_ini += "cluster_description =  " + gameConfigVo.ClusterDescription + " \n"
	// cluster_ini += "cluster_name =  " + gameConfigVo.ClusterName + " \n"
	// cluster_ini += "offline_cluster = false \n"

	// cluster_ini += "cluster_language =  zh\n"
	// cluster_ini += "\n"
	// cluster_ini += "[MISC]\n"
	// cluster_ini += "console_enabled = true\n"
	// cluster_ini += "max_snapshots = 6 \n"
	// cluster_ini += "\n"
	// cluster_ini += "\n"
	// cluster_ini += "[SHARD]\n"
	// cluster_ini += "shard_enabled = true\n"
	// cluster_ini += "bind_ip = 127.0.0.1\n"
	// cluster_ini += "master_ip = 127.0.0.1\n"
	// cluster_ini += "master_port = 10888\n"
	// cluster_ini += "cluster_key = defaultPass\n"

	cluster_ini := pareseTemplate(cluster_init_template, gameConfigVo)
	fileUtils.WriterTXT(constant.GET_CLUSTER_INI_PATH(), cluster_ini)
}

func createClusterToken(token string) {
	log.Println("生成cluster_token.txt 文件 ", constant.GET_CLUSTER_TOKEN_PATH())
	fileUtils.WriterTXT(constant.GET_CLUSTER_TOKEN_PATH(), token)
}

func createMasterServerIni() {

	fileUtils.CreateDir(constant.GET_MASTER_DIR_PATH())
	log.Println("生成 Master 目录: " + constant.GET_MASTER_DIR_PATH())

	log.Println("生成世界 Master server.ini文件: ", constant.GET_MASTER_DIR_SERVER_INI_PATH())

	// server_ini := ""
	// server_ini += "[NETWORK] \n"
	// server_ini += "server_port = " + "10999" + " \n"
	// server_ini += "\n"
	// server_ini += "\n"
	// server_ini += "[SHARD] \n"
	// server_ini += "is_master = true \n"
	// server_ini += "name = Master \n"
	// server_ini += "id = 10000 \n"
	// server_ini += "\n"
	// server_ini += "\n"
	// server_ini += "[ACCOUNT] \n"
	// server_ini += "encode_user_path = true"

	server_ini := pareseTemplate(master_server_init_template, nil)
	fileUtils.WriterTXT(constant.GET_MASTER_DIR_SERVER_INI_PATH(), server_ini)
}

func createCavesServerIni() {

	//创建洞穴设置的文件夹
	fileUtils.CreateDir(constant.GET_CAVE_DIR_PATH())
	log.Println("生成 Caves 目录: " + constant.GET_CAVE_DIR_PATH())

	log.Println("生成洞穴 Caves server.ini文件: ", constant.GET_CAVES_DIR_SERVER_INI_PATH())

	// caves_ini := ""
	// caves_ini += "[NETWORK] \n"
	// caves_ini += "server_port = 10998 \n"
	// caves_ini += "\n"
	// caves_ini += "\n"
	// caves_ini += "[SHARD]\n"
	// caves_ini += "is_master = false\n"
	// caves_ini += "name = Caves\n"
	// caves_ini += "id = 10010\n"
	// caves_ini += "\n"
	// caves_ini += "\n"
	// caves_ini += "[ACCOUNT]\n"
	// caves_ini += "encode_user_path = true\n"
	// caves_ini += "\n"
	// caves_ini += "\n"
	// caves_ini += "[STEAM]\n"
	// caves_ini += "authentication_port = 8766\n"
	// caves_ini += "master_server_port = 27016\n"

	caves_ini := pareseTemplate(caves_server_init_template, nil)
	fileUtils.WriterTXT(constant.GET_CAVES_DIR_SERVER_INI_PATH(), caves_ini)
}

func pareseTemplate(tempaltePath string, data interface{}) string {
	tmpl, err := template.ParseFiles(tempaltePath)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	tmpl.Execute(buf, data)
	fmt.Println("解析文本模板")
	fmt.Printf("buf.String():\n%v\n", buf.String())
	return buf.String()
}

func createMasteLeveldataoverride(mapConfig string) {

	log.Println("生成master_leveldataoverride.txt 文件 ", constant.GET_MASTER_LEVELDATAOVERRIDE_PATH())
	if mapConfig != "" {
		fileUtils.WriterTXT(constant.GET_MASTER_LEVELDATAOVERRIDE_PATH(), mapConfig)
	} else {
		//置空
		fileUtils.WriterTXT(constant.GET_MASTER_LEVELDATAOVERRIDE_PATH(), "")
	}
}
func createCavesLeveldataoverride(mapConfig string) {

	log.Println("生成caves_leveldataoverride.lua 文件 ", constant.GET_CAVES_LEVELDATAOVERRIDE_PATH())
	if mapConfig != "" {
		fileUtils.WriterTXT(constant.GET_CAVES_LEVELDATAOVERRIDE_PATH(), mapConfig)
	} else {
		//置空
		fileUtils.WriterTXT(constant.GET_CAVES_LEVELDATAOVERRIDE_PATH(), "")
	}
}
func createModoverrides(modConfig string) {

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
