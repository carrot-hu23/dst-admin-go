package service

import (
	"dst-admin-go/config/global"
	"dst-admin-go/constant/consts"
	"dst-admin-go/session"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

type LoginService struct {
}

func (l *LoginService) GetUserInfo() map[string]interface{} {
	user, err := fileUtils.ReadLnFile(consts.PasswordPath)

	if err != nil {
		log.Panicln("Not find password file error: " + err.Error())
	}

	username := strings.TrimSpace(strings.Split(user[0], "=")[1])
	// password := strings.TrimSpace(strings.Split(user[1], "=")[1])
	displayName := strings.TrimSpace(strings.Split(user[2], "=")[1])
	photoURL := strings.TrimSpace(strings.Split(user[3], "=")[1])

	return map[string]interface{}{
		"username":    username,
		"displayName": displayName,
		"photoURL":    photoURL,
	}
}

func (l *LoginService) Login(userVO *vo.UserVO, ctx *gin.Context, sessions *session.Manager) *vo.Response {

	response := &vo.Response{}

	user, err := fileUtils.ReadLnFile(consts.PasswordPath)
	if err != nil {
		log.Panicln("Not find password file error: " + err.Error())
	}

	// username := strings.TrimSpace(strings.Split(user[0], "=")[1])
	// password := strings.TrimSpace(strings.Split(user[1], "=")[1])
	username := strings.TrimSpace(strings.Split(user[0], "=")[1])
	password := strings.TrimSpace(strings.Split(user[1], "=")[1])
	displayName := strings.TrimSpace(strings.Split(user[2], "=")[1])
	photoURL := strings.TrimSpace(strings.Split(user[3], "=")[1])
	white := l.IsWhiteIP(ctx)
	if !white {
		if username != userVO.Username || password != userVO.Password {
			log.Panicln("User authentication failed")
			response.Code = 401
			response.Msg = "User authentication failed"
			return response
		}
	}
	session := sessions.Start(ctx.Writer, ctx.Request)

	session.Set("username", username)

	userVO.SessionID = session.SessionID()
	response.Code = 200
	response.Msg = "Login success"
	userVO.Password = ""
	response.Data = map[string]interface{}{
		"username":    username,
		"displayName": displayName,
		"photoURL":    photoURL,
	}

	return response
}
func (l *LoginService) DirectLogin(ctx *gin.Context, sessions *session.Manager) {
	user, err := fileUtils.ReadLnFile(consts.PasswordPath)
	if err != nil {
		log.Panicln("Not find password file error: " + err.Error())
	}
	username := strings.TrimSpace(strings.Split(user[0], "=")[1])
	session := sessions.Start(ctx.Writer, ctx.Request)
	session.Set("username", username)
}

func (l *LoginService) Logout(ctx *gin.Context, sessions *session.Manager) {
	sessions.Destroy(ctx.Writer, ctx.Request)
}

func (l *LoginService) ChangeUser(username, password string) {
	user, err := fileUtils.ReadLnFile(consts.PasswordPath)
	if err != nil {
		log.Panicln("Not find password file error: " + err.Error())
	}
	displayName := strings.TrimSpace(strings.Split(user[2], "=")[1])
	photoURL := strings.TrimSpace(strings.Split(user[3], "=")[1])
	fileUtils.WriterLnFile(consts.PasswordPath, []string{
		"username = " + username,
		"password = " + password,
		"displayName=" + displayName,
		"photoURL=" + photoURL,
	})
}

func (l *LoginService) ChangePassword(newPassword string) *vo.Response {

	response := &vo.Response{}
	user, err := fileUtils.ReadLnFile(consts.PasswordPath)

	if err != nil {
		log.Panicln("Not find password file error: " + err.Error())
	}
	username := strings.TrimSpace(strings.Split(user[0], "=")[1])
	displayName := strings.TrimSpace(strings.Split(user[2], "=")[1])
	photoURL := strings.TrimSpace(strings.Split(user[3], "=")[1])
	fileUtils.WriterLnFile(consts.PasswordPath, []string{
		"username = " + username,
		"password = " + newPassword,
		"displayName=" + displayName,
		"photoURL=" + photoURL,
	})

	response.Code = 200
	response.Msg = "Update user new password success"

	return response
}

func (l *LoginService) InitUserInfo(userInfo *vo.UserInfo) {
	username := "username=" + userInfo.Username
	password := "password=" + userInfo.Password
	displayName := "displayName=" + userInfo.DisplayeName
	photoURL := "photoURL=" + userInfo.PhotoURL
	fileUtils.WriterLnFile(consts.PasswordPath, []string{username, password, displayName, photoURL})
}

func (l *LoginService) IsWhiteIP(ctx *gin.Context) bool {
	WhiteAdminIP := global.Config.WhiteAdminIP
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
