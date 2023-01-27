package vo

type BlacklistVO struct {
	Blacklist []string `json:"blacklist"`
}

func NewBlacklistVO() *BlacklistVO {
	return &BlacklistVO{}
}
