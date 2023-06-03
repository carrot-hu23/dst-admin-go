package api

import (
	"dst-admin-go/service"
	"dst-admin-go/vo"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GameConsoleApi struct {
}

var gameService = service.GameService{}
var dstService = service.DstService{}

func (g *GameConsoleApi) UpdateGame(ctx *gin.Context) {

	gameService.UpdateGame()

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "update dst success",
		Data: nil,
	})
}

func (g *GameConsoleApi) SentBroadcast(ctx *gin.Context) {
	message := ctx.Query("message")
	log.Println("发送公告信息：" + message)
	service.SentBroadcast(message)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (g *GameConsoleApi) KickPlayer(ctx *gin.Context) {
	kuId := ctx.Query("kuId")
	log.Println("踢出玩家：" + kuId)
	service.KickPlayer(kuId)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (g *GameConsoleApi) KillPlayer(ctx *gin.Context) {
	kuId := ctx.Query("kuId")
	log.Println("kill玩家：" + kuId)
	service.KillPlayer(kuId)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (g *GameConsoleApi) RespawnPlayer(ctx *gin.Context) {
	kuId := ctx.Query("kuId")
	log.Println("复活玩家：" + kuId)
	service.RespawnPlayer(kuId)

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
	service.RollBack(days)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (g *GameConsoleApi) Regenerateworld(ctx *gin.Context) {

	log.Println("重置世界......")
	service.Regenerateworld()

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (g *GameConsoleApi) CleanWorld(ctx *gin.Context) {

	log.Println("删除世界......")
	service.CleanWorld()

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
	service.MasterConsole(comment)
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
	service.CavesConsole(comment)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (g *GameConsoleApi) OperatePlayer(ctx *gin.Context) {

	otype := ctx.Param("type")
	kuId := ctx.Param("kuId")
	service.OperatePlayer(otype, kuId)

	log.Printf("执行高级针对玩家的操作: type=%s,kuId=%s", otype, kuId)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

// TODO GET /game/restore
func (g *GameConsoleApi) RestoreBackup(ctx *gin.Context) {

	backupName := ctx.Query("backupName")

	backupService.RestoreBackup(backupName)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "restore backup success",
		Data: nil,
	})
}

func DeleteGame() {

}

func (g *GameConsoleApi) GetGameArchive(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: dstService.GetCurrGameArchive(),
	})
}
