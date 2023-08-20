package api

import (
	"dst-admin-go/service"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/vo"
	"dst-admin-go/vo/level"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"
)

var homeService = service.HomeService{}

type Game8LevelApi struct {
}

type LevelStatus struct {
	MasterStatus bool `json:"masterStatus"`
	Slave1Status bool `json:"slave1Status"`
	Slave2Status bool `json:"slave2Status"`
	Slave3Status bool `json:"slave3Status"`
	Slave4Status bool `json:"slave4Status"`
	Slave5Status bool `json:"slave5Status"`
	Slave6Status bool `json:"slave6Status"`
	Slave7Status bool `json:"slave7Status"`

	MasterPs *vo.DstPsVo `json:"masterPs"`
	Slave1ps *vo.DstPsVo `json:"slave1ps"`
	Slave2ps *vo.DstPsVo `json:"slave2ps"`
	Slave3ps *vo.DstPsVo `json:"slave3ps"`
	Slave4ps *vo.DstPsVo `json:"slave4ps"`
	Slave5ps *vo.DstPsVo `json:"slave5ps"`
	Slave6ps *vo.DstPsVo `json:"slave6ps"`
	Slave7ps *vo.DstPsVo `json:"slave7ps"`
}

type LevelConfig struct {
	Master *level.World `json:"master"`
	Slave1 *level.World `json:"slave1"`
	Slave2 *level.World `json:"slave2"`
	Slave3 *level.World `json:"slave3"`
	Slave4 *level.World `json:"slave4"`
	Slave5 *level.World `json:"slave5"`
	Slave6 *level.World `json:"slave6"`
	Slave7 *level.World `json:"slave7"`
}

type ClusterIni struct {
	Cluster *level.ClusterIni `json:"cluster"`
	Token   string            `json:"token"`
}

// GetStatus 获取世界状态
func (g *Game8LevelApi) GetStatus(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	levelStatus := LevelStatus{}
	var wg sync.WaitGroup
	wg.Add(8)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
			wg.Done()
		}()
		levelStatus.MasterStatus = gameService.GetLevelStatus(clusterName, "Master")
		levelStatus.MasterPs = gameService.PsAuxSpecified(clusterName, "Master")
	}()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
			wg.Done()
		}()
		levelStatus.Slave1Status = gameService.GetLevelStatus(clusterName, "Slave1")
		levelStatus.Slave1ps = gameService.PsAuxSpecified(clusterName, "Slave1")
	}()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
			wg.Done()
		}()
		levelStatus.Slave2Status = gameService.GetLevelStatus(clusterName, "Slave2")
		levelStatus.Slave2ps = gameService.PsAuxSpecified(clusterName, "Slave2")
	}()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
			wg.Done()
		}()
		levelStatus.Slave3Status = gameService.GetLevelStatus(clusterName, "Slave3")
		levelStatus.Slave3ps = gameService.PsAuxSpecified(clusterName, "Slave3")
	}()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
			wg.Done()
		}()
		levelStatus.Slave4Status = gameService.GetLevelStatus(clusterName, "Slave4")
		levelStatus.Slave4ps = gameService.PsAuxSpecified(clusterName, "Slave4")
	}()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
			wg.Done()
		}()
		levelStatus.Slave5Status = gameService.GetLevelStatus(clusterName, "Slave5")
		levelStatus.Slave5ps = gameService.PsAuxSpecified(clusterName, "Slave5")
	}()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
			wg.Done()
		}()
		levelStatus.Slave6Status = gameService.GetLevelStatus(clusterName, "Slave6")
		levelStatus.Slave6ps = gameService.PsAuxSpecified(clusterName, "Slave6")
	}()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
			wg.Done()
		}()
		levelStatus.Slave7Status = gameService.GetLevelStatus(clusterName, "Slave7")
		levelStatus.Slave7ps = gameService.PsAuxSpecified(clusterName, "Slave7")
	}()

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: levelStatus,
	})
}

// Start 启动世界
func (g *Game8LevelApi) Start(ctx *gin.Context) {
	levelName := ctx.Query("levelName")
	cluster := clusterUtils.GetClusterFromGin(ctx)
	bin := cluster.Bin
	beta := cluster.Beta
	clusterName := cluster.ClusterName
	gameService.LaunchLevel(clusterName, levelName, bin, beta)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "start " + clusterName + " " + levelName + " success",
		Data: nil,
	})
}

// Stop 停止世界
func (g *Game8LevelApi) Stop(ctx *gin.Context) {
	levelName := ctx.Query("levelName")
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	gameService.StopLevel(clusterName, levelName)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "start " + clusterName + " " + levelName + " success",
		Data: nil,
	})
}

// GetClusterIni 发送房间配置
func (g *Game8LevelApi) GetClusterIni(ctx *gin.Context) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	clusterIni := ClusterIni{}
	clusterIni.Cluster = homeService.GetClusterIni(clusterName)
	clusterIni.Token = homeService.GetClusterToken(clusterName)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: &clusterIni,
	})
}

// SaveClusterIni 保存房间配置
func (g *Game8LevelApi) SaveClusterIni(ctx *gin.Context) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	clusterIni := ClusterIni{}
	err := ctx.ShouldBind(&clusterIni)
	if err != nil {
		log.Panicln("参数解析错误", err)
	}
	homeService.SaveClusterIni(clusterName, clusterIni.Cluster)
	homeService.SaveClusterToken(clusterName, clusterIni.Token)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: &clusterIni,
	})
}

// SendCommand 发送世界指令
func (g *Game8LevelApi) SendCommand(ctx *gin.Context) {
	levelName := ctx.Query("levelName")
	command := ctx.Query("command")
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	consoleService.SendCommand(clusterName, levelName, command)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

// GetLevelConfig 获取世界配置
func (g *Game8LevelApi) GetLevelConfig(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	levelConfig := LevelConfig{}
	var wg sync.WaitGroup
	wg.Add(8)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
			wg.Done()
			levelConfig.Master = homeService.GetLevel(clusterName, "Master")
		}()
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
			wg.Done()
			levelConfig.Slave1 = homeService.GetLevel(clusterName, "Slave1")
		}()
	}()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
			wg.Done()
			levelConfig.Slave2 = homeService.GetLevel(clusterName, "Slave2")
		}()
	}()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
			wg.Done()
			levelConfig.Slave3 = homeService.GetLevel(clusterName, "Slave3")
		}()
	}()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
			wg.Done()
			levelConfig.Slave4 = homeService.GetLevel(clusterName, "Slave4")
		}()
	}()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
			wg.Done()
			levelConfig.Slave5 = homeService.GetLevel(clusterName, "Slave5")
		}()
	}()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
			wg.Done()
			levelConfig.Slave6 = homeService.GetLevel(clusterName, "Slave6")
		}()
	}()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
			wg.Done()
			levelConfig.Slave7 = homeService.GetLevel(clusterName, "Slave7")
		}()
	}()
	wg.Wait()
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: &levelConfig,
	})
}

// SaveLevelConfig 保存世界配置
func (g *Game8LevelApi) SaveLevelConfig(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	levelConfig := LevelConfig{}
	ctx.ShouldBind(&levelConfig)
	log.Println("正在保存 levelConfig", levelConfig)

	homeService.SaveLevel(clusterName, "Master", levelConfig.Master)
	homeService.SaveLevel(clusterName, "Slave1", levelConfig.Slave1)
	homeService.SaveLevel(clusterName, "Slave2", levelConfig.Slave2)
	homeService.SaveLevel(clusterName, "Slave3", levelConfig.Slave3)
	homeService.SaveLevel(clusterName, "Slave4", levelConfig.Slave4)
	homeService.SaveLevel(clusterName, "Slave5", levelConfig.Slave5)
	homeService.SaveLevel(clusterName, "Slave6", levelConfig.Slave6)
	homeService.SaveLevel(clusterName, "Slave7", levelConfig.Slave7)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

// GetOnlinePlayers 获取在线玩家
func (g *Game8LevelApi) GetOnlinePlayers(ctx *gin.Context) {
	levelName := ctx.Query("levelName")
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	playerList := playerService.GetPlayerList(clusterName, levelName)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: playerList,
	})
}

// GetAdministrators 获取管理员
func (g *Game8LevelApi) GetAdministrators(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	adminList := playerService.GetDstAdminList(clusterName)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: adminList,
	})
}

// GetWhitelist 获取白名单
func (g *Game8LevelApi) GetWhitelist(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	whitelist := playerService.GetDstWhitelistPlayerList(clusterName)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: whitelist,
	})
}

// GetBlacklist 获取黑名单
func (g *Game8LevelApi) GetBlacklist(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	blacklist := playerService.GetDstBlacklistPlayerList(clusterName)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: blacklist,
	})
}

// SaveBlacklist 保存黑名单
func (g *Game8LevelApi) SaveBlacklist(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	blacklistVO := vo.NewBlacklistVO()
	err := ctx.BindJSON(blacklistVO)
	if err != nil {
		log.Panicln("参数解析错误")
	}

	playerService.SaveBlacklist(clusterName, blacklistVO.Blacklist)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

// SaveWhitelist 保存白名单
func (g *Game8LevelApi) SaveWhitelist(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	whitelist := vo.NewWhitelistVO()
	err := ctx.BindJSON(whitelist)
	if err != nil {
		log.Panicln("参数解析错误")
	}
	playerService.SaveWhitelist(clusterName, whitelist.Whitelist)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

// SaveAdminlist 保存管理员
func (g *Game8LevelApi) SaveAdminlist(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	adminlist := vo.NewAdminListVO()
	err := ctx.BindJSON(adminlist)
	if err != nil {
		log.Panicln("参数解析错误")
	}

	playerService.SaveAdminlist(clusterName, adminlist.AdminList)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}
