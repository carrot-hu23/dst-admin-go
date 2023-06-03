package api

import (
	"dst-admin-go/service"
	"dst-admin-go/vo"
	"net/http"

	"github.com/gin-gonic/gin"
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

	response := loginService.Login(userVO, ctx, sessions)
	ctx.JSON(http.StatusOK, response)
}

func (l *LoginApi) Logout(ctx *gin.Context) {
	loginService.Logout(ctx, sessions)
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
