package service

import (
	"dst-admin-go/constant/dst"
	"dst-admin-go/utils/dstUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo/level"
	"github.com/go-ini/ini"
	"log"
	"strings"
)

type HomeService struct{}

const (
	ClusterIniTemplate      = "./static/template/cluster2.ini"
	MasterServerIniTemplate = "./static/template/master_server.ini"
	CavesServerIniTemplate  = "./static/template/caves_server.ini"
	ServerIniTemplate       = "./static/template/server.ini"
)

func (c *HomeService) GetClusterIni(clusterName string) *level.ClusterIni {
	newCluster := level.NewClusterIni()
	// 加载 INI 文件
	clusterIniPath := dst.GetClusterIniPath(clusterName)
	if !fileUtils.Exists(clusterIniPath) {
		err := fileUtils.CreateFileIfNotExists(clusterIniPath)
		if err != nil {
			return nil
		}
		return newCluster
	}
	cfg, err := ini.Load(clusterIniPath)
	if err != nil {
		log.Panicln("Failed to load INI file:", err)
	}

	// [GAMEPLAY]
	GAMEPLAY := cfg.Section("GAMEPLAY")

	newCluster.GameMode = GAMEPLAY.Key("game_mode").String()
	newCluster.MaxPlayers = GAMEPLAY.Key("max_players").MustUint(8)
	newCluster.Pvp = GAMEPLAY.Key("pvp").MustBool(false)
	newCluster.PauseWhenNobody = GAMEPLAY.Key("pause_when_empty").MustBool(true)
	newCluster.VoteEnabled = GAMEPLAY.Key("vote_enabled").MustBool(true)
	newCluster.VoteKickEnabled = GAMEPLAY.Key("vote_kick_enabled").MustBool(true)

	// [NETWORK]
	NETWORK := cfg.Section("NETWORK")

	newCluster.LanOnlyCluster = NETWORK.Key("lan_only_cluster").MustBool(false)
	newCluster.ClusterIntention = NETWORK.Key("cluster_intention").String()
	newCluster.ClusterPassword = NETWORK.Key("cluster_password").String()
	newCluster.ClusterDescription = NETWORK.Key("cluster_description").String()
	newCluster.ClusterName = NETWORK.Key("cluster_name").String()
	newCluster.OfflineCluster = NETWORK.Key("offline_cluster").MustBool(false)
	newCluster.ClusterLanguage = NETWORK.Key("cluster_language").String()
	newCluster.WhitelistSlots = NETWORK.Key("whitelist_slots").MustUint(0)
	newCluster.TickRate = NETWORK.Key("tick_rate").MustUint(15)

	// [MISC]
	MISC := cfg.Section("MISC")

	newCluster.ConsoleEnabled = MISC.Key("console_enabled").MustBool(true)
	newCluster.MaxSnapshots = MISC.Key("max_snapshots").MustUint(6)

	// [SHARD]
	SHARD := cfg.Section("SHARD")

	newCluster.ShardEnabled = SHARD.Key("shard_enabled").MustBool(true)
	newCluster.BindIp = SHARD.Key("bind_ip").MustString("127.0.0.1")
	newCluster.MasterIp = SHARD.Key("master_ip").MustString("127.0.0.1")
	newCluster.MasterPort = SHARD.Key("master_port").MustUint(10888)
	newCluster.ClusterKey = SHARD.Key("cluster_key").String()

	// [STEAM]
	STEAM := cfg.Section("STEAM")

	newCluster.SteamGroupOnly = STEAM.Key("steam_group_only").MustBool(false)
	newCluster.SteamGroupId = STEAM.Key("steam_group_id").MustString("")
	newCluster.SteamGroupAdmins = STEAM.Key("steam_group_admins").MustString("")

	return newCluster
}

func (c *HomeService) GetClusterToken(clusterName string) string {
	clusterTokenPath := dst.GetClusterTokenPath(clusterName)
	token, err := fileUtils.ReadFile(clusterTokenPath)
	if err != nil {
		return ""
	}
	return token
}

func (c *HomeService) GetAdminlist(clusterName string) (str []string) {

	adminListPath := dst.GetAdminlistPath(clusterName)
	fileUtils.CreateFileIfNotExists(adminListPath)
	str, err := fileUtils.ReadLnFile(adminListPath)
	if err != nil {
		return []string{}
	}
	return
}

func (c *HomeService) GetBlocklist(clusterName string) (str []string) {
	blocklistPath := dst.GetBlocklistPath(clusterName)
	fileUtils.CreateFileIfNotExists(blocklistPath)
	str, err := fileUtils.ReadLnFile(blocklistPath)
	if err != nil {
		return []string{}
	}
	return
}

func (c *HomeService) GetLeveldataoverride(filepath string) string {
	if !fileUtils.Exists(filepath) {
		fileUtils.CreateFileIfNotExists(filepath)
		return "return {}"
	}

	leveldataoverride, err := fileUtils.ReadFile(filepath)
	if err != nil {
		return "return {}"
	}
	if leveldataoverride == "" {
		return "return {}"
	}
	return leveldataoverride
}

func (c *HomeService) GetModoverrides(filepath string) string {
	if !fileUtils.Exists(filepath) {
		fileUtils.CreateFileIfNotExists(filepath)
		return "return {}"
	}
	modoverrides, err := fileUtils.ReadFile(filepath)
	if err != nil {
		return "return {}"
	}
	if modoverrides == "" {
		return "return {}"
	}
	return modoverrides
}

func (c *HomeService) GetServerIni(filepath string, isMaster bool) *level.ServerIni {
	fileUtils.CreateFileIfNotExists(filepath)
	var serverPortDefault uint = 10998
	idDefault := 10010

	if isMaster {
		serverPortDefault = 10999
		idDefault = 10000
	}

	serverIni := level.NewCavesServerIni()
	// 加载 INI 文件
	cfg, err := ini.Load(filepath)
	if err != nil {
		return serverIni
	}

	// [NETWORK]
	NETWORK := cfg.Section("NETWORK")

	serverIni.ServerPort = NETWORK.Key("server_port").MustUint(serverPortDefault)

	// [SHARD]
	SHARD := cfg.Section("SHARD")

	serverIni.IsMaster = SHARD.Key("is_master").MustBool(isMaster)
	serverIni.Name = SHARD.Key("name").String()
	serverIni.Id = SHARD.Key("id").MustUint(uint(idDefault))

	// [ACCOUNT]
	ACCOUNT := cfg.Section("ACCOUNT")
	serverIni.EncodeUserPath = ACCOUNT.Key("encode_user_path").MustBool(true)

	// [STEAM]
	STEAM := cfg.Section("STEAM")

	serverIni.AuthenticationPort = STEAM.Key("authentication_port").MustUint(8766)
	serverIni.MasterServerPort = STEAM.Key("master_server_port").MustUint(27016)

	return serverIni
}

func (c *HomeService) isMaster(filePath string) bool {
	return strings.Contains(filePath, "Master") || strings.Contains(filePath, "master")
}

func (c *HomeService) SaveClusterToken(clusterName, token string) {
	clusterTokenPath := dst.GetClusterTokenPath(clusterName)
	fileUtils.WriterTXT(clusterTokenPath, token)
}

func (c *HomeService) SaveClusterIni(clusterName string, cluster *level.ClusterIni) {
	clusterIniPath := dst.GetClusterIniPath(clusterName)
	fileUtils.WriterTXT(clusterIniPath, dstUtils.ParseTemplate(ClusterIniTemplate, cluster))
}

func (c *HomeService) SaveAdminlist(clusterName string, str []string) {
	adminlistPath := dst.GetAdminlistPath(clusterName)
	fileUtils.CreateFileIfNotExists(adminlistPath)
	fileUtils.WriterLnFile(adminlistPath, str)
}

func (c *HomeService) SaveBlocklist(clusterName string, str []string) {
	blocklistPath := dst.GetBlocklistPath(clusterName)
	fileUtils.CreateFileIfNotExists(blocklistPath)
	fileUtils.WriterLnFile(blocklistPath, str)
}

func (c *HomeService) GetMasterWorld(clusterName string) *level.World {
	master := level.World{}
	master.WorldName = "Master"
	master.IsMaster = true

	master.Leveldataoverride = c.GetLeveldataoverride(dst.GetMasterLeveldataoverridePath(clusterName))
	master.Modoverrides = c.GetModoverrides(dst.GetMasterModoverridesPath(clusterName))
	master.ServerIni = c.GetServerIni(dst.GetMasterServerIniPath(clusterName), true)

	return &master
}

func (c *HomeService) GetCavesWorld(clusterName string) *level.World {
	caves := level.World{}

	caves.WorldName = "Caves"
	caves.IsMaster = false

	caves.Leveldataoverride = c.GetLeveldataoverride(dst.GetCavesLeveldataoverridePath(clusterName))
	caves.Modoverrides = c.GetModoverrides(dst.GetCavesModoverridesPath(clusterName))
	caves.ServerIni = c.GetServerIni(dst.GetCavesServerIniPath(clusterName), false)

	return &caves
}

func (c *HomeService) SaveMasterWorld(clusterName string, world *level.World) {

	lPath := dst.GetMasterLeveldataoverridePath(clusterName)
	mPath := dst.GetMasterModoverridesPath(clusterName)
	sPath := dst.GetMasterServerIniPath(clusterName)

	fileUtils.CreateFileIfNotExists(lPath)
	fileUtils.CreateFileIfNotExists(mPath)
	fileUtils.CreateFileIfNotExists(sPath)

	fileUtils.WriterTXT(lPath, world.Leveldataoverride)
	fileUtils.WriterTXT(mPath, world.Modoverrides)

	serverBuf := dstUtils.ParseTemplate(MasterServerIniTemplate, world.ServerIni)

	fileUtils.WriterTXT(sPath, serverBuf)
}

func (c *HomeService) SaveCavesWorld(clusterName string, world *level.World) {

	lPath := dst.GetCavesLeveldataoverridePath(clusterName)
	mPath := dst.GetCavesModoverridesPath(clusterName)
	sPath := dst.GetCavesServerIniPath(clusterName)

	fileUtils.CreateFileIfNotExists(lPath)
	fileUtils.CreateFileIfNotExists(mPath)
	fileUtils.CreateFileIfNotExists(sPath)

	fileUtils.WriterTXT(lPath, world.Leveldataoverride)
	fileUtils.WriterTXT(mPath, world.Modoverrides)

	serverBuf := dstUtils.ParseTemplate(CavesServerIniTemplate, world.ServerIni)

	fileUtils.WriterTXT(sPath, serverBuf)
}

func (c *HomeService) GetLevel(clusterName string, levelName string) *level.World {
	level := level.World{}

	level.WorldName = levelName
	level.IsMaster = false

	level.Leveldataoverride = c.GetLeveldataoverride(dst.GetLevelLeveldataoverridePath(clusterName, levelName))
	level.Modoverrides = c.GetModoverrides(dst.GetLevelModoverridesPath(clusterName, levelName))
	level.ServerIni = c.GetServerIni(dst.GetLevelServerIniPath(clusterName, levelName), false)

	return &level
}

func (c *HomeService) SaveLevel(clusterName string, levelName string, world *level.World) {

	lPath := dst.GetLevelLeveldataoverridePath(clusterName, levelName)
	mPath := dst.GetLevelModoverridesPath(clusterName, levelName)
	sPath := dst.GetLevelServerIniPath(clusterName, levelName)

	fileUtils.CreateFileIfNotExists(lPath)
	fileUtils.CreateFileIfNotExists(mPath)
	fileUtils.CreateFileIfNotExists(sPath)

	fileUtils.WriterTXT(lPath, world.Leveldataoverride)
	fileUtils.WriterTXT(mPath, world.Modoverrides)

	serverBuf := dstUtils.ParseTemplate(ServerIniTemplate, world.ServerIni)

	fileUtils.WriterTXT(sPath, serverBuf)

	// 写入 dedicated_server_mods_setup.lua
	dstUtils.DedicatedServerModsSetup2(clusterName, world.Modoverrides)

}
