package api

import (
	"dst-admin-go/config/database"
	"dst-admin-go/config/dockerClient"
	"dst-admin-go/model"
	"dst-admin-go/service"
	"dst-admin-go/vo"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

var zoneService service.ZoneInfoService

type ZoneApi struct {
}

func (c *ZoneApi) GetZone(ctx *gin.Context) {
	db := database.DB
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: zoneService.FindAll(db),
	})
}

func (c *ZoneApi) CreateZone(ctx *gin.Context) {

	zone := model.ZoneInfo{}
	err := ctx.ShouldBind(&zone)
	if err != nil {
		log.Panicln(err)
	}
	if zone.ZoneCode == "" {
		log.Panicln("zoneCode 不能为空")
	}
	if zone.Name == "" {
		log.Panicln("Name 不能为空")
	}
	if zone.Ip == "" {
		log.Panicln("Ip 不能为空")
	}
	if zone.Port == 0 {
		log.Panicln("Port 不能为空")
	}
	db := database.DB
	tx := db.Begin()
	err = zoneService.Create(tx, zone)
	if err != nil {
		log.Panicln(err)
	}
	dockerClient.AddZone(zone)
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

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (c *ZoneApi) UpdateZone(ctx *gin.Context) {

	zone := model.ZoneInfo{}
	err := ctx.ShouldBind(&zone)
	if err != nil {
		log.Panicln(err)
	}

	db := database.DB
	err = zoneService.UpdateZone(db, zone.ID, zone.Name, zone.Ip, zone.Port)
	if err != nil {
		log.Panicln(err)
	}
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (c *ZoneApi) DeleteZone(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Query("id"))
	db := database.DB
	err = zoneService.Delete(db, uint(id))
	if err != nil {
		log.Panicln(err)
	}
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}
