package service

import (
	"dst-admin-go/constant/consts"
	"dst-admin-go/constant/dst"
	"dst-admin-go/utils/dstUtils"
	"dst-admin-go/vo/world"
	"github.com/gin-gonic/gin"

	"dst-admin-go/constant/screenKey"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/utils/shellUtils"
	"dst-admin-go/utils/systemUtils"
	"dst-admin-go/vo"
	"fmt"
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

func GetLocalDstVersion(clusterName string) string {
	cluster := clusterUtils.GetCluster(clusterName)
	versionTextPath := filepath.Join(cluster.ForceInstallDir, "version.txt")
	version, err := fileUtils.ReadFile(versionTextPath)
	if err != nil {
		return ""
	}
	return version
}

func ClearScreen() bool {
	result, err := shellUtils.Shell(consts.ClearScreenCmd)
	if err != nil {
		return false
	}
	res := strings.Split(result, "\n")[0]
	return res != ""
}

func (g *GameService) UpdateGame(clusterName string) {

	g.lock.Lock()
	defer g.lock.Unlock()

	g.stopMaster(clusterName)
	g.stopCaves(clusterName)
	updateGameCMd := dst.GetDstUpdateCmd(clusterName)
	log.Println(updateGameCMd)
	_, err := shellUtils.Shell(updateGameCMd)
	if err != nil {
		log.Panicln("update world error: " + err.Error())
	}
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
	_, err := shellUtils.Shell(cmd)
	if err != nil {
		// TODO 强制杀掉
		log.Panicln("kill " + clusterName + " " + level + " error: " + err.Error())
	}
}

func (g *GameService) launchLevel(clusterName, level string) {

	cluster := clusterUtils.GetCluster(clusterName)
	dstInstallDir := cluster.ForceInstallDir

	cmd := "cd " + dstInstallDir + "/bin ; screen -d -m -S \"" + screenKey.Key(clusterName, level) + "\"  ./dontstarve_dedicated_server_nullrenderer -console -cluster " + clusterName + " -shard " + level + "  ;"

	_, err := shellUtils.Shell(cmd)
	if err != nil {
		log.Panicln("launch " + clusterName + " " + level + " error: " + err.Error())
	}

}

func (g *GameService) stopMaster(clusterName string) {
	level := "Master"
	g.stopLevel(clusterName, level)
}

func (g *GameService) stopCaves(clusterName string) {

	level := "Caves"
	g.stopLevel(clusterName, level)
}

func (g *GameService) stopLevel(clusterName, level string) {
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

func (g *GameService) launchMaster(clusterName string) {
	level := "Master"
	g.launchLevel(clusterName, level)
}

func (g *GameService) launchCaves(clusterName string) {
	level := "Caves"
	g.launchLevel(clusterName, level)
}

func (g *GameService) StartGame(clusterName string, opType int) {
	g.lock.Lock()
	defer g.lock.Unlock()
	if opType == consts.StartGame {

		g.stopMaster(clusterName)
		g.stopCaves(clusterName)

		g.launchMaster(clusterName)
		g.launchCaves(clusterName)
	}

	if opType == consts.StartMaster {
		g.stopMaster(clusterName)
		g.launchMaster(clusterName)
	}

	if opType == consts.StartCaves {
		g.stopCaves(clusterName)
		g.launchCaves(clusterName)
	}

	ClearScreen()
}

func (g *GameService) GetClusterDashboard(clusterName string) vo.ClusterDashboardVO {
	var wg sync.WaitGroup
	wg.Add(10)

	dashboardVO := vo.NewDashboardVO(clusterName)
	start := time.Now()
	go func() {
		defer wg.Done()
		s1 := time.Now()
		dashboardVO.MasterStatus = g.GetLevelStatus(clusterName, "Master")
		elapsed := time.Since(s1)
		fmt.Println("master =", elapsed)
	}()

	go func() {
		defer wg.Done()
		s1 := time.Now()
		dashboardVO.CavesStatus = g.GetLevelStatus(clusterName, "Caves")
		elapsed := time.Since(s1)
		fmt.Println("cave =", elapsed)
	}()

	go func() {
		defer wg.Done()
		s1 := time.Now()
		dashboardVO.HostInfo = systemUtils.GetHostInfo()
		elapsed := time.Since(s1)
		fmt.Println("host =", elapsed)
	}()

	go func() {
		defer wg.Done()
		s1 := time.Now()
		dashboardVO.CpuInfo = systemUtils.GetCpuInfo()
		elapsed := time.Since(s1)
		fmt.Println("cpu =", elapsed)
	}()

	go func() {
		defer wg.Done()
		s1 := time.Now()
		dashboardVO.MemInfo = systemUtils.GetMemInfo()
		elapsed := time.Since(s1)
		fmt.Println("mem =", elapsed)
	}()

	go func() {
		defer wg.Done()
		s1 := time.Now()
		dashboardVO.DiskInfo = systemUtils.GetDiskInfo()
		elapsed := time.Since(s1)
		fmt.Println("disk =", elapsed)
	}()

	go func() {
		defer wg.Done()
		s1 := time.Now()
		dashboardVO.DiskInfo = systemUtils.GetDiskInfo()
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Printf("程序占用内存：%d Kb\n", m.Alloc/1024)
		dashboardVO.MemStates = m.Alloc / 1024
		elapsed := time.Since(s1)
		fmt.Println("disk =", elapsed)
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				dashboardVO.Version = ""
			}
			wg.Done()
		}()
		dashboardVO.Version = GetLocalDstVersion(clusterName)
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
	elapsed := time.Since(start)
	fmt.Println("Elapsed =", elapsed)

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

func (g *GameService) GetGameConfig(ctx *gin.Context) *world.GameConfig {
	gameConfig := world.GameConfig{}
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

func (g *GameService) SaveGameConfig(ctx *gin.Context, gameConfig *world.GameConfig) {
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
		// SaveAdminlist(world.Adminlist)
		wg.Done()
	}()

	go func() {
		// SaveBlocklist(world.Blocklist)
		wg.Done()
	}()

	go func() {
		g.c.SaveMasterWorld(clusterName, gameConfig.Master)
		dstUtils.DedicatedServerModsSetup(clusterName, gameConfig.Master.Modoverrides)
		wg.Done()
	}()

	go func() {
		g.c.SaveCavesWorld(clusterName, gameConfig.Caves)
		wg.Done()
	}()

	wg.Wait()
}
