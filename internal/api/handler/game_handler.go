package handler

import (
	"dst-admin-go/internal/pkg/context"
	"dst-admin-go/internal/pkg/response"
	"dst-admin-go/internal/pkg/utils/fileUtils"
	"dst-admin-go/internal/pkg/utils/systemUtils"
	"dst-admin-go/internal/service/archive"
	"dst-admin-go/internal/service/game"
	"dst-admin-go/internal/service/gameArchive"
	"dst-admin-go/internal/service/level"
	"dst-admin-go/internal/service/levelConfig"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type OperationInfo struct {
	State     string    `json:"state"`
	Message   string    `json:"message"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type GameHandler struct {
	process          game.Process
	level            *level.LevelService
	gameArchive      *gameArchive.GameArchive
	levelConfigUtils *levelConfig.LevelConfigUtils
	archive          *archive.PathResolver
	opMu             sync.Mutex
	operations       map[string]OperationInfo
}

func NewGameHandler(process game.Process, levelService *level.LevelService, gameArchive *gameArchive.GameArchive, levelConfigUtils *levelConfig.LevelConfigUtils, archive *archive.PathResolver) *GameHandler {
	return &GameHandler{
		process:          process,
		level:            levelService,
		gameArchive:      gameArchive,
		levelConfigUtils: levelConfigUtils,
		archive:          archive,
		operations:       make(map[string]OperationInfo),
	}
}

func (p *GameHandler) operationKey(clusterName, levelName string) string {
	return clusterName + "/" + levelName
}

func (p *GameHandler) setOperation(clusterName, levelName, state, message string) {
	p.opMu.Lock()
	defer p.opMu.Unlock()
	p.operations[p.operationKey(clusterName, levelName)] = OperationInfo{State: state, Message: message, UpdatedAt: time.Now()}
}

func (p *GameHandler) clearOperation(clusterName, levelName string) {
	p.opMu.Lock()
	defer p.opMu.Unlock()
	delete(p.operations, p.operationKey(clusterName, levelName))
}

func (p *GameHandler) getOperation(clusterName, levelName string) OperationInfo {
	p.opMu.Lock()
	defer p.opMu.Unlock()
	return p.operations[p.operationKey(clusterName, levelName)]
}

func (p *GameHandler) reconcileOperation(clusterName, levelName string, status bool) OperationInfo {
	op := p.getOperation(clusterName, levelName)
	switch op.State {
	case "starting":
		if status {
			p.clearOperation(clusterName, levelName)
			return OperationInfo{}
		}
	case "stopping":
		if !status {
			p.clearOperation(clusterName, levelName)
			return OperationInfo{}
		}
	}
	return op
}

func (p *GameHandler) RegisterRoute(router *gin.RouterGroup) {
	router.GET("/api/game/8level/start", p.Start)
	router.GET("/api/game/8level/stop", p.Stop)
	router.GET("/api/game/8level/start/all", p.StartAll)
	router.GET("/api/game/8level/stop/all", p.StopAll)
	router.POST("/api/game/8level/command", p.Command)
	router.GET("/api/game/8level/status", p.Status)
	router.GET("/api/game/archive", p.GameArchive)
	router.GET("/api/game/system/info", p.SystemInfo)
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

	status, err := p.process.Status(clusterName, levelName)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Response{Code: 500, Msg: "failed to read game server status: " + err.Error()})
		return
	}
	if !status {
		p.clearOperation(clusterName, levelName)
		ctx.JSON(http.StatusOK, response.Response{Code: 200, Msg: "already stopped"})
		return
	}

	p.setOperation(clusterName, levelName, "stopping", "stop command accepted")
	go func(clusterName, levelName string) {
		if err := p.process.Stop(clusterName, levelName); err != nil {
			p.setOperation(clusterName, levelName, "failed", "stop failed: "+err.Error())
			log.Printf("async stop failed: cluster=%s level=%s err=%v", clusterName, levelName, err)
			return
		}
		if !p.waitForLevelStatus(clusterName, levelName, false, 60*time.Second) {
			p.setOperation(clusterName, levelName, "failed", "stop timeout: process is still running")
			return
		}
		p.clearOperation(clusterName, levelName)
	}(clusterName, levelName)
	ctx.JSON(http.StatusOK, response.Response{Code: 200, Msg: "accepted"})
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

	status, err := p.process.Status(clusterName, levelName)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Response{Code: 500, Msg: "failed to read game server status: " + err.Error()})
		return
	}
	if status {
		p.clearOperation(clusterName, levelName)
		ctx.JSON(http.StatusOK, response.Response{Code: 200, Msg: "already running"})
		return
	}

	// Write ASCII-only markers when the panel accepts the start request.
	// DST may truncate server_log.txt during process boot, so write once immediately
	// for instant feedback and repeat shortly after launch so the marker survives
	// across zh/en/ja/ko UI locales and all language-specific level display names.
	p.writeLogMarker(clusterName, levelName, "START_REQUESTED")
	p.writeLogMarkerDelayed(clusterName, levelName, "START_REQUESTED", 1*time.Second)
	p.writeLogMarkerDelayed(clusterName, levelName, "START_REQUESTED", 3*time.Second)
	p.writeLogMarkerDelayed(clusterName, levelName, "START_REQUESTED", 8*time.Second)
	p.setOperation(clusterName, levelName, "starting", "start command accepted")
	go func(clusterName, levelName string) {
		// Do not let slow version/config preparation create a fake-start window with no DST output.
		// The user's visible feedback must be tied to the real launcher first.
		if err := p.process.Start(clusterName, levelName); err != nil {
			p.setOperation(clusterName, levelName, "failed", "start failed: "+err.Error())
			log.Printf("async start failed: cluster=%s level=%s err=%v", clusterName, levelName, err)
			return
		}
		go func() {
			if err := p.prepareSingleStart(clusterName, levelName); err != nil {
				log.Printf("async start prepare failed: cluster=%s level=%s err=%v", clusterName, levelName, err)
			}
		}()
		if !p.waitForLevelStatus(clusterName, levelName, true, 180*time.Second) {
			p.setOperation(clusterName, levelName, "failed", "start timeout: process did not stay running")
			return
		}
		p.clearOperation(clusterName, levelName)
	}(clusterName, levelName)
	ctx.JSON(http.StatusOK, response.Response{Code: 200, Msg: "accepted"})
}

func (p *GameHandler) prepareSingleStart(clusterName, levelName string) error {
	if p.archive == nil || p.levelConfigUtils == nil {
		return nil
	}
	version, _ := p.archive.GetLocalDstVersion(clusterName)
	config, err := p.levelConfigUtils.GetLevelConfig(clusterName)
	if err == nil {
		for i := range config.LevelList {
			if config.LevelList[i].File == levelName {
				config.LevelList[i].RunVersion = version
				config.LevelList[i].Version = version
				break
			}
		}
		if saveErr := p.levelConfigUtils.SaveLevelConfig(clusterName, config); saveErr != nil {
			log.Printf("async start save level config failed: cluster=%s level=%s err=%v", clusterName, levelName, saveErr)
		}
	} else {
		log.Printf("async start get level config failed: cluster=%s level=%s err=%v", clusterName, levelName, err)
	}
	levelPath := p.archive.LevelPath(clusterName, levelName)
	customPath := filepath.Join(levelPath, "customcommands.lua")
	_ = fileUtils.CreateFileIfNotExists(customPath)
	customcommands, readErr := fileUtils.ReadFile("./static/customcommands.lua")
	if readErr != nil {
		return readErr
	}
	fileUtils.WriterTXT(customPath, customcommands)
	return nil
}

func (p *GameHandler) waitForLevelStatus(clusterName, levelName string, expected bool, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for {
		status, err := p.process.Status(clusterName, levelName)
		if err == nil && status == expected {
			return true
		}
		if time.Now().After(deadline) {
			return false
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func (p *GameHandler) writeLogMarkerDelayed(clusterName, levelName, action string, delay time.Duration) {
	go func() {
		time.Sleep(delay)
		p.writeLogMarker(clusterName, levelName, action)
	}()
}

func (p *GameHandler) writeLogMarkerAfterStatus(clusterName, levelName, successAction, timeoutAction string, expected bool, timeout time.Duration) {
	go func() {
		if !p.waitForLevelStatus(clusterName, levelName, expected, timeout) {
			p.writeLogMarker(clusterName, levelName, timeoutAction)
			return
		}
		p.writeLogMarker(clusterName, levelName, successAction)
		p.writeLogMarkerDelayed(clusterName, levelName, successAction, 3*time.Second)
		if expected {
			p.writeLogMarkerDelayed(clusterName, levelName, successAction, 10*time.Second)
		}
	}()
}

func (p *GameHandler) writeLogMarker(clusterName, levelName, action string) {
	if clusterName == "" || levelName == "" {
		log.Printf("skip dst log marker: missing cluster/level action=%s cluster=%q level=%q", action, clusterName, levelName)
		return
	}
	logPath := ""
	if p.archive != nil {
		logPath = p.archive.ServerLogPath(clusterName, levelName)
	}
	if logPath == "" {
		// Fallback for production systemd environments; do not silently lose markers.
		logPath = fmt.Sprintf("/root/.klei/DoNotStarveTogether/%s/%s/server_log.txt", clusterName, levelName)
	}
	elapsed := ""
	if p.process != nil {
		ps := p.process.PsAuxSpecified(clusterName, levelName)
		elapsed = ps.Elapsed
	}
	if elapsed == "" {
		elapsed = "00:00:00"
	}
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		log.Printf("failed to create dst log marker directory: cluster=%s level=%s action=%s path=%s err=%v", clusterName, levelName, action, logPath, err)
		return
	}
	now := time.Now()
	line := fmt.Sprintf("\n========== HERMES_LOG_MARKER action=%s cluster=%s level=%s wall=%s elapsed=%s =========="+
		"\n", action, clusterName, levelName, now.Format("2006-01-02 15:04:05"), elapsed)
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("failed to write dst log marker: cluster=%s level=%s action=%s path=%s err=%v", clusterName, levelName, action, logPath, err)
		return
	}
	defer f.Close()
	if _, err := f.WriteString(line); err != nil {
		log.Printf("failed to append dst log marker: cluster=%s level=%s action=%s path=%s err=%v", clusterName, levelName, action, logPath, err)
		return
	}
	log.Printf("wrote dst log marker: cluster=%s level=%s action=%s path=%s elapsed=%s", clusterName, levelName, action, logPath, elapsed)
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
	if p.levelConfigUtils == nil {
		ctx.JSON(http.StatusOK, response.Response{Code: 500, Msg: "level config service is unavailable"})
		return
	}
	config, err := p.levelConfigUtils.GetLevelConfig(clusterName)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Response{Code: 500, Msg: "failed to read level config: " + err.Error()})
		return
	}

	for i := range config.LevelList {
		levelName := config.LevelList[i].File
		status, statusErr := p.process.Status(clusterName, levelName)
		if statusErr != nil {
			log.Printf("stop all status check failed: cluster=%s level=%s err=%v", clusterName, levelName, statusErr)
			continue
		}
		if !status {
			p.clearOperation(clusterName, levelName)
			continue
		}
		p.writeLogMarker(clusterName, levelName, "STOP_REQUESTED")
		p.setOperation(clusterName, levelName, "stopping", "stop all command accepted")
	}

	go func(clusterName string, levels []string) {
		if err := p.process.StopAll(clusterName); err != nil {
			log.Printf("async stop all failed: cluster=%s err=%v", clusterName, err)
			for _, levelName := range levels {
				p.setOperation(clusterName, levelName, "failed", "stop all failed: "+err.Error())
			}
			return
		}
		for _, levelName := range levels {
			if !p.waitForLevelStatus(clusterName, levelName, false, 60*time.Second) {
				p.setOperation(clusterName, levelName, "failed", "stop timeout: process is still running")
				continue
			}
			p.clearOperation(clusterName, levelName)
		}
	}(clusterName, func() []string {
		levels := make([]string, 0, len(config.LevelList))
		for i := range config.LevelList {
			levels = append(levels, config.LevelList[i].File)
		}
		return levels
	}())

	ctx.JSON(http.StatusOK, response.Response{Code: 200, Msg: "accepted"})
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
	Ps                 game.DstPsAux         `json:"Ps"`
	RunVersion         int64                 `json:"runVersion"`
	Status             bool                  `json:"status"`
	IsMaster           bool                  `json:"isMaster"`
	LevelName          string                `json:"levelName"`
	Uuid               string                `json:"uuid"`
	Leveldataoverride  string                `json:"leveldataoverride"`
	Modoverrides       string                `json:"modoverrides"`
	ServerIni          levelConfig.ServerIni `json:"serverIni"`
	OperationState     string                `json:"operationState"`
	OperationMessage   string                `json:"operationMessage"`
	OperationUpdatedAt string                `json:"operationUpdatedAt"`
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
	ctx.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	ctx.Header("Pragma", "no-cache")
	ctx.Header("Expires", "0")

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
				op := p.reconcileOperation(clusterName, levelItem.Uuid, status)
				result[index] = LevelStatus{
					Ps:                 ps,
					Status:             status,
					RunVersion:         levelItem.RunVersion,
					LevelName:          levelItem.LevelName,
					IsMaster:           levelItem.IsMaster || levelItem.Uuid == "Master",
					Uuid:               levelItem.Uuid,
					Leveldataoverride:  levelItem.Leveldataoverride,
					Modoverrides:       levelItem.Modoverrides,
					ServerIni:          levelItem.ServerIni,
					OperationState:     op.State,
					OperationMessage:   op.Message,
					OperationUpdatedAt: op.UpdatedAt.Format(time.RFC3339),
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
			ps := p.process.PsAuxSpecified(clusterName, levelItem.Uuid)
			status, _ := p.process.Status(clusterName, levelItem.Uuid)
			op := p.reconcileOperation(clusterName, levelItem.Uuid, status)
			result[i] = LevelStatus{
				Ps:                 ps,
				Status:             status,
				RunVersion:         levelItem.RunVersion,
				LevelName:          levelItem.LevelName,
				IsMaster:           levelItem.IsMaster || levelItem.Uuid == "Master",
				Uuid:               levelItem.Uuid,
				Leveldataoverride:  levelItem.Leveldataoverride,
				Modoverrides:       levelItem.Modoverrides,
				ServerIni:          levelItem.ServerIni,
				OperationState:     op.State,
				OperationMessage:   op.Message,
				OperationUpdatedAt: op.UpdatedAt.Format(time.RFC3339),
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

// SystemInfo 获取服务器系统信息
// @Summary 系统信息
// @Description 获取服务器系统信息
// @Tags game
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/game/system/info [get]
func (p *GameHandler) SystemInfo(ctx *gin.Context) {
	clusterName := context.GetClusterName(ctx)
	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "success",
		Data: p.GetSystemInfo(clusterName),
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
