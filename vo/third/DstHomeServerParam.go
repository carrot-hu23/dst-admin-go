package third

type DstHomeServerParam struct {
	Page           int    `json:"page"`
	Paginate       int    `json:"paginate"`
	SortType       string `json:"sort_type"`
	SortWay        int    `json:"sort_way"`
	Search_type    int    `json:"search_type"`
	Search_content string `json:"search_content"`
	Mod            string `json:"mod"`
}

func NewDstHomeServerParam() *DstHomeServerParam {
	return &DstHomeServerParam{}
}
