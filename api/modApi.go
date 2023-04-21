package api

import (
	"dst-admin-go/mod"
	"dst-admin-go/vo"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SearchModList(ctx *gin.Context) {

	//获取查询参数
	text := ctx.Query("text")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	data, err := mod.SearchModList(text, page, size)
	if err != nil {
		log.Panicln("搜索mod失败", err)
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: data,
	})
}

func GetModInfo(ctx *gin.Context) {

	moId := ctx.Param("modId")
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: mod.GetModInfo(moId),
	})
}
