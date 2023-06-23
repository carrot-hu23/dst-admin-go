package api

import (
	"dst-admin-go/config/database"
	"dst-admin-go/config/global"
	"dst-admin-go/lobbyServer"
	"dst-admin-go/utils/pageUtils"
	"dst-admin-go/vo"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LobbyServerApi struct {
}

func (l *LobbyServerApi) QueryLobbyServerList(ctx *gin.Context) {

	db := database.DB
	db2 := database.DB

	if addr, isExist := ctx.GetQuery("__addr"); isExist {
		db = db.Where("addr = ?", addr)
		db2 = db.Where("addr = ?", addr)
	}
	if rowID, isExist := ctx.GetQuery("__rowId"); isExist {
		db = db.Where("row_id = ?", rowID)
		db2 = db.Where("row_id = ?", rowID)
	}
	if host, isExist := ctx.GetQuery("host"); isExist {
		db = db.Where("host = ?", host)
		db2 = db.Where("host = ?", host)
	}
	if clanonly, isExist := ctx.GetQuery("clanonly"); isExist {
		db = db.Where("clanonly = ?", clanonly)
		db2 = db.Where("clanonly = ?", clanonly)
	}
	if platform, isExist := ctx.GetQuery("platform"); isExist {
		db = db.Where("platform = ?", platform)
		db2 = db.Where("platform = ?", platform)
	}
	if mods, isExist := ctx.GetQuery("mods"); isExist {
		db = db.Where("platform = ?", mods)
		db2 = db.Where("platform = ?", mods)
	}
	if name, isExist := ctx.GetQuery("name"); isExist {
		db = db.Where("name like ?", "%"+name+"%")
		db2 = db.Where("name like ?", "%"+name+"%")
	}
	if pvp, isExist := ctx.GetQuery("pvp"); isExist {
		db = db.Where("pvp = ?", pvp)
		db2 = db.Where("pvp = ?", pvp)
	}

	if password, isExist := ctx.GetQuery("password"); isExist {
		db = db.Where("password = ?", password)
		db2 = db.Where("password = ?", password)
	}
	if password, isExist := ctx.GetQuery("password"); isExist {
		db = db.Where("password = ?", password)
		db2 = db.Where("password = ?", password)
	}

	page, size := pageUtils.RequestPage(ctx)
	db = db.Order("created_at desc").Limit(size).Offset((page - 1) * size)

	lobbyHomeList := make([]lobbyServer.LobbyHome, 0)
	if err := db.Find(&lobbyHomeList).Error; err != nil {
		fmt.Println(err.Error())
	}

	var total int64
	db2.Model(&lobbyServer.LobbyHome{}).Count(&total)

	totalPages := total / int64(size)
	if total%int64(size) != 0 {
		totalPages++
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: vo.Page{
			Data:       lobbyHomeList,
			Page:       page,
			Size:       size,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

func (l *LobbyServerApi) QueryLobbyServerDetail(ctx *gin.Context) {

	//获取查询参数
	region := ctx.Query("region")
	rowId := ctx.Query("rowId")

	homeDetail := global.LobbyServer.QueryLobbyHomeInfo(region, rowId)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: homeDetail,
	})

}
