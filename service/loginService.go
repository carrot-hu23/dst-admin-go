package service

import (
	"dst-admin-go/config/database"
	"dst-admin-go/model"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo"
	"github.com/gin-contrib/sessions"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

type LoginService struct {
}

func (l *LoginService) GetUserInfo() map[string]interface{} {
	//user, err := fileUtils.ReadLnFile("./password.txt")
	//
	//if err != nil {
	//	log.Panicln("Not find password file error: " + err.Error())
	//}
	//
	//username := strings.TrimSpace(strings.Split(user[0], "=")[1])
	//// password := strings.TrimSpace(strings.Split(user[1], "=")[1])
	//displayName := strings.TrimSpace(strings.Split(user[2], "=")[1])
	//photoURL := strings.TrimSpace(strings.Split(user[3], "=")[1])
	//
	//return map[string]interface{}{
	//	"username":    username,
	//	"displayName": displayName,
	//	"photoURL":    photoURL,
	//}
	return nil
}

func (l *LoginService) Login(userVO *vo.UserVO, ctx *gin.Context) *vo.Response {

	session := sessions.Default(ctx)

	response := &vo.Response{}
	db := database.DB
	dbUser := model.User{}
	db.Where("username=?", userVO.Username).Find(&dbUser)
	if dbUser.Password != userVO.Password {
		log.Panicln("User authentication failed")
		response.Code = 401
		response.Msg = "User authentication failed"
		return response
	}
	session.Set("username", dbUser.Username)
	session.Set("role", dbUser.Role)
	// TODO 增加集群限制权限
	session.Set("cluster", []string{})
	session.Set("userId", dbUser.ID)
	err := session.Save()
	if err != nil {
		log.Panicln(err)
	}
	userVO.SessionID = session.ID()
	response.Code = 200
	response.Msg = "Login success"
	userVO.Password = ""
	response.Data = map[string]interface{}{
		"username":    dbUser.Username,
		"displayName": dbUser.DisplayName,
		"photoURL":    dbUser.PhotoURL,
		"role":        dbUser.Role,
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

func (l *LoginService) ChangeUser(username, password string) {
	user, err := fileUtils.ReadLnFile("./password.txt")
	if err != nil {
		log.Panicln("Not find password file error: " + err.Error())
	}
	displayName := strings.TrimSpace(strings.Split(user[2], "=")[1])
	photoURL := strings.TrimSpace(strings.Split(user[3], "=")[1])
	fileUtils.WriterLnFile("./password.txt", []string{
		"username = " + username,
		"password = " + password,
		"displayName=" + displayName,
		"photoURL=" + photoURL,
	})
}

func (l *LoginService) ChangePassword(newPassword string) *vo.Response {

	response := &vo.Response{}
	user, err := fileUtils.ReadLnFile("./password.txt")

	if err != nil {
		log.Panicln("Not find password file error: " + err.Error())
	}
	username := strings.TrimSpace(strings.Split(user[0], "=")[1])
	displayName := strings.TrimSpace(strings.Split(user[2], "=")[1])
	photoURL := strings.TrimSpace(strings.Split(user[3], "=")[1])
	fileUtils.WriterLnFile("./password.txt", []string{
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
	fileUtils.WriterLnFile("./password.txt", []string{username, password, displayName, photoURL})
}
