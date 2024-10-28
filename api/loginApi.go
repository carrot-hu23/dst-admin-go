package api

import (
	"dst-admin-go/constant"
	"dst-admin-go/service"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type LoginApi struct {
}

var loginService = service.LoginService{}

func (l *LoginApi) GetUserInfo(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "Init user success",
		Data: loginService.GetUserInfo(),
	})
}

func (l *LoginApi) Login(ctx *gin.Context) {

	userVO := vo.NewUserVO()
	ctx.BindJSON(userVO)

	response := loginService.Login(userVO, ctx)
	ctx.JSON(http.StatusOK, response)
}

func (l *LoginApi) Logout(ctx *gin.Context) {
	loginService.Logout(ctx)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "Logout success",
	})
}

func (l *LoginApi) ChangePassword(ctx *gin.Context) {

	var body struct {
		NewPassword string `json:"newPassword"`
	}
	if err := ctx.BindJSON(&body); err != nil {
		return
	}
	newPassword := body.NewPassword
	response := loginService.ChangePassword(newPassword)

	ctx.JSON(http.StatusOK, response)
}

func (l *LoginApi) UpdateUserInfo(ctx *gin.Context) {

	var body struct {
		Username    string `json:"username"`
		DisplayName string `json:"displayName"`
		PhotoURL    string `json:"photoURL"`
		Password    string `json:"password"`
	}
	if err := ctx.BindJSON(&body); err != nil {
		log.Panicln("参数解析错误: " + err.Error())
		return
	}
	err := fileUtils.WriterLnFile(constant.PASSWORD_PATH, []string{
		"username = " + body.Username,
		"password = " + body.Password,
		"displayName=" + body.DisplayName,
		"photoURL=" + body.PhotoURL,
	})
	if err != nil {
		log.Panicln("修改用户信息失败: " + err.Error())
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "Logout success",
		Data: nil,
	})
}
