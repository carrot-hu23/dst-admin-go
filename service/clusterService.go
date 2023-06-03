package service

import (
	"dst-admin-go/constant"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo/cluster"
	"log"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/go-ini/ini"
)

type ClusterService struct {
	DstHelper
}

const (
	CLUSTER_INI_TEMPLATE       = "./static/template/cluster2.ini"
	MASTER_SERVER_INI_TEMPLATE = "./static/template/master_server.ini"
	CAVES_SERVER_INI_TEMPLATE  = "./static/template/caves_server.ini"
)

func (c *ClusterService) ReadClusterIniFile() *cluster.Cluster {
	newCluster := cluster.NewCluster()
	// 加载 INI 文件
	clusterIniPath := constant.GET_CLUSTER_INI_PATH()
	if !fileUtils.Exists(clusterIniPath) {
		fileUtils.CreateFileIfNotExists(clusterIniPath)
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

func (c *ClusterService) ReadClusterTokenFile() string {
	clusterTokenPath := constant.GET_CLUSTER_TOKEN_PATH()
	if !fileUtils.Exists(clusterTokenPath) {
		fileUtils.CreateFileIfNotExists(clusterTokenPath)
		return ""
	}

	token, err := fileUtils.ReadFile(clusterTokenPath)
	if err != nil {
		panic("read cluster_token.txt file error: " + err.Error())
	}
	return token
}

func (c *ClusterService) ReadAdminlistFile() (str []string) {
	adminListPath := constant.GET_DST_ADMIN_LIST_PATH()
	fileUtils.CreateFileIfNotExists(adminListPath)
	str, err := fileUtils.ReadLnFile(adminListPath)
	log.Println("str:", str)
	if err != nil {
		panic("read dst adminlist.txt error: \n" + err.Error())
	}
	return
}

func (c *ClusterService) ReadBlocklistFile() (str []string) {
	blocklistPath := constant.GET_DST_BLOCKLIST_PATH()
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

func (c *ClusterService) GetMultiLevelWorldConfig() *cluster.MultiLevelWorldConfig {
	multiLevelWorldConfig := &cluster.MultiLevelWorldConfig{}

	clusterPath := constant.GET_DST_USER_GAME_CONFG_PATH()
	worldDirs, err := fileUtils.FindWorldDirs(clusterPath)
	if err != nil {
		log.Panicln("路径查找失败", err)
	}
	size := len(worldDirs)
	worlds := make([]cluster.World, size)
	var wg sync.WaitGroup
	wg.Add(size + 4)
	for i, worldPath := range worldDirs {
		//查找
		go func(i int, worldPath string) {
			worlds[i].WorldName = filepath.Base(worldPath)
			worlds[i].Leveldataoverride = c.ReadLeveldataoverrideFile(path.Join(worldPath, "leveldataoverride.lua"))
			worlds[i].Modoverrides = c.ReadModoverridesFile(path.Join(worldPath, "modoverrides.lua"))
			worlds[i].IsMaster = c.isMaster(worldPath)
			worlds[i].ServerIni = c.ReadServerIniFile(path.Join(worldPath, "server.ini"), c.isMaster(worldPath))
			wg.Done()
		}(i, worldPath)
	}

	go func() {
		multiLevelWorldConfig.ClusterToken = c.ReadClusterTokenFile()
		wg.Done()
	}()
	go func() {
		multiLevelWorldConfig.Cluster = c.ReadClusterIniFile()
		wg.Done()
	}()
	go func() {
		multiLevelWorldConfig.Adminlist = c.ReadAdminlistFile()
		wg.Done()
	}()
	go func() {
		multiLevelWorldConfig.Blocklist = c.ReadBlocklistFile()
		wg.Done()
	}()

	multiLevelWorldConfig.Worlds = worlds
	wg.Wait()
	return multiLevelWorldConfig
}

func (c *ClusterService) SaveClusterToken(token string) {
	clusterTokenPath := constant.GET_CLUSTER_TOKEN_PATH()
	fileUtils.CreateFileIfNotExists(clusterTokenPath)
	fileUtils.WriterTXT(clusterTokenPath, token)
}

func (c *ClusterService) SaveClusterIni(cluster *cluster.Cluster) {
	clusterIniPath := constant.GET_CLUSTER_INI_PATH()
	fileUtils.CreateFileIfNotExists(clusterIniPath)
	fileUtils.WriterTXT(clusterIniPath, c.ParseTemplate(CLUSTER_INI_TEMPLATE, cluster))
}

func (c *ClusterService) SaveAdminlist(str []string) {
	adminlistPath := constant.GET_DST_ADMIN_LIST_PATH()
	fileUtils.CreateFileIfNotExists(adminlistPath)
	fileUtils.WriterLnFile(adminlistPath, str)
}

func (c *ClusterService) SaveBlocklist(str []string) {
	blocklistPath := constant.GET_DST_BLOCKLIST_PATH()
	fileUtils.CreateFileIfNotExists(blocklistPath)
	fileUtils.WriterLnFile(blocklistPath, str)
}

func (c *ClusterService) handleWorldDir(basePath string, world *cluster.World) {

	fileUtils.CreateDir(basePath)
	leveldataoverridePath := path.Join(basePath, "leveldataoverride.lua")
	modoverridesPath := path.Join(basePath, "modoverrides.lua")
	serverIniPath := path.Join(basePath, "server.ini")

	fileUtils.CreateFileIfNotExists(leveldataoverridePath)
	fileUtils.CreateFileIfNotExists(modoverridesPath)
	fileUtils.CreateFileIfNotExists(serverIniPath)

	fileUtils.WriterTXT(leveldataoverridePath, world.Leveldataoverride)
	fileUtils.WriterTXT(modoverridesPath, world.Modoverrides)

	// TODO 写入 constant.DST_MOD_SETTING_PATH

	serverIniBuf := ""
	if world.IsMaster {
		serverIniBuf = c.ParseTemplate(MASTER_SERVER_INI_TEMPLATE, world.ServerIni)
	} else {
		serverIniBuf = c.ParseTemplate(CAVES_SERVER_INI_TEMPLATE, world.ServerIni)
	}

	fileUtils.WriterTXT(serverIniPath, serverIniBuf)
}

func (c *ClusterService) SaveWorlds(worlds []cluster.World) {
	var wg sync.WaitGroup
	wg.Add(len(worlds))
	basePath := constant.GET_DST_USER_GAME_CONFG_PATH()
	for _, world := range worlds {
		go func(world *cluster.World) {
			// 判断文件是否存在，如果不存在则创建
			filePath := path.Join(basePath, world.WorldName)
			c.handleWorldDir(filePath, world)
			wg.Done()
		}(&world)
	}
	wg.Wait()
}

func (c *ClusterService) SaveMultiLevelWorldConfig(multiLevelWorldConfig *cluster.MultiLevelWorldConfig) {

	var wg sync.WaitGroup
	wg.Add(5)
	var size int32
	go func() {
		c.SaveClusterToken(multiLevelWorldConfig.ClusterToken)
		wg.Done()
		atomic.AddInt32(&size, 1)
	}()

	go func() {
		c.SaveClusterIni(multiLevelWorldConfig.Cluster)
		wg.Done()
		atomic.AddInt32(&size, 1)
	}()

	go func() {
		c.SaveAdminlist(multiLevelWorldConfig.Adminlist)
		wg.Done()
		atomic.AddInt32(&size, 1)
	}()
	go func() {
		c.SaveBlocklist(multiLevelWorldConfig.Blocklist)
		wg.Done()
		atomic.AddInt32(&size, 1)
	}()
	go func() {
		c.SaveWorlds(multiLevelWorldConfig.Worlds)
		wg.Done()
		atomic.AddInt32(&size, 1)
	}()

	wg.Wait()
	if size != 5 {
		// 设置失败
	}
}

func (c *ClusterService) ReadMaster() *cluster.World {
	master := cluster.World{}

	master.WorldName = "Master"
	master.IsMaster = true

	master.Leveldataoverride = c.ReadLeveldataoverrideFile(constant.GET_MASTER_LEVELDATAOVERRIDE_PATH())
	master.Modoverrides = c.ReadModoverridesFile(constant.GET_MASTER_MOD_PATH())
	master.ServerIni = c.ReadServerIniFile(constant.GET_MASTER_DIR_SERVER_INI_PATH(), true)

	return &master
}

func (c *ClusterService) ReadCaves() *cluster.World {
	caves := cluster.World{}

	caves.WorldName = "Caves"
	caves.IsMaster = false

	caves.Leveldataoverride = c.ReadLeveldataoverrideFile(constant.GET_CAVES_LEVELDATAOVERRIDE_PATH())
	caves.Modoverrides = c.ReadModoverridesFile(constant.GET_CAVES_MOD_PATH())
	caves.ServerIni = c.ReadServerIniFile(constant.GET_CAVES_DIR_SERVER_INI_PATH(), false)

	return &caves
}

func (c *ClusterService) SaveMaster(world *cluster.World) {

	lPath := constant.GET_MASTER_LEVELDATAOVERRIDE_PATH()
	mPath := constant.GET_MASTER_MOD_PATH()
	sPath := constant.GET_MASTER_DIR_SERVER_INI_PATH()

	fileUtils.CreateFileIfNotExists(lPath)
	fileUtils.CreateFileIfNotExists(mPath)
	fileUtils.CreateFileIfNotExists(sPath)

	fileUtils.WriterTXT(lPath, world.Leveldataoverride)
	fileUtils.WriterTXT(mPath, world.Modoverrides)

	serverBuf := c.ParseTemplate(MASTER_SERVER_INI_TEMPLATE, world.ServerIni)

	fileUtils.WriterTXT(sPath, serverBuf)
}

func (c *ClusterService) SaveCaves(world *cluster.World) {

	lPath := constant.GET_CAVES_LEVELDATAOVERRIDE_PATH()
	mPath := constant.GET_CAVES_MOD_PATH()
	sPath := constant.GET_CAVES_DIR_SERVER_INI_PATH()

	fileUtils.CreateFileIfNotExists(lPath)
	fileUtils.CreateFileIfNotExists(mPath)
	fileUtils.CreateFileIfNotExists(sPath)

	fileUtils.WriterTXT(lPath, world.Leveldataoverride)
	fileUtils.WriterTXT(mPath, world.Modoverrides)

	serverBuf := c.ParseTemplate(CAVES_SERVER_INI_TEMPLATE, world.ServerIni)

	fileUtils.WriterTXT(sPath, serverBuf)
}

func (c *ClusterService) GetGameConfog() *cluster.GameConfig {
	gameConfig := cluster.GameConfig{}
	var wg sync.WaitGroup
	wg.Add(6)

	go func() {
		gameConfig.ClusterToken = c.ReadClusterTokenFile()
		wg.Done()
	}()
	go func() {
		gameConfig.Cluster = c.ReadClusterIniFile()
		wg.Done()
	}()
	go func() {
		gameConfig.Adminlist = c.ReadAdminlistFile()
		wg.Done()
	}()
	go func() {
		gameConfig.Blocklist = c.ReadBlocklistFile()
		wg.Done()
	}()
	go func() {
		gameConfig.Master = c.ReadMaster()
		wg.Done()
	}()
	go func() {
		gameConfig.Caves = c.ReadCaves()
		wg.Done()
	}()
	wg.Wait()
	return &gameConfig
}

func (c *ClusterService) SaveGameConfig(gameConfig *cluster.GameConfig) {

	var wg sync.WaitGroup
	wg.Add(6)

	go func() {
		c.SaveClusterToken(gameConfig.ClusterToken)
		wg.Done()
	}()

	go func() {
		c.SaveClusterIni(gameConfig.Cluster)
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
		c.SaveMaster(gameConfig.Master)
		c.DedicatedServerModsSetup(gameConfig.Master.Modoverrides)
		wg.Done()
	}()

	go func() {
		c.SaveCaves(gameConfig.Caves)
		wg.Done()
	}()

	wg.Wait()
}
