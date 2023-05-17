package service

import (
	"dst-admin-go/constant"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo/cluster"
	"log"
	"path"
	"strings"
	"sync"

	"github.com/go-ini/ini"
)

func ReadClusterIniFile() *cluster.Cluster {
	cluster := cluster.NewCluster()
	cluster_ini, err := fileUtils.ReadLnFile(constant.GET_CLUSTER_INI_PATH())
	if err != nil {
		panic("read cluster.ini file error: " + err.Error())
	}

	// 加载 INI 文件
	cfg, err := ini.Load(cluster_ini)
	if err != nil {
		log.Panicln("Failed to load INI file:", err)
	}

	// [GAMEPLAY]
	GAMEPLAY := cfg.Section("GAMEPLAY")

	cluster.GameMode = GAMEPLAY.Key("game_mode").String()
	cluster.MaxPlayers = GAMEPLAY.Key("max_players").MustUint(8)
	cluster.Pvp = GAMEPLAY.Key("pvp").MustBool(false)
	cluster.PauseWhenNobody = GAMEPLAY.Key("pause_when_empty").MustBool(true)
	cluster.VoteEnabled = GAMEPLAY.Key("vote_enabled").MustBool(true)
	cluster.VoteKickEnabled = GAMEPLAY.Key("vote_kick_enabled").MustBool(true)

	// [NETWORK]
	NETWORK := cfg.Section("NETWORK")

	cluster.LanOnlyCluster = NETWORK.Key("lan_only_cluster").MustBool(false)
	cluster.ClusterIntention = NETWORK.Key("cluster_intention").String()
	cluster.ClusterPassword = NETWORK.Key("cluster_password").String()
	cluster.ClusterDescription = NETWORK.Key("cluster_description").String()
	cluster.ClusterName = NETWORK.Key("cluster_name").String()
	cluster.OfflineCluster = NETWORK.Key("offline_cluster").MustBool(false)
	cluster.ClusterLanguage = NETWORK.Key("cluster_language").String()
	cluster.WhitelistSlots = NETWORK.Key("whitelist_slots").MustUint(0)
	cluster.TickRate = NETWORK.Key("tick_rate").MustUint(15)

	// [MISC]
	MISC := cfg.Section("MISC")

	cluster.ConsoleEnabled = MISC.Key("console_enabled").MustBool(true)
	cluster.MaxSnapshots = MISC.Key("max_snapshots").MustUint(6)

	// [SHARD]
	SHARD := cfg.Section("SHARD")

	cluster.ShardEnabled = SHARD.Key("shard_enabled").MustBool(true)
	cluster.BindIp = SHARD.Key("bind_ip").MustString("127.0.0.1")
	cluster.MasterIp = SHARD.Key("master_ip").MustString("127.0.0.1")
	cluster.MasterPort = SHARD.Key("master_ip").MustUint(108888)
	cluster.ClusterKey = SHARD.Key("cluster_key").String()

	// [STEAM]
	STEAM := cfg.Section("STEAM")

	cluster.SteamGroupOnly = STEAM.Key("steam_group_only").MustBool(false)
	cluster.SteamGroupId = STEAM.Key("steam_group_id").MustUint(0)
	cluster.SteamGroupAdmins = STEAM.Key("steam_group_admins").MustBool(false)

	return cluster
}

func ReadClusterTokenFile() string {
	token, err := fileUtils.ReadFile(constant.GET_CLUSTER_TOKEN_PATH())
	if err != nil {
		panic("read cluster_token.txt file error: " + err.Error())
	}
	return token
}

func ReadAdminlistFile() (str []string) {
	path := constant.GET_DST_ADMIN_LIST_PATH()
	if !fileUtils.Exists(path) {
		log.Panicln("路径不存在", path)
	}
	str, err := fileUtils.ReadLnFile(path)
	log.Println("str:", str)
	if err != nil {
		panic("read dst adminlist.txt error: \n" + err.Error())
	}
	return
}

func ReadBlocklistFile() (str []string) {
	path := constant.GET_DST_BLOCKLIST_PATH()
	if !fileUtils.Exists(path) {
		log.Println("路径不存在", path)
	}
	str, err := fileUtils.ReadLnFile(path)
	log.Println("str:", str)
	if err != nil {
		panic("read dst blocklist.txt error: \n" + err.Error())
	}
	return
}

func ReadLeveldataoverrideFile(filepath string) string {
	leveldataoverride, err := fileUtils.ReadFile(filepath)
	if err != nil {
		panic("read leveldataoverride.lua file error: " + err.Error())
	}
	return leveldataoverride
}

func ReadModoverridesFile(filepath string) string {
	modoverrides, err := fileUtils.ReadFile(filepath)
	if err != nil {
		panic("read modoverrides.lua file error: " + err.Error())
	}
	return modoverrides
}

func ReadServerIniFile(filepath string, isMaster bool) *cluster.ServerIni {

	var server_port_default uint = 10998
	id_default := 10010

	if isMaster {
		server_port_default = 10999
		id_default = 10000
	}

	serverIni := cluster.NewCavesServerIni()
	// 加载 INI 文件
	cfg, err := ini.Load(filepath)
	if err != nil {
		log.Panicln("Failed to load INI file:", err)
	}

	// [NETWORK]
	NETWORK := cfg.Section("NETWORK")

	serverIni.ServerPort = NETWORK.Key("server_port").MustUint(server_port_default)

	// [SHARD]
	SHARD := cfg.Section("SHARD")

	serverIni.IsMaster = SHARD.Key("is_master").MustBool(isMaster)
	serverIni.Name = SHARD.Key("name").String()
	serverIni.Id = SHARD.Key("id").MustUint(uint(id_default))

	// [ACCOUNT]
	ACCOUNT := cfg.Section("ACCOUNT")
	serverIni.EncodeUserPath = ACCOUNT.Key("encode_user_path").MustBool(true)

	// [STEAM]
	STEAM := cfg.Section("STEAM")

	serverIni.AuthenticationPort = STEAM.Key("authentication_port").MustUint()
	serverIni.Master_serverPort = STEAM.Key("master_server_port").MustUint()

	return serverIni
}

func isMaster(filePath string) bool {
	return strings.Contains(filePath, "Master") || strings.Contains(filePath, "master")
}

func GetmultiLevelWorldConfig() *cluster.MultiLevelWorldConfig {
	multiLevelWorldConfig := &cluster.MultiLevelWorldConfig{}

	cluster_path := constant.GET_CLUSTER_INI_PATH()
	worldDirs, err := fileUtils.FindWorldDirs(cluster_path)
	if err != nil {
		log.Panicln("路径查找失败", err)
	}
	size := len(worldDirs)
	worlds := make([]cluster.World, size)
	var wg sync.WaitGroup
	wg.Add(size)
	for i, worldPath := range worldDirs {
		//查找
		go func(word cluster.World, worldPath string) {
			word.Leveldataoverride = ReadLeveldataoverrideFile(path.Join(worldPath, "leveldataoverride.lua"))
			word.Modoverrides = ReadModoverridesFile(path.Join(worldPath, "modoverrides.lua"))
			word.IsMaster = isMaster(worldPath)
			word.ServerIni = ReadServerIniFile(path.Join(worldPath, "server_ini.lua"), isMaster(worldPath))
			wg.Done()
		}(worlds[i], worldPath)
	}
	wg.Wait()
	multiLevelWorldConfig.Worlds = worlds
	return multiLevelWorldConfig
}
