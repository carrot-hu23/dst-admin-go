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
