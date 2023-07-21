package api

import (
	"dst-admin-go/config/database"
	"dst-admin-go/constant/consts"
	"dst-admin-go/model"
	"dst-admin-go/vo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync"
)

type AutoCheckApi struct {
}

const (
	TURN = 1
	OFF  = 0
)

var lock = sync.Mutex{}

func (m *AutoCheckApi) EnableAutoCheckUpdateVersion(ctx *gin.Context) {
	lock.Lock()
	defer lock.Unlock()
	enable, _ := strconv.Atoi(ctx.DefaultQuery("enable", "0"))
	db := database.DB
	autoCheck := model.AutoCheck{}
	db.Where("name = ?", consts.UpdateGameVersion).Find(&autoCheck)
	autoCheck.Enable = enable
	autoCheck.Name = consts.UpdateGameVersion
	db.Save(&autoCheck)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (m *AutoCheckApi) EnableAutoCheckRun(ctx *gin.Context) {
	lock.Lock()
	defer lock.Unlock()

	enable, _ := strconv.Atoi(ctx.DefaultQuery("enable", "0"))
	db := database.DB
	autoCheck := model.AutoCheck{}
	db.Where("name = ?", consts.GameRunning).Find(&autoCheck)
	autoCheck.Name = consts.GameRunning
	autoCheck.Enable = enable

	db.Save(&autoCheck)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (m *AutoCheckApi) EnableAutoCheckGameMod(ctx *gin.Context) {
	lock.Lock()
	defer lock.Unlock()
	enable, _ := strconv.Atoi(ctx.DefaultQuery("enable", "0"))
	db := database.DB
	autoCheck := model.AutoCheck{}
	db.Where("name = ?", consts.UpdateGameMod).Find(&autoCheck)
	autoCheck.Enable = enable
	autoCheck.Name = consts.UpdateGameVersion
	db.Save(&autoCheck)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (m *AutoCheckApi) GetAutoCheckStatus(ctx *gin.Context) {

	db1 := database.DB
	autoCheck1 := model.AutoCheck{}
	db1.Where("name = ?", consts.GameRunning).Find(&autoCheck1)

	db2 := database.DB
	autoCheck2 := model.AutoCheck{}
	db2.Where("name = ?", consts.UpdateGameVersion).Find(&autoCheck2)

	db3 := database.DB
	autoCheck3 := model.AutoCheck{}
	db3.Where("name = ?", consts.UpdateGameMod).Find(&autoCheck3)

	res := map[string]int{}
	res[consts.GameRunning] = autoCheck1.Enable
	res[consts.UpdateGameVersion] = autoCheck2.Enable
	res[consts.UpdateGameMod] = autoCheck3.Enable

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: res,
	})
}
