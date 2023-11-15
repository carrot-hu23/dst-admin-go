package vo

type GameArchive struct {
	ClusterName        string     `json:"clusterName"`
	ClusterDescription string     `json:"clusterDescription"`
	ClusterPassword    string     `json:"clusterPassword"`
	GameMod            string     `json:"gameMod"`
	Players            []PlayerVO `json:"players"`
	MaxPlayers         int        `json:"maxPlayers"`
	Mods               int        `json:"mods"`
	IpConnect          string     `json:"ipConnect"`
	Port               uint       `json:"port"`
	Ip                 string     `json:"ip"`
	Meta               Meta       `json:"meta"`
	Version            int64      `json:"version"`
	LastVersion        int64      `json:"lastVersion"`
}

func NewGameArchie() *GameArchive {
	return &GameArchive{}
}

type Clock struct {
	TotalTimeInPhase     int     `lua:"totaltimeinphase"`
	Cycles               int     `lua:"cycles"`
	Phase                string  `lua:"phase"`
	RemainingTimeInPhase float64 `lua:"remainingtimeinphase"`
	MooomPhaseCycle      int     `lua:"mooomphasecycle"`
	Segs                 Segs    `lua:"segs"`
}

type Segs struct {
	Night int `lua:"night"`
	Day   int `lua:"day"`
	Dusk  int `lua:"dusk"`
}

type IsRandom struct {
	Summer bool `lua:"summer"`
	Autumn bool `lua:"autumn"`
	Spring bool `lua:"spring"`
	Winter bool `lua:"winter"`
}

type Lengths struct {
	Summer int `lua:"summer"`
	Autumn int `lua:"autumn"`
	Spring int `lua:"spring"`
	Winter int `lua:"winter"`
}

type Seasons struct {
	Premode               bool                   `lua:"premode"`
	Season                string                 `lua:"season"`
	ElapsedDaysInSeason   int                    `lua:"elapseddaysinseason"`
	IsRandom              IsRandom               `lua:"israndom"`
	Lengths               Lengths                `lua:"lengths"`
	RemainingDaysInSeason int                    `lua:"remainingdaysinseason"`
	Mode                  string                 `lua:"mode"`
	TotalDaysInSeason     int                    `lua:"totaldaysinseason"`
	Segs                  map[string]interface{} `lua:"segs"`
}

type Meta struct {
	Clock   Clock   `lua:"clock"`
	Seasons Seasons `lua:"seasons"`
}
