package api

import (
	"dst-admin-go/config/database"
	"dst-admin-go/model"
	"dst-admin-go/vo"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PlayerLogApi struct {
}

func (l *PlayerLogApi) PlayerLogQueryPage(ctx *gin.Context) {

	//获取查询参数
	name := ctx.Query("name")
	kuId := ctx.Query("kuId")
	steamId := ctx.Query("steamId")

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	if page <= 0 {
		page = 1
	}
	if size < 0 {
		size = 10
	}

	db := database.DB

	if name, isExist := ctx.GetQuery("name"); isExist {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if kuId, isExist := ctx.GetQuery("kuId"); isExist {
		db = db.Where("ku_id LIKE ?", "%"+kuId+"%")
	}
	if steamId, isExist := ctx.GetQuery("steamId"); isExist {
		db = db.Where("steamId LIKE ?", "%"+steamId+"%")
	}

	db = db.Order("created_at desc").Limit(size).Offset((page - 1) * size)

	playerLogs := make([]model.PlayerLog, 0)

	if err := db.Find(&playerLogs).Error; err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("name:", name, "kuId", kuId, "steamId", steamId)
	var total int64
	db2 := database.DB
	if name != "" {
		db2.Model(&model.PlayerLog{}).Where("name like ?", "%"+name+"%").Count(&total)
	} else {
		db2.Model(&model.PlayerLog{}).Count(&total)
	}
	totalPages := total / int64(size)
	if total%int64(size) != 0 {
		totalPages++
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: vo.Page{
			Data:       playerLogs,
			Page:       page,
			Size:       size,
			Total:      total,
			TotalPages: totalPages,
		},
	})

}
