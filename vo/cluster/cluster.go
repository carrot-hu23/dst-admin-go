package cluster

type Cluster struct {
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
	ClusterPassword    string `json:"clusterPassword"`
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
	SteamGroupOnly   bool `json:"steam_group_only"`
	SteamGroupId     uint `json:"steam_group_id"`
	SteamGroupAdmins bool `json:"steam_group_admins"`
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
	Master_serverPort  uint `json:"master_server_port"`
}

func NewCluster() *Cluster {
	return &Cluster{
		Pvp:              false,
		PauseWhenNobody:  true,
		VoteEnabled:      true,
		VoteKickEnabled:  true,
		LanOnlyCluster:   false,
		ClusterLanguage:  "zh",
		WhitelistSlots:   0,
		TickRate:         15,
		ConsoleEnabled:   true,
		MaxSnapshots:     6,
		ShardEnabled:     true,
		BindIp:           "127.0.0.1",
		MasterIp:         "127.0.0.1",
		MasterPort:       108888,
		ClusterKey:       "",
		SteamGroupOnly:   false,
		SteamGroupId:     0,
		SteamGroupAdmins: false,
	}
}

func NewMasterServerIni() *ServerIni {
	return &ServerIni{
		ServerPort:     10999,
		IsMaster:       true,
		Name:           "Master",
		Id:             10000,
		EncodeUserPath: true,
	}
}

func NewCavesServerIni() *ServerIni {
	return &ServerIni{
		ServerPort:         10998,
		IsMaster:           false,
		Name:               "Caves",
		Id:                 10010,
		EncodeUserPath:     true,
		AuthenticationPort: 8766,
		Master_serverPort:  27016,
	}
}

type World struct {
	WorldName         string     `json:"world_name"`
	IsMaster          bool       `json:"is_master"`
	Leveldataoverride string     `json:"leveldataoverride"`
	Modoverrides      string     `json:"modoverrides"`
	ServerIni         *ServerIni `json:"server_ini"`
}

type MultiLevelWorldConfig struct {
	Cluster      *Cluster `json:"cluster"`
	ClusterToken string   `json:"cluster_token"`
	Adminlist    string   `json:"adminlist"`
	Blocklist    string   `json:"blocklist"`
	Worlds       []World  `json:"worlds"`
}
