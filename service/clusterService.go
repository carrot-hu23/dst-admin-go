package service

import (
	"dst-admin-go/constant"
	"dst-admin-go/utils/fileUtils"
	"log"
	"path"
	"strconv"
	"strings"
	"sync"
)

type Cluster struct {
	ClusterIntention   string `json:"clusterIntention"`
	ClusterName        string `json:"clusterName"`
	ClusterDescription string `json:"clusterDescription"`
	GameMode           string `json:"gameMode"`
	Pvp                bool   `json:"pvp"`
	MaxPlayers         uint8  `json:"maxPlayers"`
	MaxSnapshots       uint8  `json:"max_snapshots"`
	ClusterPassword    string `json:"clusterPassword"`
	Token              string `json:"token"`
	PauseWhenNobody    bool   `json:"pause_when_nobody"`
	VoteEnabled        bool   `json:"vote_enabled"`
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
}

type World struct {
	WorldName         string `json:"world_name"`
	Description       string `json:"description"`
	Type              string `json:"type"`
	Leveldataoverride string `json:"leveldataoverride"`
	Modoverrides      string `json:"modoverrides"`
	ServerIni         `json:"server_ini"`
}

type MultiLevelWorldConfig struct {
	Cluster      *Cluster `json:"cluster"`
	ClusterToken string   `json:"cluster_token"`
	Adminlist    string   `json:"adminlist"`
	Blocklist    string   `json:"blocklist"`
	Worlds       []World  `json:"worlds"`
}

func ReadClusterIniFile() (cluster *Cluster) {
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
				cluster.GameMode = s
			}
		}
		if strings.Contains(value, "max_players") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				n, err := strconv.ParseUint(s, 10, 8)
				if err == nil {
					cluster.MaxPlayers = uint8(n)
				}
			}
		}
		if strings.Contains(value, "pvp") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				b, err := strconv.ParseBool(s)
				if err == nil {
					cluster.Pvp = b
				}
			}
		}
		if strings.Contains(value, "pause_when_empty") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				b, err := strconv.ParseBool(s)
				if err == nil {
					cluster.PauseWhenNobody = b
				}
			}
		}
		if strings.Contains(value, "cluster_intention") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				cluster.ClusterIntention = s
			}
		}
		if strings.Contains(value, "cluster_password") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				cluster.ClusterPassword = s
			}
		}
		if strings.Contains(value, "cluster_description") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				cluster.ClusterDescription = s
			}
		}
		if strings.Contains(value, "cluster_name") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				cluster.ClusterName = s
			}
		}
	}
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
		log.Println("路径不存在", path)
		return
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
		return
	}
	str, err := fileUtils.ReadLnFile(path)
	log.Println("str:", str)
	if err != nil {
		panic("read dst blocklist.txt error: \n" + err.Error())
	}
	return
}

func ReadLeveldataoverrideFile(filepath string) string {
	return ""
}

func ReadModoverridesFile(filepath string) string {
	return ""
}

func ReadServerIniFile(filepath string) *ServerIni {

	return &ServerIni{}
}

func GetmultiLevelWorldConfig() *MultiLevelWorldConfig {
	multiLevelWorldConfig := &MultiLevelWorldConfig{}

	cluster_path := constant.GET_CLUSTER_INI_PATH()
	worldDirs, err := fileUtils.FindWorldDirs(cluster_path)
	if err != nil {
		log.Panicln("路径查找失败", err)
	}
	size := len(worldDirs)
	worlds := make([]World, size)
	var wg sync.WaitGroup
	wg.Add(size)
	for i, worldPath := range worldDirs {
		//查找
		go func(word World, worldPath string) {
			word.Leveldataoverride = ReadLeveldataoverrideFile(path.Join(worldPath, "leveldataoverride.lua"))
			word.Modoverrides = ReadModoverridesFile(path.Join(worldPath, "modoverrides.lua"))
			word.ServerIni = *ReadServerIniFile(path.Join(worldPath, "server_ini.lua"))
			wg.Done()
		}(worlds[i], worldPath)
	}
	wg.Wait()
	multiLevelWorldConfig.Worlds = worlds
	return multiLevelWorldConfig
}
