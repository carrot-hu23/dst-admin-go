package gameConfig

import (
	"dst-admin-go/internal/pkg/utils/collectionUtils"
	"dst-admin-go/internal/pkg/utils/dstUtils"
	"dst-admin-go/internal/pkg/utils/fileUtils"
	"dst-admin-go/internal/service/archive"
	"dst-admin-go/internal/service/levelConfig"
	"log"
	"path/filepath"
	"strings"

	"github.com/go-ini/ini"
)

const (
	ClusterIniTemplate      = "./static/template/cluster2.ini"
	MasterServerIniTemplate = "./static/template/master_server.ini"
	CavesServerIniTemplate  = "./static/template/caves_server.ini"
	ServerIniTemplate       = "./static/template/server.ini"
)

type ClusterIni struct {
	// [GAEMPLAY]
	GameMode        string `json:"game_mode"`
	MaxPlayers      uint   `json:"max_players"`
	Pvp             bool   `json:"pvp"`
	PauseWhenNobody bool   `json:"pause_when_nobody"`
	VoteEnabled     bool   `json:"vote_enabled"`
	VoteKickEnabled bool   `json:"vote_kick_enabled"`

	// [NETWORK]
	LanOnlyCluster     bool   `json:"lan_only_cluster"`
	ClusterIntention   string `json:"cluster_intention"`
	ClusterDescription string `json:"cluster_description"`
	ClusterPassword    string `json:"cluster_password"`
	ClusterName        string `json:"cluster_name"`
	OfflineCluster     bool   `json:"offline_cluster"`
	ClusterLanguage    string `json:"cluster_language"`
	WhitelistSlots     uint   `json:"whitelist_slots"`
	TickRate           uint   `json:"tick_rate"`

	// [MISC]
	ConsoleEnabled bool `json:"console_enabled"`
	MaxSnapshots   uint `json:"max_snapshots"`

	// [SHARD]
	ShardEnabled bool   `json:"shard_enabled"`
	BindIp       string `json:"bind_ip"`
	MasterIp     string `json:"master_ip"`
	MasterPort   uint   `json:"master_port"`
	ClusterKey   string `json:"cluster_key"`

	// [STEAM]
	SteamGroupId     string `json:"steam_group_id"`
	SteamGroupOnly   bool   `json:"steam_group_only"`
	SteamGroupAdmins bool   `json:"steam_group_admins"`
}
type ServerIni struct {

	// [NETWORK]
	ServerPort uint `json:"server_port"`
	// [SHARD]
	IsMaster bool   `json:"is_master"`
	Name     string `json:"name"`
	Id       uint   `json:"id"`

	// [ACCOUNT]
	EncodeUserPath bool `json:"encode_user_path"`

	// [STEAM]
	AuthenticationPort uint `json:"authentication_port"`
	MasterServerPort   uint `json:"master_server_port"`
}

type ClusterIniConfig struct {
	ClusterIni *ClusterIni `json:"cluster"`
	Token      string      `json:"token"`
}
type GameConfig struct {
	archive          *archive.PathResolver
	levelConfigUtils *levelConfig.LevelConfigUtils
}

func NewGameConfig(archive *archive.PathResolver, levelConfigUtils *levelConfig.LevelConfigUtils) *GameConfig {
	return &GameConfig{
		archive:          archive,
		levelConfigUtils: levelConfigUtils,
	}
}

func (p *GameConfig) GetClusterIniConfig(clusterName string) (ClusterIniConfig, error) {
	clusterIni, err := p.GetClusterIni(clusterName)
	if err != nil {
		return ClusterIniConfig{}, err
	}
	clusterToken, err := p.GetClusterToken(clusterName)
	if err != nil {
		return ClusterIniConfig{}, err
	}
	return ClusterIniConfig{
		ClusterIni: &clusterIni,
		Token:      clusterToken,
	}, nil
}

func (p *GameConfig) SaveClusterIniConfig(clusterName string, config *ClusterIniConfig) error {
	err := p.SaveClusterIni(clusterName, config.ClusterIni)
	if err != nil {
		return err
	}
	err = p.SaveClusterToken(clusterName, config.Token)
	if err != nil {
		return err
	}
	return nil
}

func (p *GameConfig) GetClusterIni(clusterName string) (ClusterIni, error) {
	// return fileUtils.ReadLnFile(p.archive.ClusterPath(clusterName))
	// 加载 INI 文件
	clusterIniPath := p.archive.ClusterIniPath(clusterName)
	if !fileUtils.Exists(clusterIniPath) {
		err := fileUtils.CreateFileIfNotExists(clusterIniPath)
		if err != nil {
			return ClusterIni{}, err
		}
	}
	cfg, err := ini.Load(clusterIniPath)
	if err != nil {
		log.Panicln("Failed to load INI file:", err)
	}

	// [GAMEPLAY]
	GAMEPLAY := cfg.Section("GAMEPLAY")

	newClusterIni := ClusterIni{}

	newClusterIni.GameMode = GAMEPLAY.Key("game_mode").String()
	newClusterIni.MaxPlayers = GAMEPLAY.Key("max_players").MustUint(8)
	newClusterIni.Pvp = GAMEPLAY.Key("pvp").MustBool(false)
	newClusterIni.PauseWhenNobody = GAMEPLAY.Key("pause_when_empty").MustBool(true)
	newClusterIni.VoteEnabled = GAMEPLAY.Key("vote_enabled").MustBool(true)
	newClusterIni.VoteKickEnabled = GAMEPLAY.Key("vote_kick_enabled").MustBool(true)

	// [NETWORK]
	NETWORK := cfg.Section("NETWORK")

	newClusterIni.LanOnlyCluster = NETWORK.Key("lan_only_cluster").MustBool(false)
	newClusterIni.ClusterIntention = NETWORK.Key("cluster_intention").String()
	newClusterIni.ClusterPassword = NETWORK.Key("cluster_password").String()
	newClusterIni.ClusterDescription = NETWORK.Key("cluster_description").String()
	newClusterIni.ClusterName = NETWORK.Key("cluster_name").String()
	newClusterIni.OfflineCluster = NETWORK.Key("offline_cluster").MustBool(false)
	newClusterIni.ClusterLanguage = NETWORK.Key("cluster_language").MustString("zh")
	newClusterIni.WhitelistSlots = NETWORK.Key("whitelist_slots").MustUint(0)
	newClusterIni.TickRate = NETWORK.Key("tick_rate").MustUint(15)

	// [MISC]
	MISC := cfg.Section("MISC")

	newClusterIni.ConsoleEnabled = MISC.Key("console_enabled").MustBool(true)
	newClusterIni.MaxSnapshots = MISC.Key("max_snapshots").MustUint(6)

	// [SHARD]
	SHARD := cfg.Section("SHARD")

	newClusterIni.ShardEnabled = SHARD.Key("shard_enabled").MustBool(true)
	newClusterIni.BindIp = SHARD.Key("bind_ip").MustString("127.0.0.1")
	newClusterIni.MasterIp = SHARD.Key("master_ip").MustString("127.0.0.1")
	newClusterIni.MasterPort = SHARD.Key("master_port").MustUint(10888)
	newClusterIni.ClusterKey = SHARD.Key("cluster_key").String()

	// [STEAM]
	STEAM := cfg.Section("STEAM")

	newClusterIni.SteamGroupId = STEAM.Key("steam_group_id").MustString("")
	newClusterIni.SteamGroupOnly = STEAM.Key("steam_group_only").MustBool(false)
	newClusterIni.SteamGroupAdmins = STEAM.Key("steam_group_admins").MustBool(false)

	clusterIni, err := fileUtils.ReadLnFile(clusterIniPath)
	if err != nil {
		panic("read cluster.ini file error: " + err.Error())
	}
	for _, value := range clusterIni {
		if value == "" {
			continue
		}
		if strings.Contains(value, "cluster_password") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				newClusterIni.ClusterPassword = s
			}
		}
		if strings.Contains(value, "cluster_description") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				newClusterIni.ClusterDescription = s
			}
		}
		if strings.Contains(value, "cluster_name") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				newClusterIni.ClusterName = s
			}
		}
	}
	return newClusterIni, nil
}

func (p *GameConfig) SaveClusterIni(clusterName string, clusterIni *ClusterIni) error {
	clusterIniPath := p.archive.ClusterIniPath(clusterName)
	err := fileUtils.WriterTXT(clusterIniPath, dstUtils.ParseTemplate(ClusterIniTemplate, clusterIni))
	return err
}

func (p *GameConfig) GetClusterToken(clusterName string) (string, error) {
	return fileUtils.ReadFile(p.archive.ClusterTokenPath(clusterName))
}

func (p *GameConfig) SaveClusterToken(clusterName string, token string) error {
	return fileUtils.WriterTXT(p.archive.ClusterTokenPath(clusterName), token)
}

func (p *GameConfig) GetAdminList(clusterName string) ([]string, error) {
	return fileUtils.ReadLnFile(p.archive.AdminlistPath(clusterName))
}

func (p *GameConfig) GetBlackList(clusterName string) ([]string, error) {
	return fileUtils.ReadLnFile(p.archive.BlacklistPath(clusterName))
}

func (p *GameConfig) GetWhithList(clusterName string) ([]string, error) {
	return fileUtils.ReadLnFile(p.archive.WhitelistPath(clusterName))
}

func (p *GameConfig) SaveAdminList(clusterName string, list []string) error {
	path := p.archive.AdminlistPath(clusterName)
	err := fileUtils.CreateFileIfNotExists(path)
	if err != nil {
		return err
	}
	lnFile, err := fileUtils.ReadLnFile(path)
	set := collectionUtils.ToSet(append(lnFile, list...))
	err = fileUtils.WriterLnFile(path, set)
	return err
}

func (p *GameConfig) SaveBlackList(clusterName string, list []string) error {
	path := p.archive.BlacklistPath(clusterName)
	err := fileUtils.CreateFileIfNotExists(path)
	if err != nil {
		return err
	}
	lnFile, err := fileUtils.ReadLnFile(path)
	set := collectionUtils.ToSet(append(lnFile, list...))
	err = fileUtils.WriterLnFile(path, set)
	return err
}

func (p *GameConfig) SaveWhithList(clusterName string, list []string) error {
	path := p.archive.WhitelistPath(clusterName)
	err := fileUtils.CreateFileIfNotExists(path)
	if err != nil {
		return err
	}
	lnFile, err := fileUtils.ReadLnFile(path)
	set := collectionUtils.ToSet(append(lnFile, list...))
	err = fileUtils.WriterLnFile(path, set)
	return err
}

type HomeConfigVO struct {
	ClusterIntention   string `json:"clusterIntention"`
	ClusterName        string `json:"clusterName"`
	ClusterDescription string `json:"clusterDescription"`
	GameMode           string `json:"gameMode"`
	Pvp                bool   `json:"pvp"`
	MaxPlayers         uint   `json:"maxPlayers"`
	MaxSnapshots       uint   `json:"max_snapshots"`
	ClusterPassword    string `json:"clusterPassword"`
	Token              string `json:"token"`
	MasterMapData      string `json:"masterMapData"`
	CavesMapData       string `json:"cavesMapData"`
	ModData            string `json:"modData"`
	Otype              int64  `json:"type"`
	PauseWhenNobody    bool   `json:"pause_when_nobody"`
	VoteEnabled        bool   `json:"vote_enabled"`
}

func (p *GameConfig) GetHomeConfig(clusterName string) (HomeConfigVO, error) {
	clusterToken, err := p.GetClusterToken(clusterName)
	if err != nil {
		return HomeConfigVO{}, err
	}
	clusterIni, err := p.GetClusterIni(clusterName)
	if err != nil {
		return HomeConfigVO{}, err
	}
	masterData, err := fileUtils.ReadFile(p.archive.DataFilePath(clusterName, "Master", "leveldataoverride.lua"))
	if err != nil {
		return HomeConfigVO{}, err
	}
	cavesData, err := fileUtils.ReadFile(p.archive.DataFilePath(clusterName, "Caves", "leveldataoverride.lua"))
	if err != nil {
		return HomeConfigVO{}, err
	}
	modData, err := fileUtils.ReadFile(p.archive.DataFilePath(clusterName, "Master", "modoverrides.lua"))
	if err != nil {
		return HomeConfigVO{}, err
	}

	homeConfigVo := HomeConfigVO{
		ClusterIntention:   clusterIni.ClusterIntention,
		ClusterName:        clusterIni.ClusterName,
		ClusterDescription: clusterIni.ClusterDescription,
		GameMode:           clusterIni.GameMode,
		Pvp:                clusterIni.Pvp,
		MaxPlayers:         clusterIni.MaxPlayers,
		MaxSnapshots:       clusterIni.MaxSnapshots,
		ClusterPassword:    clusterIni.ClusterPassword,
		Token:              clusterToken,
		PauseWhenNobody:    clusterIni.PauseWhenNobody,
		VoteEnabled:        clusterIni.VoteEnabled,
		MasterMapData:      masterData,
		CavesMapData:       cavesData,
		ModData:            modData,
	}
	return homeConfigVo, nil
}

func (p *GameConfig) SaveConfig(clusterName string, homeConfig HomeConfigVO) {
	modConfig := homeConfig.ModData
	if modConfig != "" {
		config, _ := p.levelConfigUtils.GetLevelConfig(clusterName)
		for i := range config.LevelList {
			clusterPath := p.archive.ClusterPath(clusterName)
			fileUtils.WriterTXT(filepath.Join(clusterPath, config.LevelList[i].File, "modoverrides.lua"), modConfig)
		}
		var serverModSetup = ""
		workshopIds := dstUtils.WorkshopIds(modConfig)
		for _, workshopId := range workshopIds {
			serverModSetup += "ServerModSetup(\"" + workshopId + "\")\n"
		}
		fileUtils.WriterTXT(p.archive.GetModSetup(clusterName), serverModSetup)
	}
}
