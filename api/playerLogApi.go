package api

import (
	"dst-admin-go/config/database"
	"dst-admin-go/model"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/utils/systemUtils"
	"dst-admin-go/vo"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type PlayerLogApi struct {
}

func (l *PlayerLogApi) PlayerLogQueryPage(ctx *gin.Context) {

	//获取查询参数
	//name := ctx.Query("name")
	//kuId := ctx.Query("kuId")
	//steamId := ctx.Query("steamId")
	//role := ctx.Query("role")
	//action := ctx.Query("action")

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	if page <= 0 {
		page = 1
	}
	if size < 0 {
		size = 10
	}

	db := database.DB
	db2 := database.DB
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
	db = db.Where("cluster_name = ?", clusterName)
	db2 = db2.Where("cluster_name = ?", clusterName)
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

	var wg sync.WaitGroup
	wg.Add(len(playerLogs))
	for i := range playerLogs {
		go func(index int) {
			defer func() {
				wg.Done()
				if r := recover(); r != nil {
					log.Println(r)
				}
			}()
			ipStr := playerLogs[index].Ip
			split := strings.Split(ipStr, "|")
			if len(split) == 2 {
				ip := split[0]
				region := systemUtils.FindIpRegion(ip)
				playerLogs[index].Region = region
			}
		}(i)
	}
	wg.Wait()

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

func (l *PlayerLogApi) DeletePlayerLog(ctx *gin.Context) {
	var payload struct {
		Ids []int64 `json:"ids"`
	}

	err := ctx.ShouldBind(&payload)
	if err != nil {
		log.Panicln(err)
	}
	db := database.DB

	db.Delete(&model.PlayerLog{}, payload.Ids)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}
