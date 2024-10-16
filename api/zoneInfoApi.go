package api

import (
	"dst-admin-go/config/database"
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
	infos, err := zoneService.FindAll(db)
	if err != nil {
		log.Panicln(err)
	}
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: infos,
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

	err = zoneService.Create(tx, zone)
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

func (c *ZoneApi) UpdateZone(ctx *gin.Context) {

	zone := model.ZoneInfo{}
	err := ctx.ShouldBind(&zone)
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
	zoneInfo := model.ZoneInfo{}
	if err = tx.First(&zoneInfo, zone.ID).Error; err != nil {
		log.Panicln(err)
	}

	err = zoneService.UpdateZone(tx, zone.ID, zone.Name)
	if err != nil {
		log.Panicln(err)
	}

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

func (c *ZoneApi) DeleteZone(ctx *gin.Context) {

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
	zoneInfo := model.ZoneInfo{}
	if err = tx.First(&zoneInfo, id).Error; err != nil {
		log.Panicln(err)
	}
	err = zoneService.Delete(tx, uint(id))
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
