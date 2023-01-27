package vo

type UserVO struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	SessionID string `json:"sessionId"`
}

func NewUserVO() *UserVO {
	return &UserVO{}
}
