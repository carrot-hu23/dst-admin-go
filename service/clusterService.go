package service

import (
	"dst-admin-go/constant/dst"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo/cluster"
	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
	"log"
	"strings"
	"sync"
)

type ClusterService struct {
	DstHelper
}

const (
	CLUSTER_INI_TEMPLATE       = "./static/template/cluster2.ini"
	MASTER_SERVER_INI_TEMPLATE = "./static/template/master_server.ini"
	CAVES_SERVER_INI_TEMPLATE  = "./static/template/caves_server.ini"
)

func (c *ClusterService) ReadClusterIniFile(clusterName string) *cluster.ClusterIni {
	newCluster := cluster.NewCluster()
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
	newCluster.SteamGroupId = STEAM.Key("steam_group_id").MustUint(0)
	newCluster.SteamGroupAdmins = STEAM.Key("steam_group_admins").MustString("")

	return newCluster
}

func (c *ClusterService) ReadClusterTokenFile(clusterName string) string {
	clusterTokenPath := dst.GetClusterTokenPath(clusterName)
	token, err := fileUtils.ReadFile(clusterTokenPath)
	if err != nil {
		panic("read cluster_token.txt file error: " + err.Error())
	}
	return token
}

func (c *ClusterService) ReadAdminlistFile(clusterName string) (str []string) {

	adminListPath := dst.GetAdminlistPath(clusterName)
	fileUtils.CreateFileIfNotExists(adminListPath)
	str, err := fileUtils.ReadLnFile(adminListPath)
	log.Println("str:", str)
	if err != nil {
		panic("read dst adminlist.txt error: \n" + err.Error())
	}
	return
}

func (c *ClusterService) ReadBlocklistFile(clusterName string) (str []string) {
	blocklistPath := dst.GetBlocklistPath(clusterName)
	fileUtils.CreateFileIfNotExists(blocklistPath)
	str, err := fileUtils.ReadLnFile(blocklistPath)
	log.Println("str:", str)
	if err != nil {
		panic("read dst blocklist.txt error: \n" + err.Error())
	}
	return
}

func (c *ClusterService) ReadLeveldataoverrideFile(filepath string) string {
	if !fileUtils.Exists(filepath) {
		fileUtils.CreateFileIfNotExists(filepath)
		return "return {}"
	}

	leveldataoverride, err := fileUtils.ReadFile(filepath)
	if err != nil {
		panic("read leveldataoverride.lua file error: " + err.Error())
	}
	return leveldataoverride
}

func (c *ClusterService) ReadModoverridesFile(filepath string) string {
	if !fileUtils.Exists(filepath) {
		fileUtils.CreateFileIfNotExists(filepath)
		return "return {}"
	}
	modoverrides, err := fileUtils.ReadFile(filepath)
	if err != nil {
		panic("read modoverrides.lua file error: " + err.Error())
	}
	return modoverrides
}

func (c *ClusterService) ReadServerIniFile(filepath string, isMaster bool) *cluster.ServerIni {
	fileUtils.CreateFileIfNotExists(filepath)
	var serverPortDefault uint = 10998
	idDefault := 10010

	if isMaster {
		serverPortDefault = 10999
		idDefault = 10000
	}

	serverIni := cluster.NewCavesServerIni()
	// 加载 INI 文件
	cfg, err := ini.Load(filepath)
	if err != nil {
		log.Panicln("Failed to load INI file:", err)
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

	serverIni.AuthenticationPort = STEAM.Key("authentication_port").String()
	serverIni.MasterServerPort = STEAM.Key("master_server_port").String()

	return serverIni
}

func (c *ClusterService) isMaster(filePath string) bool {
	return strings.Contains(filePath, "Master") || strings.Contains(filePath, "master")
}

func (c *ClusterService) SaveClusterToken(clusterName, token string) {
	clusterTokenPath := dst.GetClusterTokenPath(clusterName)
	fileUtils.WriterTXT(clusterTokenPath, token)
}

func (c *ClusterService) SaveClusterIni(clusterName string, cluster *cluster.ClusterIni) {
	clusterIniPath := dst.GetClusterIniPath(clusterName)
	fileUtils.WriterTXT(clusterIniPath, c.ParseTemplate(CLUSTER_INI_TEMPLATE, cluster))
}

func (c *ClusterService) SaveAdminlist(clusterName string, str []string) {
	adminlistPath := dst.GetAdminlistPath(clusterName)
	fileUtils.CreateFileIfNotExists(adminlistPath)
	fileUtils.WriterLnFile(adminlistPath, str)
}

func (c *ClusterService) SaveBlocklist(clusterName string, str []string) {
	blocklistPath := dst.GetBlocklistPath(clusterName)
	fileUtils.CreateFileIfNotExists(blocklistPath)
	fileUtils.WriterLnFile(blocklistPath, str)
}

func (c *ClusterService) ReadMaster(clusterName string) *cluster.World {
	master := cluster.World{}

	master.WorldName = "Master"
	master.IsMaster = true

	master.Leveldataoverride = c.ReadLeveldataoverrideFile(dst.GetMasterLeveldataoverridePath(clusterName))
	master.Modoverrides = c.ReadModoverridesFile(dst.GetMasterModoverridesPath(clusterName))
	master.ServerIni = c.ReadServerIniFile(dst.GetMasterServerIniPath(clusterName), true)

	return &master
}

func (c *ClusterService) ReadCaves(clusterName string) *cluster.World {
	caves := cluster.World{}

	caves.WorldName = "Caves"
	caves.IsMaster = false

	caves.Leveldataoverride = c.ReadLeveldataoverrideFile(dst.GetCavesLeveldataoverridePath(clusterName))
	caves.Modoverrides = c.ReadModoverridesFile(dst.GetCavesModoverridesPath(clusterName))
	caves.ServerIni = c.ReadServerIniFile(dst.GetCavesServerIniPath(clusterName), false)

	return &caves
}

func (c *ClusterService) SaveMaster(clusterName string, world *cluster.World) {

	lPath := dst.GetMasterLeveldataoverridePath(clusterName)
	mPath := dst.GetMasterModoverridesPath(clusterName)
	sPath := dst.GetMasterServerIniPath(clusterName)

	fileUtils.CreateFileIfNotExists(lPath)
	fileUtils.CreateFileIfNotExists(mPath)
	fileUtils.CreateFileIfNotExists(sPath)

	fileUtils.WriterTXT(lPath, world.Leveldataoverride)
	fileUtils.WriterTXT(mPath, world.Modoverrides)

	serverBuf := c.ParseTemplate(MASTER_SERVER_INI_TEMPLATE, world.ServerIni)

	fileUtils.WriterTXT(sPath, serverBuf)
}

func (c *ClusterService) SaveCaves(clusterName string, world *cluster.World) {

	lPath := dst.GetCavesLeveldataoverridePath(clusterName)
	mPath := dst.GetCavesModoverridesPath(clusterName)
	sPath := dst.GetCavesServerIniPath(clusterName)

	fileUtils.CreateFileIfNotExists(lPath)
	fileUtils.CreateFileIfNotExists(mPath)
	fileUtils.CreateFileIfNotExists(sPath)

	fileUtils.WriterTXT(lPath, world.Leveldataoverride)
	fileUtils.WriterTXT(mPath, world.Modoverrides)

	serverBuf := c.ParseTemplate(CAVES_SERVER_INI_TEMPLATE, world.ServerIni)

	fileUtils.WriterTXT(sPath, serverBuf)
}

func (c *ClusterService) GetGameConfig(ctx *gin.Context) *cluster.GameConfig {
	gameConfig := cluster.GameConfig{}
	var wg sync.WaitGroup
	wg.Add(6)
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	go func() {
		gameConfig.ClusterToken = c.ReadClusterTokenFile(clusterName)
		wg.Done()
	}()
	go func() {
		gameConfig.ClusterIni = c.ReadClusterIniFile(clusterName)
		wg.Done()
	}()
	go func() {
		gameConfig.Adminlist = c.ReadAdminlistFile(clusterName)
		wg.Done()
	}()
	go func() {
		gameConfig.Blocklist = c.ReadBlocklistFile(clusterName)
		wg.Done()
	}()
	go func() {
		gameConfig.Master = c.ReadMaster(clusterName)
		wg.Done()
	}()
	go func() {
		gameConfig.Caves = c.ReadCaves(clusterName)
		wg.Done()
	}()
	wg.Wait()
	return &gameConfig
}

func (c *ClusterService) SaveGameConfig(ctx *gin.Context, gameConfig *cluster.GameConfig) {

	var wg sync.WaitGroup
	wg.Add(6)
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	go func() {
		c.SaveClusterToken(clusterName, gameConfig.ClusterToken)
		wg.Done()
	}()

	go func() {
		c.SaveClusterIni(clusterName, gameConfig.ClusterIni)
		wg.Done()
	}()

	go func() {
		// SaveAdminlist(gameConfig.Adminlist)
		wg.Done()
	}()

	go func() {
		// SaveBlocklist(gameConfig.Blocklist)
		wg.Done()
	}()

	go func() {
		c.SaveMaster(clusterName, gameConfig.Master)
		c.DedicatedServerModsSetup(clusterName, gameConfig.Master.Modoverrides)
		wg.Done()
	}()

	go func() {
		c.SaveCaves(clusterName, gameConfig.Caves)
		wg.Done()
	}()

	wg.Wait()
}
