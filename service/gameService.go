package service

import (
	"dst-admin-go/constant/consts"
	"dst-admin-go/constant/dst"
	"dst-admin-go/utils/dstUtils"
	"dst-admin-go/vo/level"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"

	"dst-admin-go/constant/screenKey"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/utils/shellUtils"
	"dst-admin-go/utils/systemUtils"
	"dst-admin-go/vo"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type GameService struct {
	lock sync.Mutex
	c    HomeService
}

func (g *GameService) GetLastDstVersion() int64 {

	url := "http://ver.tugos.cn/getLocalVersion"
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	s := string(body)
	veriosn, err := strconv.Atoi(s)
	if err != nil {
		veriosn = 0
	}
	return int64(veriosn)
}

func (g *GameService) GetLocalDstVersion(clusterName string) int64 {
	cluster := clusterUtils.GetCluster(clusterName)
	versionTextPath := filepath.Join(cluster.ForceInstallDir, "version.txt")
	version, err := fileUtils.ReadFile(versionTextPath)
	if err != nil {
		log.Println(err)
		return 0
	}
	version = strings.Replace(version, "\n", "", -1)
	l, err := strconv.ParseInt(version, 10, 64)
	if err != nil {
		log.Println(err)
		return 0
	}
	return l
}

func ClearScreen() bool {
	result, err := shellUtils.Shell(consts.ClearScreenCmd)
	if err != nil {
		return false
	}
	res := strings.Split(result, "\n")[0]
	return res != ""
}

func (g *GameService) UpdateGame(clusterName string) error {

	g.lock.Lock()
	defer g.lock.Unlock()

	g.stopMaster(clusterName)
	g.stopCaves(clusterName)
	updateGameCMd := dst.GetDstUpdateCmd(clusterName)
	log.Println("正在更新游戏", "cluster: ", clusterName, "command: ", updateGameCMd)
	_, err := shellUtils.Shell(updateGameCMd)
	if err != nil {
		return err
	}
	return nil
}

func (g *GameService) GetLevelStatus(clusterName, level string) bool {
	cmd := " ps -ef | grep -v grep | grep -v tail |grep '" + clusterName + "'|grep " + level + " |sed -n '1P'|awk '{print $2}' "
	result, err := shellUtils.Shell(cmd)
	if err != nil {
		return false
	}
	res := strings.Split(result, "\n")[0]
	return res != ""
}

func (g *GameService) shutdownLevel(clusterName, level string) {
	if !g.GetLevelStatus(clusterName, level) {
		return
	}

	shell := "screen -S \"" + screenKey.Key(clusterName, level) + "\" -p 0 -X stuff \"c_shutdown(true)\\n\""
	log.Println("正在shutdown世界", "cluster: ", clusterName, "level: ", level, "command: ", shell)
	_, err := shellUtils.Shell(shell)
	if err != nil {
		log.Println("shut down " + clusterName + " " + level + " error: " + err.Error())
		log.Println("shutdown 失败，将强制杀掉世界")
	}
}

/*
STOP_CAVES_CMD = "ps -ef | grep -v grep |grep '" + DST_CAVES + "' |sed -n '1P'|awk '{print $2}' |xargs kill -9"
*/
func (g *GameService) killLevel(clusterName, level string) {

	if !g.GetLevelStatus(clusterName, level) {
		return
	}
	cmd := " ps -ef | grep -v grep | grep -v tail |grep '" + clusterName + "'|grep " + level + " |sed -n '1P'|awk '{print $2}' |xargs kill -9"
	log.Println("正在kill世界", "cluster: ", clusterName, "level: ", level, "command: ", cmd)
	_, err := shellUtils.Shell(cmd)
	if err != nil {
		// TODO 强制杀掉
		log.Println("kill "+clusterName+" "+level+" error: ", err)
	}
}

func (g *GameService) LaunchLevel(clusterName, level string, bin, beta int) {

	cluster := clusterUtils.GetCluster(clusterName)
	dstInstallDir := cluster.ForceInstallDir

	var startCmd = ""

	if bin == 64 {
		startCmd = "cd " + dstInstallDir + "/bin64 ; screen -d -m -S \"" + screenKey.Key(clusterName, level) + "\"  ./dontstarve_dedicated_server_nullrenderer_x64 -console -cluster " + clusterName + " -shard " + level + "  ;"
	} else {
		startCmd = "cd " + dstInstallDir + "/bin ; screen -d -m -S \"" + screenKey.Key(clusterName, level) + "\"  ./dontstarve_dedicated_server_nullrenderer -console -cluster " + clusterName + " -shard " + level + "  ;"
	}

	log.Println("正在启动世界", "cluster: ", clusterName, "level: ", level, "command: ", startCmd)
	_, err := shellUtils.Shell(startCmd)
	if err != nil {
		log.Panicln("启动 "+clusterName+" "+level+" error,", err)
	}

}

func (g *GameService) stopMaster(clusterName string) {
	level := "Master"
	g.StopLevel(clusterName, level)
}

func (g *GameService) stopCaves(clusterName string) {

	level := "Caves"
	g.StopLevel(clusterName, level)
}

func (g *GameService) StopLevel(clusterName, level string) {
	g.shutdownLevel(clusterName, level)

	time.Sleep(3 * time.Second)

	if g.GetLevelStatus(clusterName, level) {
		var i uint8 = 1
		for {
			if g.GetLevelStatus(clusterName, level) {
				break
			}
			g.shutdownLevel(clusterName, level)
			time.Sleep(1 * time.Second)
			i++
			if i > 3 {
				break
			}
		}
	}
	g.killLevel(clusterName, level)
}

func (g *GameService) StopGame(clusterName string, opType int) {
	if opType == consts.StopGame {
		g.stopMaster(clusterName)
		g.stopCaves(clusterName)
	}

	if opType == consts.StopMaster {
		g.stopMaster(clusterName)
	}

	if opType == consts.StopCaves {
		g.stopCaves(clusterName)
	}
}

func (g *GameService) launchMaster(clusterName string, bin, beta int) {
	level := "Master"
	g.LaunchLevel(clusterName, level, bin, beta)
}

func (g *GameService) launchCaves(clusterName string, bin, beta int) {
	level := "Caves"
	g.LaunchLevel(clusterName, level, bin, beta)
}

func (g *GameService) StartGame(clusterName string, bin, beta, opType int) {
	g.lock.Lock()
	defer g.lock.Unlock()
	if opType == consts.StartGame {

		g.stopMaster(clusterName)
		g.stopCaves(clusterName)

		g.launchMaster(clusterName, bin, beta)
		g.launchCaves(clusterName, bin, beta)
	}

	if opType == consts.StartMaster {
		g.stopMaster(clusterName)
		g.launchMaster(clusterName, bin, beta)
	}

	if opType == consts.StartCaves {
		g.stopCaves(clusterName)
		g.launchCaves(clusterName, bin, beta)
	}

	ClearScreen()
}

func (g *GameService) GetClusterDashboard(clusterName string) vo.ClusterDashboardVO {
	var wg sync.WaitGroup
	wg.Add(10)

	dashboardVO := vo.NewDashboardVO(clusterName)
	go func() {
		defer wg.Done()
		dashboardVO.MasterStatus = g.GetLevelStatus(clusterName, "Master")
	}()

	go func() {
		defer wg.Done()
		dashboardVO.CavesStatus = g.GetLevelStatus(clusterName, "Caves")

	}()

	go func() {
		defer wg.Done()
		dashboardVO.HostInfo = systemUtils.GetHostInfo()
	}()

	go func() {
		defer wg.Done()
		dashboardVO.CpuInfo = systemUtils.GetCpuInfo()
	}()

	go func() {
		defer wg.Done()
		dashboardVO.MemInfo = systemUtils.GetMemInfo()
	}()

	go func() {
		defer wg.Done()
		dashboardVO.DiskInfo = systemUtils.GetDiskInfo()
	}()

	go func() {
		defer wg.Done()
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		dashboardVO.MemStates = m.Alloc / 1024
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				dashboardVO.Version = 0
			}
			wg.Done()
		}()
		dashboardVO.Version = g.GetLocalDstVersion(clusterName)
	}()

	// 获取master进程占用情况
	go func() {
		defer wg.Done()
		dashboardVO.MasterPs = g.PsAuxSpecified(clusterName, "Master")
	}()
	// 获取caves进程占用情况
	go func() {
		defer wg.Done()
		dashboardVO.CavesPs = g.PsAuxSpecified(clusterName, "Caves")
	}()

	wg.Wait()
	return *dashboardVO
}

func (g *GameService) PsAuxSpecified(clusterName, level string) *vo.DstPsVo {
	dstPsVo := vo.NewDstPsVo()
	cmd := "ps -aux | grep -v grep | grep -v tail | grep " + clusterName + "  | grep " + level + " | sed -n '2P' |awk '{print $3, $4, $5, $6}'"

	info, err := shellUtils.Shell(cmd)
	if err != nil {
		log.Println(cmd + " error: " + err.Error())
		return dstPsVo
	}
	if info == "" {
		return dstPsVo
	}

	arr := strings.Split(info, " ")
	dstPsVo.CpuUage = strings.Replace(arr[0], "\n", "", -1)
	dstPsVo.MemUage = strings.Replace(arr[1], "\n", "", -1)
	dstPsVo.VSZ = strings.Replace(arr[2], "\n", "", -1)
	dstPsVo.RSS = strings.Replace(arr[3], "\n", "", -1)

	return dstPsVo
}

func (g *GameService) GetGameConfig(ctx *gin.Context) *level.GameConfig {
	gameConfig := level.GameConfig{}
	var wg sync.WaitGroup
	wg.Add(6)
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	go func() {
		gameConfig.ClusterToken = g.c.GetClusterToken(clusterName)
		wg.Done()
	}()
	go func() {
		gameConfig.ClusterIni = g.c.GetClusterIni(clusterName)
		wg.Done()
	}()
	go func() {
		gameConfig.Adminlist = g.c.GetAdminlist(clusterName)
		wg.Done()
	}()
	go func() {
		gameConfig.Blocklist = g.c.GetBlocklist(clusterName)
		wg.Done()
	}()
	go func() {
		gameConfig.Master = g.c.GetMasterWorld(clusterName)
		wg.Done()
	}()
	go func() {
		gameConfig.Caves = g.c.GetCavesWorld(clusterName)
		wg.Done()
	}()
	wg.Wait()
	return &gameConfig
}

func (g *GameService) SaveGameConfig(ctx *gin.Context, gameConfig *level.GameConfig) {
	var wg sync.WaitGroup
	wg.Add(6)
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	go func() {
		g.c.SaveClusterToken(clusterName, gameConfig.ClusterToken)
		wg.Done()
	}()

	go func() {
		g.c.SaveClusterIni(clusterName, gameConfig.ClusterIni)
		wg.Done()
	}()

	go func() {
		// SaveAdminlist(level.Adminlist)
		wg.Done()
	}()

	go func() {
		// SaveBlocklist(level.Blocklist)
		wg.Done()
	}()

	go func() {
		g.c.SaveMasterWorld(clusterName, gameConfig.Master)
		dstUtils.DedicatedServerModsSetup(clusterName, gameConfig.Master.Modoverrides)
		dstUtils.DedicatedServerModsSetup2(clusterName, gameConfig.Caves.Modoverrides)
		wg.Done()
	}()

	go func() {
		g.c.SaveCavesWorld(clusterName, gameConfig.Caves)
		wg.Done()
	}()

	wg.Wait()
}
