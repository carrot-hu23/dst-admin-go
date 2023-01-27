package api

import (
	"dst-admin-go/service"
	"dst-admin-go/vo"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitUser(ctx *gin.Context) {

	userVO := vo.NewUserVO()
	ctx.BindJSON(userVO)

	service.Inituser(userVO)

	userVO.Password = ""
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "Init user success",
		Data: userVO,
	})
}

func Login(ctx *gin.Context) {

	userVO := vo.NewUserVO()
	ctx.BindJSON(userVO)

	response := service.Login(userVO, ctx, sessions)
	ctx.JSON(http.StatusOK, response)
}

func Logout(ctx *gin.Context) {
	service.Logout(ctx, sessions)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "Logout success",
	})
}

func ChangePassword(ctx *gin.Context) {

	var body struct {
		NewPassword string `json:"newPassword"`
	}
	if err := ctx.BindJSON(&body); err != nil {
		return
	}
	newPassword := body.NewPassword
	response := service.ChangePassword(newPassword)

	ctx.JSON(http.StatusOK, response)
}
