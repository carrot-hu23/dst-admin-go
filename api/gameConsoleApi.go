package api

import (
	"dst-admin-go/constant/dst"
	"dst-admin-go/model"
	"dst-admin-go/service"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GameConsoleApi struct {
}

var consoleService = service.GameConsoleService{}
var gameArchiveService = service.GameArchive{}
var announceService = service.AnnounceService{}

func (g *GameConsoleApi) SentBroadcast(ctx *gin.Context) {
	message := ctx.Query("message")
	log.Println("发送公告信息：" + message)
	cluster := clusterUtils.GetClusterFromGin(ctx)

	consoleService.SentBroadcast(cluster.ClusterName, message)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (g *GameConsoleApi) KickPlayer(ctx *gin.Context) {

	kuId := ctx.Query("kuId")
	log.Println("踢出玩家：" + kuId)

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	consoleService.KickPlayer(clusterName, kuId)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (g *GameConsoleApi) KillPlayer(ctx *gin.Context) {

	kuId := ctx.Query("kuId")
	log.Println("kill玩家：" + kuId)

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	consoleService.KillPlayer(clusterName, kuId)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (g *GameConsoleApi) RespawnPlayer(ctx *gin.Context) {

	kuId := ctx.Query("kuId")
	log.Println("复活玩家：" + kuId)

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	consoleService.RespawnPlayer(clusterName, kuId)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (g *GameConsoleApi) RollBack(ctx *gin.Context) {

	dayNums := ctx.Query("dayNums")
	days, err := strconv.Atoi(dayNums)
	if err != nil {
		log.Panicln("参数解析错误：" + dayNums)
	}
	log.Println("回滚指定的天数：" + dayNums)

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	consoleService.RollBack(clusterName, days)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (g *GameConsoleApi) Regenerateworld(ctx *gin.Context) {

	log.Println("重置世界......")

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	consoleService.Regenerateworld(clusterName)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (g *GameConsoleApi) CleanWorld(ctx *gin.Context) {

	log.Println("删除世界......")

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	consoleService.CleanWorld(clusterName)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (g *GameConsoleApi) CleanLevel(ctx *gin.Context) {

	defer func() {
		if r := recover(); r != nil {
			log.Panicln("删除世界存档失败", r)
		}
	}()

	log.Println("删除世界level......")

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	basePath := dst.GetClusterBasePath(clusterName)
	levels := ctx.QueryArray("level")

	for _, level := range levels {
		fileUtils.DeleteDir(path.Join(basePath, level, "backup"))
		fileUtils.DeleteDir(path.Join(basePath, level, "save"))
		fileUtils.DeleteDir(path.Join(basePath, level, "server_chat_log.txt"))
		fileUtils.DeleteDir(path.Join(basePath, level, "server_log.txt"))
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (g *GameConsoleApi) MasterConsole(ctx *gin.Context) {
	var body struct {
		Command string `json:"command"`
	}
	if err := ctx.BindJSON(&body); err != nil {
		return
	}

	comment := body.Command

	log.Println("地面控制台: " + comment)

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	consoleService.MasterConsole(clusterName, comment)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (g *GameConsoleApi) CavesConsole(ctx *gin.Context) {
	var body struct {
		Command string `json:"command"`
	}
	if err := ctx.BindJSON(&body); err != nil {
		return
	}

	comment := body.Command

	log.Println("洞穴控制台: " + comment)

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	consoleService.CavesConsole(clusterName, comment)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (g *GameConsoleApi) OperatePlayer(ctx *gin.Context) {

	otype := ctx.Param("type")
	kuId := ctx.Param("kuId")

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	consoleService.OperatePlayer(clusterName, otype, kuId)

	log.Printf("执行高级针对玩家的操作: type=%s,kuId=%s", otype, kuId)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (g *GameConsoleApi) RestoreBackup(ctx *gin.Context) {

	backupName := ctx.Query("backupName")

	backupService.RestoreBackup(ctx, backupName)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "restore backup success",
		Data: nil,
	})
}

func (g *GameConsoleApi) GetGameArchive(ctx *gin.Context) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: gameArchiveService.GetGameArchive(clusterName),
	})
}

func (g *GameConsoleApi) GetAnnounceSetting(ctx *gin.Context) {

	//cluster := clusterUtils.GetClusterFromGin(ctx)
	//clusterName := cluster.ClusterName

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: announceService.GetAnnounceSetting(),
	})
}

func (g *GameConsoleApi) SaveAnnounceSetting(ctx *gin.Context) {

	//cluster := clusterUtils.GetClusterFromGin(ctx)
	//clusterName := cluster.ClusterName

	var announce model.Announce
	err := ctx.ShouldBind(&announce)
	log.Println(announce)
	if err != nil {
		log.Println("参数错误", err)
	}
	announceService.SaveAnnounceSetting(&announce)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: announce,
	})

}

func (g *GameConsoleApi) ReadLevelServeLog(ctx *gin.Context) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	lines, _ := strconv.Atoi(ctx.DefaultQuery("lines", "100"))
	levelName := ctx.Query("levelName")
	if levelName == "" {
		log.Panicln("levelName 不能为空")
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: consoleService.ReadLevelServerLog(clusterName, levelName, uint(lines)),
	})
}

func (g *GameConsoleApi) ReadLevelServeChatLog(ctx *gin.Context) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	lines, _ := strconv.Atoi(ctx.DefaultQuery("lines", "100"))
	levelName := ctx.Query("levelName")
	if levelName == "" {
		log.Panicln("levelName 不能为空")
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: consoleService.ReadLevelServerChatLog(clusterName, levelName, uint(lines)),
	})
}

func (g *GameConsoleApi) ReadServerLog(ctx *gin.Context) {

	lines, _ := strconv.Atoi(ctx.DefaultQuery("lines", "100"))
	logs, err := fileUtils.ReverseRead("./log.log", uint(lines))
	if err != nil {
		log.Panicln("读取面板日志失败")
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: logs,
	})
}

func (g *GameConsoleApi) DownloadDstLogFile(ctx *gin.Context) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	fileName := ctx.Query("fileName")
	if fileName == "" {
		log.Panicln("fileName 不能为空")
	}
	levelName := ctx.Query("levelName")
	if fileName == "" {
		log.Panicln("levelName 不能为空")
	}

	filePath := filepath.Join(dst.GetClusterBasePath(clusterName), levelName, fileName)
	//打开文件
	_, err := os.Open(filePath)
	//非空处理
	if err != nil {
		log.Panicln("download filePath error", err)
	}
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment; filename="+fileName)
	ctx.Header("Content-Transfer-Encoding", "binary")
	// c.Header("Content-Length", strconv.FormatInt(f.Size(), 10))
	ctx.File(filePath)
}

func (g *GameConsoleApi) DownloadServerLogFile(ctx *gin.Context) {

	filePath := "./log.log"
	//打开文件
	_, err := os.Open(filePath)
	//非空处理
	if err != nil {
		log.Panicln("download filePath error", err)
	}
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment; filename="+"log.log")
	ctx.Header("Content-Transfer-Encoding", "binary")
	// c.Header("Content-Length", strconv.FormatInt(f.Size(), 10))
	ctx.File(filePath)
}
