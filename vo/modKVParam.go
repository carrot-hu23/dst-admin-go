package vo

type ModKVParam struct {
	UserId    string `json:"UserId"`
	ModId     int    `json:"modId"`
	ModConfig string `json:"ModConfig"`
	Version   string `json:"Version"`
}
