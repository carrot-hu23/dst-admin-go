package vo

type WhitelistVO struct {
	Whitelist []string `json:"whitelist"`
}

func NewWhitelistVO() *WhitelistVO {
	return &WhitelistVO{}
}
