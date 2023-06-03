package service

import (
	"dst-admin-go/constant"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo/cluster"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/go-ini/ini"
)

const (
	CLUSTER_INI_TEMPLATE       = "./static/template/cluster2.ini"
	MASTER_SERVER_INI_TEMPLATE = "./static/template/master_server.ini"
	CAVES_SERVER_INI_TEMPLATE  = "./static/template/caves_server.ini"
)

func ReadClusterIniFile() *cluster.Cluster {
	cluster := cluster.NewCluster()
	// 加载 INI 文件
	cluster_ini_path := constant.GET_CLUSTER_INI_PATH()
	if !fileUtils.Exists(cluster_ini_path) {
		createFileIfNotExsists(cluster_ini_path)
		return cluster
	}
	cfg, err := ini.Load(cluster_ini_path)
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
	cluster.MasterPort = SHARD.Key("master_port").MustUint(10888)
	cluster.ClusterKey = SHARD.Key("cluster_key").String()

	// [STEAM]
	STEAM := cfg.Section("STEAM")

	cluster.SteamGroupOnly = STEAM.Key("steam_group_only").MustBool(false)
	cluster.SteamGroupId = STEAM.Key("steam_group_id").MustUint(0)
	cluster.SteamGroupAdmins = STEAM.Key("steam_group_admins").MustString("")

	return cluster
}

func ReadClusterTokenFile() string {
	cluster_token_path := constant.GET_CLUSTER_TOKEN_PATH()
	if !fileUtils.Exists(cluster_token_path) {
		createFileIfNotExsists(cluster_token_path)
		return ""
	}

	token, err := fileUtils.ReadFile(cluster_token_path)
	if err != nil {
		panic("read cluster_token.txt file error: " + err.Error())
	}
	return token
}

func ReadAdminlistFile() (str []string) {
	path := constant.GET_DST_ADMIN_LIST_PATH()
	createFileIfNotExsists(path)
	str, err := fileUtils.ReadLnFile(path)
	log.Println("str:", str)
	if err != nil {
		panic("read dst adminlist.txt error: \n" + err.Error())
	}
	return
}

func ReadBlocklistFile() (str []string) {
	path := constant.GET_DST_BLOCKLIST_PATH()
	createFileIfNotExsists(path)
	str, err := fileUtils.ReadLnFile(path)
	log.Println("str:", str)
	if err != nil {
		panic("read dst blocklist.txt error: \n" + err.Error())
	}
	return
}

func ReadLeveldataoverrideFile(filepath string) string {
	if !fileUtils.Exists(filepath) {
		createFileIfNotExsists(filepath)
		return "return {}"
	}

	leveldataoverride, err := fileUtils.ReadFile(filepath)
	if err != nil {
		panic("read leveldataoverride.lua file error: " + err.Error())
	}
	return leveldataoverride
}

func ReadModoverridesFile(filepath string) string {
	if !fileUtils.Exists(filepath) {
		createFileIfNotExsists(filepath)
		return "return {}"
	}
	modoverrides, err := fileUtils.ReadFile(filepath)
	if err != nil {
		panic("read modoverrides.lua file error: " + err.Error())
	}
	return modoverrides
}

func ReadServerIniFile(filepath string, isMaster bool) *cluster.ServerIni {
	createFileIfNotExsists(filepath)
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

	serverIni.AuthenticationPort = STEAM.Key("authentication_port").String()
	serverIni.MasterServerPort = STEAM.Key("master_server_port").String()

	return serverIni
}

func isMaster(filePath string) bool {
	return strings.Contains(filePath, "Master") || strings.Contains(filePath, "master")
}

func GetmultiLevelWorldConfig() *cluster.MultiLevelWorldConfig {
	multiLevelWorldConfig := &cluster.MultiLevelWorldConfig{}

	cluster_path := constant.GET_DST_USER_GAME_CONFG_PATH()
	worldDirs, err := fileUtils.FindWorldDirs(cluster_path)
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
			worlds[i].Leveldataoverride = ReadLeveldataoverrideFile(path.Join(worldPath, "leveldataoverride.lua"))
			worlds[i].Modoverrides = ReadModoverridesFile(path.Join(worldPath, "modoverrides.lua"))
			worlds[i].IsMaster = isMaster(worldPath)
			worlds[i].ServerIni = ReadServerIniFile(path.Join(worldPath, "server.ini"), isMaster(worldPath))
			wg.Done()
		}(i, worldPath)
	}

	go func() {
		multiLevelWorldConfig.ClusterToken = ReadClusterTokenFile()
		wg.Done()
	}()
	go func() {
		multiLevelWorldConfig.Cluster = ReadClusterIniFile()
		wg.Done()
	}()
	go func() {
		multiLevelWorldConfig.Adminlist = ReadAdminlistFile()
		wg.Done()
	}()
	go func() {
		multiLevelWorldConfig.Blocklist = ReadBlocklistFile()
		wg.Done()
	}()

	multiLevelWorldConfig.Worlds = worlds
	wg.Wait()
	return multiLevelWorldConfig
}

func SaveClusterToken(token string) {
	cluster_token_path := constant.GET_CLUSTER_TOKEN_PATH()
	createFileIfNotExsists(cluster_token_path)
	fileUtils.WriterTXT(cluster_token_path, token)
}

func SaveClusterIni(cluster *cluster.Cluster) {
	cluster_ini_path := constant.GET_CLUSTER_INI_PATH()
	createFileIfNotExsists(cluster_ini_path)
	fileUtils.WriterTXT(cluster_ini_path, pareseTemplate(CLUSTER_INI_TEMPLATE, cluster))
}

func SaveAdminlist(str []string) {
	adminlist_path := constant.GET_DST_ADMIN_LIST_PATH()
	createFileIfNotExsists(adminlist_path)
	fileUtils.WriterLnFile(adminlist_path, str)
}

func SaveBlocklist(str []string) {
	blocklist_path := constant.GET_DST_BLOCKLIST_PATH()
	createFileIfNotExsists(blocklist_path)
	fileUtils.WriterLnFile(blocklist_path, str)
}

func createFileIfNotExsists(path string) error {

	// 检查文件是否存在
	_, err := os.Stat(path)
	if err == nil {
		// 文件已经存在，直接返回
		return nil
	}
	if !os.IsNotExist(err) {
		// 其他错误，返回错误信息
		return err
	}

	// 创建文件所在的目录
	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		// 创建目录失败，返回错误信息
		return err
	}

	// 创建文件
	_, err = os.Create(path)
	if err != nil {
		// 创建文件失败，返回错误信息
		return err
	}

	// 创建成功，返回 nil
	return nil
}

func createDirIfNotExsists(filepath string) {
	if !fileUtils.Exists(filepath) {
		fileUtils.CreateDir(filepath)
	}
}

func handleWorldDir(basePath string, world *cluster.World) {

	fileUtils.CreateDir(basePath)
	leveldataoverride_path := path.Join(basePath, "leveldataoverride.lua")
	modoverrides_path := path.Join(basePath, "modoverrides.lua")
	server_ini_path := path.Join(basePath, "server.ini")

	createFileIfNotExsists(leveldataoverride_path)
	createFileIfNotExsists(modoverrides_path)
	createFileIfNotExsists(server_ini_path)

	fileUtils.WriterTXT(leveldataoverride_path, world.Leveldataoverride)
	fileUtils.WriterTXT(modoverrides_path, world.Modoverrides)

	// TODO 写入 constant.DST_MOD_SETTING_PATH

	serverIniBuf := ""
	if world.IsMaster {
		serverIniBuf = pareseTemplate(MASTER_SERVER_INI_TEMPLATE, world.ServerIni)
	} else {
		serverIniBuf = pareseTemplate(CAVES_SERVER_INI_TEMPLATE, world.ServerIni)
	}

	fileUtils.WriterTXT(server_ini_path, serverIniBuf)
}

func SaveWorlds(worlds []cluster.World) {
	var wg sync.WaitGroup
	wg.Add(len(worlds))
	basePath := constant.GET_DST_USER_GAME_CONFG_PATH()
	for _, world := range worlds {
		go func(world *cluster.World) {
			// 判断文件是否存在，如果不存在则创建
			filePath := path.Join(basePath, world.WorldName)
			handleWorldDir(filePath, world)
			wg.Done()
		}(&world)
	}
	wg.Wait()
}

// TODO not test
func SaveMultiLevelWorldConfig(multiLevelWorldConfig *cluster.MultiLevelWorldConfig) {

	var wg sync.WaitGroup
	wg.Add(5)
	var size int32
	go func() {
		SaveClusterToken(multiLevelWorldConfig.ClusterToken)
		wg.Done()
		atomic.AddInt32(&size, 1)
	}()

	go func() {
		SaveClusterIni(multiLevelWorldConfig.Cluster)
		wg.Done()
		atomic.AddInt32(&size, 1)
	}()

	go func() {
		SaveAdminlist(multiLevelWorldConfig.Adminlist)
		wg.Done()
		atomic.AddInt32(&size, 1)
	}()
	go func() {
		SaveBlocklist(multiLevelWorldConfig.Blocklist)
		wg.Done()
		atomic.AddInt32(&size, 1)
	}()
	go func() {
		SaveWorlds(multiLevelWorldConfig.Worlds)
		wg.Done()
		atomic.AddInt32(&size, 1)
	}()

	wg.Wait()
	if size != 5 {
		// 设置失败
	}
}

func ReadMaster() *cluster.World {
	master := cluster.World{}

	master.WorldName = "Master"
	master.IsMaster = true

	master.Leveldataoverride = ReadLeveldataoverrideFile(constant.GET_MASTER_LEVELDATAOVERRIDE_PATH())
	master.Modoverrides = ReadModoverridesFile(constant.GET_MASTER_MOD_PATH())
	master.ServerIni = ReadServerIniFile(constant.GET_MASTER_DIR_SERVER_INI_PATH(), true)

	return &master
}

func ReadCaves() *cluster.World {
	caves := cluster.World{}

	caves.WorldName = "Caves"
	caves.IsMaster = false

	caves.Leveldataoverride = ReadLeveldataoverrideFile(constant.GET_CAVES_LEVELDATAOVERRIDE_PATH())
	caves.Modoverrides = ReadModoverridesFile(constant.GET_CAVES_MOD_PATH())
	caves.ServerIni = ReadServerIniFile(constant.GET_CAVES_DIR_SERVER_INI_PATH(), false)

	return &caves
}

func SaveMaster(world *cluster.World) {

	l_path := constant.GET_MASTER_LEVELDATAOVERRIDE_PATH()
	m_path := constant.GET_MASTER_MOD_PATH()
	s_path := constant.GET_MASTER_DIR_SERVER_INI_PATH()

	createFileIfNotExsists(l_path)
	createFileIfNotExsists(m_path)
	createFileIfNotExsists(s_path)

	fileUtils.WriterTXT(l_path, world.Leveldataoverride)
	fileUtils.WriterTXT(m_path, world.Modoverrides)

	serverBuf := pareseTemplate(MASTER_SERVER_INI_TEMPLATE, world.ServerIni)

	fileUtils.WriterTXT(s_path, serverBuf)
}

func SaveCavesr(world *cluster.World) {

	l_path := constant.GET_CAVES_LEVELDATAOVERRIDE_PATH()
	m_path := constant.GET_CAVES_MOD_PATH()
	s_path := constant.GET_CAVES_DIR_SERVER_INI_PATH()

	createFileIfNotExsists(l_path)
	createFileIfNotExsists(m_path)
	createFileIfNotExsists(s_path)

	fileUtils.WriterTXT(l_path, world.Leveldataoverride)
	fileUtils.WriterTXT(m_path, world.Modoverrides)

	serverBuf := pareseTemplate(CAVES_SERVER_INI_TEMPLATE, world.ServerIni)

	fileUtils.WriterTXT(s_path, serverBuf)
}

func GetGameConfog() *cluster.GameConfig {
	gameConfig := cluster.GameConfig{}
	var wg sync.WaitGroup
	wg.Add(6)

	go func() {
		gameConfig.ClusterToken = ReadClusterTokenFile()
		wg.Done()
	}()
	go func() {
		gameConfig.Cluster = ReadClusterIniFile()
		wg.Done()
	}()
	go func() {
		gameConfig.Adminlist = ReadAdminlistFile()
		wg.Done()
	}()
	go func() {
		gameConfig.Blocklist = ReadBlocklistFile()
		wg.Done()
	}()
	go func() {
		gameConfig.Master = ReadMaster()
		wg.Done()
	}()
	go func() {
		gameConfig.Caves = ReadCaves()
		wg.Done()
	}()
	wg.Wait()
	return &gameConfig
}

func SaveGameConfig(gameConfig *cluster.GameConfig) {

	var wg sync.WaitGroup
	wg.Add(6)

	go func() {
		SaveClusterToken(gameConfig.ClusterToken)
		wg.Done()
	}()

	go func() {
		SaveClusterIni(gameConfig.Cluster)
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
		SaveMaster(gameConfig.Master)
		UpdateDedicatedServerModsSetup(gameConfig.Master.Modoverrides)
		wg.Done()
	}()

	go func() {
		SaveCavesr(gameConfig.Caves)
		wg.Done()
	}()

	wg.Wait()
}
