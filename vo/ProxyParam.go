package vo

type ProxyParam struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Ip          string `json:"ip"`
	Port        string `json:"port"`
}
