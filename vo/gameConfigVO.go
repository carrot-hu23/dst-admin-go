package vo

type GameConfigVO struct {
	ClusterIntention   string `json:"clusterIntention"`
	ClusterName        string `json:"clusterName"`
	ClusterDescription string `json:"clusterDescription"`
	GameMode           string `json:"gameMode"`
	Pvp                bool   `json:"pvp"`
	MaxPlayers         uint8  `json:"maxPlayers"`
	MaxSnapshots       uint8  `json:"max_snapshots"`
	ClusterPassword    string `json:"clusterPassword"`
	Token              string `json:"token"`
	MasterMapData      string `json:"masterMapData"`
	CavesMapData       string `json:"cavesMapData"`
	ModData            string `json:"modData"`
	Otype              uint8  `json:"type"`
	PauseWhenNobody    bool   `json:"pause_when_nobody"`
	VoteEnabled        bool   `json:"vote_enabled"`
}

func NewGameConfigVO() *GameConfigVO {
	return &GameConfigVO{
		Pvp:             false,
		MaxPlayers:      6,
		MaxSnapshots:    6,
		PauseWhenNobody: false,
		VoteEnabled:     true,
	}
}
