package third

type DstHomeDetailParam struct {
	RowId  string `json:"rowId"`
	Region string `json:"region"`
}

func NewDstHomeDetailParam() *DstHomeDetailParam {
	return &DstHomeDetailParam{}
}
