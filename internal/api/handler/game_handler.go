package handler

import (
	"dst-admin-go/internal/middleware"
	"dst-admin-go/internal/pkg/context"
	"dst-admin-go/internal/pkg/response"
	"dst-admin-go/internal/pkg/utils/shellUtils"
	"dst-admin-go/internal/pkg/utils/systemUtils"
	"dst-admin-go/internal/service/archive"
	"dst-admin-go/internal/service/game"
	"dst-admin-go/internal/service/gameArchive"
	"dst-admin-go/internal/service/level"
	"dst-admin-go/internal/service/levelConfig"
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	process          game.Process
	level            *level.LevelService
	gameArchive      *gameArchive.GameArchive
	levelConfigUtils *levelConfig.LevelConfigUtils
	archive          *archive.PathResolver
}

func NewGameHandler(process game.Process, levelService *level.LevelService, gameArchive *gameArchive.GameArchive, levelConfigUtils *levelConfig.LevelConfigUtils, archive *archive.PathResolver) *GameHandler {
	return &GameHandler{
		process:          process,
		level:            levelService,
		gameArchive:      gameArchive,
		levelConfigUtils: levelConfigUtils,
		archive:          archive,
	}
}

func (p *GameHandler) RegisterRoute(router *gin.RouterGroup) {
	router.GET("/api/game/8level/start", p.Start, middleware.StartBeforeMiddleware(p.archive, p.levelConfigUtils))
	router.GET("/api/game/8level/stop", p.Stop)
	router.GET("/api/game/8level/start/all", p.StartAll, middleware.StartBeforeMiddleware(p.archive, p.levelConfigUtils))
	router.GET("/api/game/8level/stop/all", p.StopAll)
	router.POST("/api/game/8level/command", p.Command)
	router.GET("/api/game/8level/status", p.Status)
	router.GET("/api/game/archive", p.GameArchive)
	router.GET("/api/game/system/info/stream", p.SystemInfoStream)
}

// Stop 停止世界 swagger 注释
// @Summary 停止世界
// @Description 停止世界
// @Tags game
// @Accept json
// @Produce json
// @Param levelName query string true "世界名称"
// @Success 200 {object} response.Response
// @Router /api/game/stop [get]
func (p *GameHandler) Stop(ctx *gin.Context) {

	clusterName := context.GetClusterName(ctx)
	levelName := ctx.Query("levelName")
	if levelName == "" {
		ctx.JSON(400, response.Response{Code: 400, Msg: "levelName query parameter is required"})
		return
	}
	err := p.process.Stop(clusterName, levelName)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Response{Code: 500, Msg: "failed to stop game server: " + err.Error()})
	} else {
		ctx.JSON(http.StatusOK, response.Response{Code: 200, Msg: "success"})
	}
}

// Start 启动世界 swagger 注释
// @Summary 启动世界
// @Description 启动世界
// @Tags game
// @Accept json
// @Produce json
// @Param levelName query string true "世界名称"
// @Success 200 {object} response.Response
// @Router /api/game/start [get]
func (p *GameHandler) Start(ctx *gin.Context) {
	clusterName := context.GetClusterName(ctx)
	levelName := ctx.Query("levelName")
	if levelName == "" {
		ctx.JSON(400, response.Response{Code: 400, Msg: "levelName query parameter is required"})
		return
	}
	err := p.process.Start(clusterName, levelName)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Response{Code: 500, Msg: "failed to start game server: " + err.Error()})
	} else {
		ctx.JSON(http.StatusOK, response.Response{Code: 200, Msg: "success"})
	}
}

// StartAll 启动所有世界 swagger 注释
// @Summary 启动所有世界
// @Description 启动所有世界
// @Tags game
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/game/start/all [get]
func (p *GameHandler) StartAll(ctx *gin.Context) {
	clusterName := context.GetClusterName(ctx)
	err := p.process.StartAll(clusterName)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Response{Code: 500, Msg: "failed to start all game servers: " + err.Error()})
	} else {
		ctx.JSON(http.StatusOK, response.Response{Code: 200, Msg: "success"})
	}
}

// StopAll 停止所有世界 swagger 注释
// @Summary 停止所有世界
// @Description 停止所有世界
// @Tags game
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/game/stop/all [get]
func (p *GameHandler) StopAll(ctx *gin.Context) {
	clusterName := context.GetClusterName(ctx)
	err := p.process.StopAll(clusterName)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Response{Code: 500, Msg: "failed to stop all game servers: " + err.Error()})
	} else {
		ctx.JSON(http.StatusOK, response.Response{Code: 200, Msg: "success"})
	}
}

// Command 运行命令 swagger 注释
// @Summary 运行命令
// @Description 运行命令
// @Tags game
// @Accept json
// @Produce json
// @Param command query string true "命令"
// @Success 200 {object} response.Response
// @Router /api/game/command [post]
func (p *GameHandler) Command(ctx *gin.Context) {
	clusterName := context.GetClusterName(ctx)
	type payload struct {
		Command   string `json:"command"`
		LevelName string `json:"levelName"`
	}
	var command payload
	if err := ctx.ShouldBindJSON(&command); err != nil {
		ctx.JSON(400, response.Response{Code: 400, Msg: "Invalid request body"})
		return
	}
	if command.LevelName == "" {
		ctx.JSON(400, response.Response{Code: 400, Msg: "levelName query parameter is required"})
		return
	}
	status, err := p.process.Status(clusterName, command.LevelName)
	if !status {
		ctx.JSON(http.StatusOK, response.Response{Code: 200, Msg: "game server is not running"})
		return
	}
	err = p.process.Command(clusterName, command.LevelName, command.Command)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Response{Code: 500, Msg: "failed to run command: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response.Response{Code: 200, Msg: "success"})
}

type LevelStatus struct {
	Ps                game.DstPsAux         `json:"Ps"`
	RunVersion        int64                 `json:"runVersion"`
	Status            bool                  `json:"status"`
	IsMaster          bool                  `json:"isMaster"`
	LevelName         string                `json:"levelName"`
	Uuid              string                `json:"uuid"`
	Leveldataoverride string                `json:"leveldataoverride"`
	Modoverrides      string                `json:"modoverrides"`
	ServerIni         levelConfig.ServerIni `json:"serverIni"`
}

// Status 获取服务器状态
// @Summary 获取服务器状态
// @Description 获取所有世界的运行状态信息
// @Tags game
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]LevelStatus}
// @Router /api/game/8level/status [get]
func (p *GameHandler) Status(ctx *gin.Context) {

	clusterName := context.GetClusterName(ctx)
	levelList := p.level.GetLevelList(clusterName)
	length := len(levelList)
	result := make([]LevelStatus, length)

	if runtime.GOOS == "windows" {
		var wg sync.WaitGroup
		wg.Add(length)
		for i := range levelList {
			go func(index int) {
				defer func() {
					wg.Done()
					if r := recover(); r != nil {

					}
				}()
				levelItem := levelList[index]
				ps := p.process.PsAuxSpecified(clusterName, levelItem.Uuid)
				status, _ := p.process.Status(clusterName, levelItem.Uuid)
				result[index] = LevelStatus{
					Ps:                ps,
					Status:            status,
					RunVersion:        levelItem.RunVersion,
					LevelName:         levelItem.LevelName,
					IsMaster:          levelItem.IsMaster,
					Uuid:              levelItem.Uuid,
					Leveldataoverride: levelItem.Leveldataoverride,
					Modoverrides:      levelItem.Modoverrides,
					ServerIni:         levelItem.ServerIni,
				}
			}(i)
		}
		wg.Wait()
		ctx.JSON(http.StatusOK, response.Response{
			Code: 200,
			Msg:  "success",
			Data: result,
		})
	} else {
		for i := range levelList {
			levelItem := levelList[i]
			result[i] = LevelStatus{
				Status:            false,
				RunVersion:        levelItem.RunVersion,
				LevelName:         levelItem.LevelName,
				IsMaster:          levelItem.IsMaster,
				Uuid:              levelItem.Uuid,
				Leveldataoverride: levelItem.Leveldataoverride,
				Modoverrides:      levelItem.Modoverrides,
				ServerIni:         levelItem.ServerIni,
			}
		}

		cmd := "ps -aux | grep -v grep | grep -v tail | grep -v SCREEN | grep " + clusterName + " |awk '{print $3, $4, $5, $6,$16}'"
		info, err := shellUtils.Shell(cmd)
		if err != nil {
			log.Println(cmd + " error: " + err.Error())
		} else {
			lines := strings.Split(info, "\n")
			for lineIndex := range lines {
				dstPsVo := game.DstPsAux{}
				arr := strings.Split(lines[lineIndex], " ")
				if len(arr) > 4 {
					dstPsVo.CpuUage = strings.Replace(arr[0], "\n", "", -1)
					dstPsVo.MemUage = strings.Replace(arr[1], "\n", "", -1)
					dstPsVo.VSZ = strings.Replace(arr[2], "\n", "", -1)
					dstPsVo.RSS = strings.Replace(arr[3], "\n", "", -1)
					for i := range result {
						levelName := result[i].Uuid
						if strings.Contains(arr[4], levelName) {
							result[i].Ps = dstPsVo
							result[i].Status = true
						}
					}
				}

			}
		}
		ctx.JSON(http.StatusOK, response.Response{
			Code: 200,
			Msg:  "success",
			Data: result,
		})
	}
}

// GameArchive 获取游戏存档列表
// @Summary 获取游戏存档列表
// @Description 获取当前集群的游戏存档列表
// @Tags game
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{gameArchive.GameArchiveInfo}
// @Router /api/game/archive [get]
func (p *GameHandler) GameArchive(ctx *gin.Context) {
	clusterName := context.GetClusterName(ctx)
	archiveInfo := p.gameArchive.GetGameArchive(clusterName)
	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "success",
		Data: archiveInfo,
	})
}

// SystemInfoStream 系统信息流
// @Summary 系统信息流
// @Description 获取服务器系统信息的实时流 (SSE)
// @Tags game
// @Accept text/event-stream
// @Produce text/event-stream
// @Success 200 {string} string "SSE 格式的系统信息流"
// @Router /api/game/system/info/stream [get]
func (p *GameHandler) SystemInfoStream(ctx *gin.Context) {
	clusterName := context.GetClusterName(ctx)

	// 设置SSE响应头
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("X-Accel-Buffering", "no")

	// 创建一个ticker,每2秒推送一次数据
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// 使用context来检测客户端断开连接
	clientGone := ctx.Request.Context().Done()

	// 立即发送第一次数据
	p.sendSystemInfoData(ctx, clusterName)

	for {
		select {
		case <-clientGone:
			log.Println("Client disconnected from system info stream")
			return
		case <-ticker.C:
			p.sendSystemInfoData(ctx, clusterName)
		}
	}
}

func (p *GameHandler) sendSystemInfoData(ctx *gin.Context, clusterName string) {
	systemInfo := p.GetSystemInfo(clusterName)

	// 构造响应数据
	response := response.Response{
		Code: 200,
		Msg:  "success",
		Data: systemInfo,
	}

	// 将数据序列化为JSON
	data, err := json.Marshal(response)
	if err != nil {
		log.Println("Failed to marshal system info data:", err)
		return
	}

	// 发送SSE数据
	ctx.SSEvent("message", string(data))
	ctx.Writer.Flush()
}

type SystemInfo struct {
	HostInfo      *systemUtils.HostInfo `json:"host"`
	CpuInfo       *systemUtils.CpuInfo  `json:"cpu"`
	MemInfo       *systemUtils.MemInfo  `json:"mem"`
	DiskInfo      *systemUtils.DiskInfo `json:"disk"`
	PanelMemUsage uint64                `json:"panelMemUsage"`
	PanelCpuUsage float64               `json:"panelCpuUsage"`
}

func (p *GameHandler) GetSystemInfo(clusterName string) *SystemInfo {
	var wg sync.WaitGroup
	wg.Add(5)

	dashboardVO := SystemInfo{}
	go func() {
		defer func() {
			wg.Done()
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		dashboardVO.HostInfo = systemUtils.GetHostInfo()
	}()

	go func() {
		defer func() {
			wg.Done()
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		dashboardVO.CpuInfo = systemUtils.GetCpuInfo()
	}()

	go func() {
		defer func() {
			wg.Done()
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		dashboardVO.MemInfo = systemUtils.GetMemInfo()
	}()

	go func() {
		defer func() {
			wg.Done()
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		dashboardVO.DiskInfo = systemUtils.GetDiskInfo()
	}()

	go func() {
		defer func() {
			wg.Done()
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		dashboardVO.PanelMemUsage = m.Alloc / 1024 // 将字节转换为MB
	}()

	wg.Wait()
	return &dashboardVO
}
