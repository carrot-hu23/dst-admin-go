package api

import (
	"dst-admin-go/config/database"
	"dst-admin-go/config/dockerClient"
	"dst-admin-go/model"
	"dst-admin-go/service"
	"dst-admin-go/vo"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

var queueInfoService service.QueueInfoService

type QueueApi struct {
}

func (c *QueueApi) GetQueue(ctx *gin.Context) {
	db := database.DB
	infos, err := queueInfoService.FindAll(db)
	if err != nil {
		log.Panicln(err)
	}
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: infos,
	})
}

func (c *QueueApi) CreateQueue(ctx *gin.Context) {

	checkAdmin(ctx)

	queue := model.QueueInfo{}
	err := ctx.ShouldBind(&queue)
	if err != nil {
		log.Panicln(err)
	}
	if queue.QueueCode == "" {
		log.Panicln("queueCode 不能为空")
	}
	if queue.Name == "" {
		log.Panicln("Name 不能为空")
	}
	if queue.Ip == "" {
		log.Panicln("Ip 不能为空")
	}
	if queue.Port == 0 {
		log.Panicln("Port 不能为空")
	}
	db := database.DB
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			ctx.JSON(http.StatusOK, vo.Response{
				Code: 500,
				Msg:  "创建失败",
				Data: nil,
			})
		}
	}()

	err = queueInfoService.Create(tx, queue)
	if err != nil {
		log.Panicln(err)
	}
	err = dockerClient.AddQueue(queue)
	if err != nil {
		log.Panicln(err)
	}
	tx.Commit()

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (c *QueueApi) UpdateQueue(ctx *gin.Context) {

	checkAdmin(ctx)

	queue := model.QueueInfo{}
	err := ctx.ShouldBind(&queue)
	if err != nil {
		log.Panicln(err)
	}

	db := database.DB
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			ctx.JSON(http.StatusOK, vo.Response{
				Code: 500,
				Msg:  "更新失败",
				Data: nil,
			})
		}
	}()
	queueInfo := model.QueueInfo{}
	if err = tx.First(&queueInfo, queue.ID).Error; err != nil {
		log.Panicln(err)
	}

	err = queueInfoService.UpdateQueueInfo(tx, queue.ID, queue.Name, queue.Ip, queue.Port)
	if err != nil {
		log.Panicln(err)
	}

	err = dockerClient.UpdateQueue(queueInfo)
	if err != nil {
		log.Panicln(err)
	}
	tx.Commit()

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (c *QueueApi) DeleteQueue(ctx *gin.Context) {

	checkAdmin(ctx)

	id, err := strconv.Atoi(ctx.Query("id"))
	db := database.DB
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			ctx.JSON(http.StatusOK, vo.Response{
				Code: 500,
				Msg:  "删除失败",
				Data: nil,
			})
		}
	}()
	queueInfo := model.QueueInfo{}
	if err = tx.First(&queueInfo, id).Error; err != nil {
		log.Panicln(err)
	}
	err = queueInfoService.Delete(tx, uint(id))
	if err != nil {
		log.Panicln(err)
	}
	dockerClient.DeleteQueue(queueInfo.QueueCode)

	tx.Commit()

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (c *QueueApi) BindQueue2Zone(ctx *gin.Context) {

	checkAdmin(ctx)

	var payload struct {
		ZoneCode  string `json:"ZoneCode"`
		QueueCode string `json:"queueCode"`
	}
	err := ctx.ShouldBind(&payload)
	if err != nil {
		log.Panicln(err)
	}
	if payload.ZoneCode == "" {
		log.Panicln("zoneCode 不能为空")
	}
	if payload.QueueCode == "" {
		log.Panicln("queueCode 不能为空")
	}
	db := database.DB
	// 查找 Zone
	var zoneInfo model.ZoneInfo
	if err := db.Where("zone_code = ?", payload.ZoneCode).First(&zoneInfo).Error; err != nil {
		log.Panicln(err)
	}

	var queueInfo model.QueueInfo
	// 查找 Queue
	if err := db.Where("queue_code = ?", payload.QueueCode).First(&queueInfo).Error; err != nil {
		log.Panicln(err)
	}

	zoneQueue := model.ZoneQueue{
		ZoneCode:  payload.ZoneCode,
		QueueCode: payload.QueueCode,
	}

	db.Create(&zoneQueue)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (c *QueueApi) UnbindQueueFromZone(ctx *gin.Context) {
	checkAdmin(ctx)
	var payload struct {
		ZoneCode  string `json:"ZoneCode"`
		QueueCode string `json:"queueCode"`
	}
	err := ctx.ShouldBind(&payload)
	if err != nil {
		log.Panicln(err)
	}

	db := database.DB
	// 查找绑定关系
	var zoneQueue model.ZoneQueue
	if err := db.Where("zone_code = ? AND queue_code = ?", payload.ZoneCode, payload.QueueCode).First(&zoneQueue).Error; err != nil {
		log.Panicln(errors.New("未找到当前绑定关系"))
	}
	// 删除绑定关系
	if err := db.Unscoped().Delete(&zoneQueue).Error; err != nil {
		log.Panicln(err)
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (c *QueueApi) GetQueuesByZone(ctx *gin.Context) {
	zoneCode := ctx.Query("zoneCode")
	db := database.DB
	// 查询绑定的 QueueCodes
	var zoneQueues []model.ZoneQueue
	if err := db.Where("zone_code = ?", zoneCode).Find(&zoneQueues).Error; err != nil {
		log.Panicln(err)
	}

	// 获取所有的 QueueInfo
	var queueCodes []string
	for _, zq := range zoneQueues {
		queueCodes = append(queueCodes, zq.QueueCode)
	}

	var queueInfos []model.QueueInfo
	if err := db.Where("queue_code IN (?)", queueCodes).Find(&queueInfos).Error; err != nil {
		log.Panicln(err)
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: queueInfos,
	})

}
