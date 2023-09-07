package api

import (
	"dst-admin-go/config/database"
	"dst-admin-go/constant/consts"
	"dst-admin-go/model"
	"dst-admin-go/vo"
	"github.com/gin-gonic/gin"
	"log"
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

func (m *AutoCheckApi) SaveAutoCheck(ctx *gin.Context) {
	lock.Lock()
	defer lock.Unlock()
	var autoCheck model.AutoCheck
	err := ctx.ShouldBind(&autoCheck)
	if err != nil {
		log.Panicln("参数错误", err)
	}
	db := database.DB
	db.Save(&autoCheck)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (m *AutoCheckApi) GetAutoCheck(ctx *gin.Context) {

	name := ctx.Query("name")

	db1 := database.DB
	autoCheck1 := model.AutoCheck{}
	db1.Where("name = ?", name).Find(&autoCheck1)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: autoCheck1,
	})
}

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

func (m *AutoCheckApi) EnableAutoCheckMasterRun(ctx *gin.Context) {
	lock.Lock()
	defer lock.Unlock()

	enable, _ := strconv.Atoi(ctx.DefaultQuery("enable", "0"))
	db := database.DB
	autoCheck := model.AutoCheck{}
	db.Where("name = ?", consts.MasterRunning).Find(&autoCheck)
	autoCheck.Name = consts.MasterRunning
	autoCheck.Enable = enable

	db.Save(&autoCheck)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (m *AutoCheckApi) EnableAutoCheckCavesRun(ctx *gin.Context) {
	lock.Lock()
	defer lock.Unlock()

	enable, _ := strconv.Atoi(ctx.DefaultQuery("enable", "0"))
	db := database.DB
	autoCheck := model.AutoCheck{}
	db.Where("name = ?", consts.CavesRunning).Find(&autoCheck)
	autoCheck.Name = consts.CavesRunning
	autoCheck.Enable = enable

	db.Save(&autoCheck)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (m *AutoCheckApi) EnableAutoCheckMasterMod(ctx *gin.Context) {
	lock.Lock()
	defer lock.Unlock()
	enable, _ := strconv.Atoi(ctx.DefaultQuery("enable", "0"))
	db := database.DB
	autoCheck := model.AutoCheck{}
	db.Where("name = ?", consts.UpdateMasterMod).Find(&autoCheck)
	autoCheck.Enable = enable
	autoCheck.Name = consts.UpdateMasterMod
	db.Save(&autoCheck)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (m *AutoCheckApi) EnableAutoCheckCavesMod(ctx *gin.Context) {
	lock.Lock()
	defer lock.Unlock()
	enable, _ := strconv.Atoi(ctx.DefaultQuery("enable", "0"))
	db := database.DB
	autoCheck := model.AutoCheck{}
	db.Where("name = ?", consts.UpdateCavesMod).Find(&autoCheck)
	autoCheck.Enable = enable
	autoCheck.Name = consts.UpdateCavesMod
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
	db1.Where("name = ?", consts.MasterRunning).Find(&autoCheck1)

	db2 := database.DB
	autoCheck2 := model.AutoCheck{}
	db2.Where("name = ?", consts.UpdateGameVersion).Find(&autoCheck2)

	db3 := database.DB
	autoCheck3 := model.AutoCheck{}
	db3.Where("name = ?", consts.UpdateGameMod).Find(&autoCheck3)

	db4 := database.DB
	autoCheck4 := model.AutoCheck{}
	db4.Where("name = ?", consts.CavesRunning).Find(&autoCheck4)

	db5 := database.DB
	autoCheck5 := model.AutoCheck{}
	db5.Where("name = ?", consts.UpdateMasterMod).Find(&autoCheck5)

	db6 := database.DB
	autoCheck6 := model.AutoCheck{}
	db6.Where("name = ?", consts.UpdateCavesMod).Find(&autoCheck6)

	res := map[string]int{}
	res[consts.MasterRunning] = autoCheck1.Enable
	res[consts.UpdateGameVersion] = autoCheck2.Enable
	res[consts.UpdateGameMod] = autoCheck3.Enable
	res[consts.CavesRunning] = autoCheck4.Enable
	res[consts.UpdateMasterMod] = autoCheck5.Enable
	res[consts.UpdateCavesMod] = autoCheck6.Enable

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: res,
	})
}
