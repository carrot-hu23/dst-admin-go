package login

import (
	"dst-admin-go/internal/config"
	"dst-admin-go/internal/pkg/response"
	"dst-admin-go/internal/pkg/utils/fileUtils"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

const (
	PasswordPath = "./password.txt"
)

type LoginService struct {
	config *config.Config
}

type UserInfo struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName"`
	PhotoURL    string `json:"photoURL"`
}

func NewLoginService(config *config.Config) *LoginService {
	return &LoginService{
		config: config,
	}
}

func (l *LoginService) GetUserInfo() UserInfo {
	user, err := fileUtils.ReadLnFile(PasswordPath)

	if err != nil {
		log.Panicln("Not find password file error: " + err.Error())
	}

	username := strings.TrimSpace(strings.Split(user[0], "=")[1])
	// password := strings.TrimSpace(strings.Split(user[1], "=")[1])
	displayName := strings.TrimSpace(strings.Split(user[2], "=")[1])
	photoURL := strings.TrimSpace(strings.Split(user[3], "=")[1])

	return UserInfo{
		Username:    username,
		DisplayName: displayName,
		PhotoURL:    photoURL,
	}
}

func (l *LoginService) Login(userInfo UserInfo, ctx *gin.Context) *response.Response {

	response := &response.Response{}

	user, err := fileUtils.ReadLnFile(PasswordPath)
	if err != nil {
		log.Panicln("Not find password file error: " + err.Error())
	}

	username := strings.TrimSpace(strings.Split(user[0], "=")[1])
	password := strings.TrimSpace(strings.Split(user[1], "=")[1])
	displayName := strings.TrimSpace(strings.Split(user[2], "=")[1])
	photoURL := strings.TrimSpace(strings.Split(user[3], "=")[1])
	white := l.IsWhiteIP(ctx)
	if !white {
		if username != userInfo.Username || password != userInfo.Password {
			log.Panicln("User authentication failed")
			response.Code = 401
			response.Msg = "User authentication failed"
			return response
		}
	}
	session := sessions.Default(ctx)
	session.Set("username", username)
	err = session.Save()
	if err != nil {
		log.Panicln(err)
	}

	response.Code = 200
	response.Msg = "Login success"
	response.Data = map[string]interface{}{
		"username":    username,
		"displayName": displayName,
		"photoURL":    photoURL,
	}

	return response
}

func (l *LoginService) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	err := session.Save()
	if err != nil {
		log.Panicln(err)
	}
}

func (l *LoginService) DirectLogin(ctx *gin.Context) {
	user, err := fileUtils.ReadLnFile(PasswordPath)
	if err != nil {
		log.Panicln("Not find password file error: " + err.Error())
	}
	username := strings.TrimSpace(strings.Split(user[0], "=")[1])
	session := sessions.Default(ctx)
	session.Set("username", username)
}

func (l *LoginService) ChangeUser(username, password string) {
	user, err := fileUtils.ReadLnFile(PasswordPath)
	if err != nil {
		log.Panicln("Not find password file error: " + err.Error())
	}
	displayName := strings.TrimSpace(strings.Split(user[2], "=")[1])
	photoURL := strings.TrimSpace(strings.Split(user[3], "=")[1])
	fileUtils.WriterLnFile(PasswordPath, []string{
		"username = " + username,
		"password = " + password,
		"displayName=" + displayName,
		"photoURL=" + photoURL,
	})
}

func (l *LoginService) ChangePassword(newPassword string) *response.Response {

	response := &response.Response{}
	user, err := fileUtils.ReadLnFile(PasswordPath)

	if err != nil {
		log.Panicln("Not find password file error: " + err.Error())
	}
	username := strings.TrimSpace(strings.Split(user[0], "=")[1])
	displayName := strings.TrimSpace(strings.Split(user[2], "=")[1])
	photoURL := strings.TrimSpace(strings.Split(user[3], "=")[1])
	fileUtils.WriterLnFile(PasswordPath, []string{
		"username = " + username,
		"password = " + newPassword,
		"displayName=" + displayName,
		"photoURL=" + photoURL,
	})

	response.Code = 200
	response.Msg = "Update user new password success"

	return response
}

func (l *LoginService) InitUserInfo(userInfo UserInfo) {
	username := "username=" + userInfo.Username
	password := "password=" + userInfo.Password
	displayName := "displayName=" + userInfo.DisplayName
	photoURL := "photoURL=" + userInfo.PhotoURL
	fileUtils.WriterLnFile(PasswordPath, []string{username, password, displayName, photoURL})
}

func (l *LoginService) IsWhiteIP(ctx *gin.Context) bool {
	if l.config == nil {
		return false
	}
	WhiteAdminIP := l.config.WhiteAdminIP
	if WhiteAdminIP != "" {
		//
		ipaddr := ctx.Request.RemoteAddr
		ip, _, _ := net.SplitHostPort(ipaddr)
		if ip != "" {
			ipnet := net.ParseIP(ip)
			adminips := strings.Split(WhiteAdminIP, ",")
			//fmt.Println(ipnet)
			for _, s := range adminips {
				if strings.Count(s, "/") > 0 {
					_, netadmin, err := net.ParseCIDR(s)
					//fmt.Println(netadmin)
					//fmt.Println(len)
					if err != nil {
						fmt.Printf("Error parsing CIDR: %v\n", err)
					}
					if netadmin.Contains(ipnet) {
						return true
					}
				} else {
					netadmin := net.ParseIP(s)
					if netadmin != nil && netadmin.Equal(ipnet) {
						return true
					}
				}
			}

		}
	}
	return false
}
