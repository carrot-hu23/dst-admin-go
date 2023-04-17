package api

import (
	"dst-admin-go/service"
	"dst-admin-go/vo"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UpdateGame(ctx *gin.Context) {

	service.UpdateGame()

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "update dst success",
		Data: nil,
	})
}

func StartGame(ctx *gin.Context) {

	opType, _ := strconv.Atoi(ctx.DefaultQuery("type", "0"))
	log.Println("正在启动游戏服务 type:", opType)
	service.StartGame(opType)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "start game success",
		Data: nil,
	})
}

func StoptGame(ctx *gin.Context) {

	opType, _ := strconv.Atoi(ctx.DefaultQuery("type", "0"))
	log.Println("正在停止游戏服务 type:", opType)
	service.StopGame(opType)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "stop game success",
		Data: nil,
	})
}

func StartMaster(ctx *gin.Context) {

	service.StartMaster()

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "start master success",
		Data: nil,
	})
}

func StopMaster(ctx *gin.Context) {

	service.ElegantShutdownMaster()

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "stop master success",
		Data: nil,
	})
}

func StartCaves(ctx *gin.Context) {

	service.ElegantShutdownCaves()

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "start caves success",
		Data: nil,
	})
}

func StopCaves(ctx *gin.Context) {

	service.ElegantShutdownCaves()

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "stop caves success",
		Data: nil,
	})
}

func SentBroadcast(ctx *gin.Context) {
	message := ctx.Query("message")
	log.Println("发送公告信息：" + message)
	service.SentBroadcast(message)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func KickPlayer(ctx *gin.Context) {
	kuId := ctx.Query("kuId")
	log.Println("踢出玩家：" + kuId)
	service.KickPlayer(kuId)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func KillPlayer(ctx *gin.Context) {
	kuId := ctx.Query("kuId")
	log.Println("kill玩家：" + kuId)
	service.KillPlayer(kuId)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func RespawnPlayer(ctx *gin.Context) {
	kuId := ctx.Query("kuId")
	log.Println("复活玩家：" + kuId)
	service.RespawnPlayer(kuId)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func RollBack(ctx *gin.Context) {
	dayNums := ctx.Param("dayNums")
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

func Regenerateworld(ctx *gin.Context) {

	log.Println("重置世界......")
	service.Regenerateworld()

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func MasterConsole(ctx *gin.Context) {
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

func CavesConsole(ctx *gin.Context) {
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

func OperatePlayer(ctx *gin.Context) {

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

// TODO GET /game/backup
func Backup(ctx *gin.Context) {
	ctx.String(200, "test")
}

// TODO GET /game/restore
func RestoreBackup(ctx *gin.Context) {

	backupName := ctx.Query("backupName")

	service.RestoreBackup(backupName)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "restore backup success",
		Data: nil,
	})
}

func DeleteGame() {

}

func GetGameArchive(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: service.GetGameArchive(),
	})
}
