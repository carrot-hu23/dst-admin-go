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

type ClusterIni struct {
	Cluster *level.ClusterIni `json:"cluster"`
	Token   string            `json:"token"`
}

type LevelInfo struct {
	Ps                *vo.DstPsVo      `json:"Ps"`
	Status            bool             `json:"status"`
	LevelName         string           `json:"levelName"`
	IsMaster          bool             `json:"is_master"`
	Uuid              string           `json:"uuid"`
	Leveldataoverride string           `json:"leveldataoverride"`
	Modoverrides      string           `json:"modoverrides"`
	ServerIni         *level.ServerIni `json:"server_ini"`
}

// GetStatus 获取世界状态
func (g *Game8LevelApi) GetStatus(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	levelList := gameLevel2Service.GetLevelList(clusterName)
	length := len(levelList)
	result := make([]LevelInfo, length)
	var wg sync.WaitGroup
	wg.Add(length)
	for i := range levelList {
		go func(index int) {
			defer func() {
				wg.Done()
				if r := recover(); r != nil {

				}
			}()
			world := levelList[index]
			ps := gameService.PsAuxSpecified(clusterName, world.Uuid)
			status := gameService.GetLevelStatus(clusterName, world.Uuid)
			result[index] = LevelInfo{
				Ps:                ps,
				Status:            status,
				LevelName:         world.LevelName,
				IsMaster:          world.IsMaster,
				Uuid:              world.Uuid,
				Leveldataoverride: world.Leveldataoverride,
				Modoverrides:      world.Modoverrides,
				ServerIni:         world.ServerIni,
			}
		}(i)
	}
	wg.Wait()
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: result,
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
