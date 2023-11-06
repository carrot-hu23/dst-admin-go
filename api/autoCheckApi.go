package api

import (
	"dst-admin-go/config/database"
	"dst-admin-go/constant/consts"
	"dst-admin-go/model"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/utils/levelConfigUtils"
	"dst-admin-go/vo"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type AutoCheckApi struct{}

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

func (m *AutoCheckApi) GetAutoCheckList2(ctx *gin.Context) {

	checkType := ctx.Query("checkType")

	cluster := clusterUtils.GetClusterFromGin(ctx)
	config, _ := levelConfigUtils.GetLevelConfig(cluster.ClusterName)

	var uuidSet []string
	var result []model.AutoCheck
	if checkType == "" {
		for i := range config.LevelList {
			level := config.LevelList[i]
			uuidSet = append(uuidSet, level.File)
			autoCheck1 := model.AutoCheck{
				ClusterName:  cluster.ClusterName,
				LevelName:    level.Name,
				Uuid:         level.File,
				Enable:       0,
				Announcement: "",
				Times:        1,
				Sleep:        5,
				Interval:     5,
				CheckType:    consts.LEVEL_MOD,
			}
			autoCheck2 := model.AutoCheck{
				ClusterName:  cluster.ClusterName,
				LevelName:    level.Name,
				Uuid:         level.File,
				Enable:       0,
				Announcement: "",
				Times:        1,
				Sleep:        5,
				Interval:     5,
				CheckType:    consts.LEVEL_DOWN,
			}
			result = append(result, autoCheck1)
			result = append(result, autoCheck2)
		}
		autoCheck3 := model.AutoCheck{
			ClusterName:  cluster.ClusterName,
			LevelName:    cluster.ClusterName,
			Uuid:         "",
			Enable:       0,
			Announcement: "",
			Times:        1,
			Sleep:        5,
			Interval:     5,
			CheckType:    consts.UPDATE_GAME,
		}
		result = append(result, autoCheck3)
	} else if checkType == consts.LEVEL_DOWN {
		for i := range config.LevelList {
			level := config.LevelList[i]
			uuidSet = append(uuidSet, level.File)
			autoCheck1 := model.AutoCheck{
				ClusterName:  cluster.ClusterName,
				LevelName:    level.Name,
				Uuid:         level.File,
				Enable:       0,
				Announcement: "",
				Times:        1,
				Sleep:        5,
				Interval:     5,
				CheckType:    consts.LEVEL_DOWN,
			}
			result = append(result, autoCheck1)
		}
	} else if checkType == consts.LEVEL_MOD {
		for i := range config.LevelList {
			level := config.LevelList[i]
			uuidSet = append(uuidSet, level.File)
			autoCheck1 := model.AutoCheck{
				ClusterName:  cluster.ClusterName,
				LevelName:    level.Name,
				Uuid:         level.File,
				Enable:       0,
				Announcement: "",
				Times:        1,
				Sleep:        5,
				Interval:     5,
				CheckType:    consts.LEVEL_MOD,
			}
			result = append(result, autoCheck1)
		}
	} else {
		autoCheck3 := model.AutoCheck{
			ClusterName:  cluster.ClusterName,
			LevelName:    "",
			Uuid:         "",
			Enable:       0,
			Announcement: "",
			Times:        1,
			Sleep:        5,
			Interval:     5,
			CheckType:    consts.UPDATE_GAME,
		}
		result = append(result, autoCheck3)
	}

	db := database.DB
	var dbAutoChecks []model.AutoCheck
	if checkType == "" {
		db.Where("uuid in ?", uuidSet).Find(&dbAutoChecks)
	} else {
		db.Where("check_type = ? and uuid in ?", checkType, uuidSet).Find(&dbAutoChecks)
	}

	for i := range result {
		for j := range dbAutoChecks {
			if result[i].Uuid == dbAutoChecks[j].Uuid && result[i].CheckType == dbAutoChecks[j].CheckType {
				result[i].Enable = dbAutoChecks[j].Enable
				result[i].Announcement = dbAutoChecks[j].Announcement
				result[i].Times = dbAutoChecks[j].Times
				result[i].Sleep = dbAutoChecks[j].Sleep
				result[i].Interval = dbAutoChecks[j].Interval
			}
		}
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: result,
	})

}

func (m *AutoCheckApi) SaveAutoCheck2(ctx *gin.Context) {
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
