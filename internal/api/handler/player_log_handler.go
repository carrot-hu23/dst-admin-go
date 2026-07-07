package handler

import (
	"dst-admin-go/internal/database"
	"dst-admin-go/internal/model"
	"dst-admin-go/internal/pkg/response"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PlayerLogHandler struct {
}

func NewPlayerLogHandler() *PlayerLogHandler {
	return &PlayerLogHandler{}
}

func (l *PlayerLogHandler) RegisterRoute(router *gin.RouterGroup) {
	router.GET("/api/player/log", l.PlayerLogQueryPage)
	router.POST("/api/player/log/delete", l.DeletePlayerLog)
	router.GET("/api/player/log/delete/all", l.DeletePlayerLogAll)
}

// PlayerLogQueryPage 生成 swagger 文档注释
// @Summary 分页查询玩家日志
// @Description 分页查询玩家日志
// @Tags playerLog
// @Param name query string false "玩家名称"
// @Param kuId query string false "KuId"
// @Param steamId query string false "SteamId"
// @Param role query string false "角色"
// @Param action query string false "操作"
// @Param ip query string false "IP地址"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} response.Response
// @Router /api/player/log [get]
func (l *PlayerLogHandler) PlayerLogQueryPage(ctx *gin.Context) {

	//获取查询参数
	//name := ctx.Query("name")
	//kuId := ctx.Query("kuId")
	//steamId := ctx.Query("steamId")
	//role := ctx.Query("role")
	//action := ctx.Query("action")

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	if page <= 0 {
		page = 1
	}
	if size < 0 {
		size = 10
	}

	db := database.Db
	db2 := database.Db
	if name, isExist := ctx.GetQuery("name"); isExist {
		db = db.Where("name LIKE ?", "%"+name+"%")
		db2 = db2.Where("name LIKE ?", "%"+name+"%")
	}
	if kuId, isExist := ctx.GetQuery("kuId"); isExist {
		db = db.Where("ku_id LIKE ?", "%"+kuId+"%")
		db2 = db2.Where("ku_id LIKE ?", "%"+kuId+"%")
	}
	if steamId, isExist := ctx.GetQuery("steamId"); isExist {
		db = db.Where("steamId LIKE ?", "%"+steamId+"%")
		db2 = db2.Where("steamId LIKE ?", "%"+steamId+"%")
	}
	if role, isExist := ctx.GetQuery("role"); isExist {
		db = db.Where("role LIKE ?", "%"+role+"%")
		db2 = db2.Where("role LIKE ?", "%"+role+"%")
	}
	if action, isExist := ctx.GetQuery("action"); isExist {
		db = db.Where("action LIKE ?", "%"+action+"%")
		db2 = db2.Where("action LIKE ?", "%"+action+"%")
	}
	if ip, isExist := ctx.GetQuery("ip"); isExist {
		db = db.Where("ip LIKE ?", "%"+ip+"%")
		db2 = db2.Where("ip LIKE ?", "%"+ip+"%")
	}

	db = db.Order("created_at desc").Limit(size).Offset((page - 1) * size)

	playerLogs := make([]model.PlayerLog, 0)

	if err := db.Find(&playerLogs).Error; err != nil {
		fmt.Println(err.Error())
	}

	var total int64
	db2.Model(&model.PlayerLog{}).Count(&total)

	totalPages := total / int64(size)
	if total%int64(size) != 0 {
		totalPages++
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "success",
		Data: response.Page{
			Data:       playerLogs,
			Page:       page,
			Size:       size,
			Total:      total,
			TotalPages: totalPages,
		},
	})

}

// DeletePlayerLog 删除玩家日志
// @Summary 删除玩家日志
// @Description 删除玩家日志
// @Tags playerLog
// @Param ids body []int64 true "ID列表"
// @Success 200 {object} response.Response
// @Router /api/player/log/delete [post]
func (l *PlayerLogHandler) DeletePlayerLog(ctx *gin.Context) {
	var payload struct {
		Ids []int64 `json:"ids"`
	}
	err := ctx.ShouldBind(&payload)
	if err != nil {
		log.Panicln(err)
	}
	db := database.Db

	db.Delete(&model.PlayerLog{}, payload.Ids)

	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

// DeletePlayerLogAll 删除所有玩家日志
// @Summary 删除所有玩家日志
// @Description 删除所有玩家
// @Tags playerLog
// @Success 200 {object} response.Response
// @Router /api/player/log/delete/all [get]
func (l *PlayerLogHandler) DeletePlayerLogAll(ctx *gin.Context) {

	db := database.Db
	db.Delete(&model.PlayerLog{})

	// 删除所有记录
	result := db.Where("1 = 1").Delete(&model.PlayerLog{})

	if result.Error != nil {
		log.Panicln(result.Error)
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}
