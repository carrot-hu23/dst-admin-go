package vo

type AdminListVO struct {
	AdminList []string `json:"adminList"`
}

func NewAdminListVO() *AdminListVO {
	return &AdminListVO{}
}
