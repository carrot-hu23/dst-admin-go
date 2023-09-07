package api

import (
	"dst-admin-go/config/database"
	"dst-admin-go/model"
	"dst-admin-go/vo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type WebLinkApi struct {
}

func (w *WebLinkApi) GetWebLinkList(ctx *gin.Context) {

	db := database.DB
	var webLinkList []model.WebLink
	db.Find(&webLinkList)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: webLinkList,
	})
}

func (w *WebLinkApi) AddWebLink(ctx *gin.Context) {

	// cluster := clusterUtils.GetClusterFromGin(ctx)
	var webLink model.WebLink
	if err := ctx.ShouldBindJSON(&webLink); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}

	db := database.DB
	db.Create(&webLink)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (w *WebLinkApi) DeleteWebLink(ctx *gin.Context) {

	id, _ := strconv.Atoi(ctx.DefaultQuery("ID", "0"))
	db := database.DB
	db.Where("ID = ?", id).Delete(&model.WebLink{})

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}
