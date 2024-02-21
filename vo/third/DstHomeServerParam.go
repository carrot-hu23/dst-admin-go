package third

type DstHomeServerParam struct {
	Page           int    `json:"page"`
	Paginate       int    `json:"paginate"`
	SortType       string `json:"sort_type"`
	SortWay        int    `json:"sort_way"`
	Search_type    int    `json:"search_type"`
	Search_content string `json:"search_content"`
	Mode           string `json:"mode"`
	Mod            int    `json:"mod"`
	Season         string `json:"season"`
	Pvp            int    `json:"pvp"`
	Password       int    `json:"password"`
	World          int    `json:"world"`
	Playerpercent  string `json:"playerpercent"`
}

func NewDstHomeServerParam() *DstHomeServerParam {
	return &DstHomeServerParam{}
}
