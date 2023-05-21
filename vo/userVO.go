package vo

type UserVO struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	SessionID string `json:"sessionId"`
}

func NewUserVO() *UserVO {
	return &UserVO{}
}

type UserInfo struct {
	Username     string `json:"username"`
	DisplayeName string `json:"displayeName"`
	Password     string `json:"password"`
	PhotoURL     string `json:"photoURL"`
}
