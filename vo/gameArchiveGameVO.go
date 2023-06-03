package vo

type GameArchiveVO struct {
	ClusterName     string   `json:"clusterName"`
	GameMode        string   `json:"gameMode"`
	MaxPlayers      uint8    `json:"maxPlayers"`
	ClusterPassword string   `json:"clusterPassword"`
	PlayDay         string   `json:"playDay"`
	Season          string   `json:"season"`
	TotalModNum     int      `json:"totalModNum"`
	WorkshopIds     []string `json:"workshopIds"`
	Modoverrides    string   `json:"modoverrides"`
}

func NewGameArchieVO() *GameArchiveVO {
	return &GameArchiveVO{}
}

type GameArchive struct {
	ClusterName     string     `json:"clusterName"`
	ClusterPassword string     `json:"clusterPassword"`
	GameMod         string     `json:"gameMod"`
	Players         []PlayerVO `json:"players"`
	MaxPlayers      int        `json:"maxPlayers"`
	Days            int        `json:"days"`
	Season          string     `json:"season"`
	Mods            int        `json:"mods"`
	IpConnect       string     `json:"ipConnect"`
	Meta            string     `json:"meta"`
}

func NewGameArchie() *GameArchive {
	return &GameArchive{
		Days:   0,
		Season: "unknown",
		Mods:   0,
	}
}
