package service

import (
	"dst-admin-go/constant"
	"dst-admin-go/session"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

func Inituser(userVO *vo.UserVO) {
	username := "username=" + userVO.Username
	password := "password=" + userVO.Password
	fileUtils.WriterLnFile(constant.PASSWORD_PATH, []string{username, password})
}

func Login(userVO *vo.UserVO, ctx *gin.Context, sessions *session.Manager) *vo.Response {

	response := &vo.Response{}

	user, err := fileUtils.ReadLnFile(constant.PASSWORD_PATH)

	if err != nil {
		log.Panicln("Not find password file error: " + err.Error())
	}

	username := strings.TrimSpace(strings.Split(user[0], "=")[1])
	password := strings.TrimSpace(strings.Split(user[1], "=")[1])

	if username != userVO.Username || password != userVO.Password {
		log.Panicln("User authentication failed")
		response.Code = 401
		response.Msg = "User authentication failed"
		return response
	}

	session := sessions.Start(ctx.Writer, ctx.Request)

	session.Set("username", username)

	userVO.SessionID = session.SessionID()
	response.Code = 200
	response.Msg = "Login success"
	userVO.Password = ""
	response.Data = userVO

	return response
}

func Logout(ctx *gin.Context, sessions *session.Manager) {
	sessions.Destroy(ctx.Writer, ctx.Request)
}

func ChangeUser(username, password string) {
	fileUtils.WriterLnFile(constant.PASSWORD_PATH, []string{
		"username = " + username,
		"password = " + password,
	})
}

func ChangePassword(newPassword string) *vo.Response {

	response := &vo.Response{}
	user, err := fileUtils.ReadLnFile(constant.PASSWORD_PATH)

	if err != nil {
		log.Panicln("Not find password file error: " + err.Error())
	}
	username := strings.TrimSpace(strings.Split(user[0], "=")[1])
	//password := strings.TrimSpace(strings.Split(user[1], "=")[1])

	fileUtils.WriterLnFile(constant.PASSWORD_PATH, []string{
		"username = " + username,
		"password = " + newPassword,
	})

	response.Code = 200
	response.Msg = "Update user new password success"

	return response
}
