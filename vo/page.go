package vo

type Page struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	TotalPages int64       `json:"totalPages"`
	Page       int         `json:"page"`
	Size       int         `json:"size"`
}
