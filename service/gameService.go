package service

import (
	"dst-admin-go/constant"
	"dst-admin-go/constant/dst"
	"dst-admin-go/constant/screenKey"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/utils/shellUtils"
	"dst-admin-go/utils/systemUtils"
	"dst-admin-go/vo"
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"
)

type GameService struct {
	lock sync.Mutex
}

func ClearScreen() bool {
	result, err := shellUtils.Shell(constant.CLEAR_SCREEN_CMD)
	if err != nil {
		return false
	}
	res := strings.Split(result, "\n")[0]
	return res != ""
}

func (s *GameService) UpdateGame(clusterName string) {

	s.lock.Lock()
	defer s.lock.Unlock()

	s.stopSpecifiedMaster(clusterName)
	s.stopSpecifiedCaves(clusterName)
	updateGameCMd := dst.GetDstUpdateCmd(clusterName)
	log.Println(updateGameCMd)
	_, err := shellUtils.Shell(updateGameCMd)
	if err != nil {
		log.Panicln("update game error: " + err.Error())
	}
}

func (s *GameService) GetSpecifiedLevelStatus(clusterName, level string) bool {
	cmd := " ps -ef | grep -v grep | grep -v tail |grep '" + clusterName + "'|grep " + level + " |sed -n '1P'|awk '{print $2}' "
	result, err := shellUtils.Shell(cmd)
	if err != nil {
		return false
	}
	res := strings.Split(result, "\n")[0]
	return res != ""
}

func (s *GameService) shutdownSpecifiedLevel(clusterName, level string) {
	if !s.GetSpecifiedLevelStatus(clusterName, level) {
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
func (s *GameService) killSpecifiedLevel(clusterName, level string) {

	if !s.GetSpecifiedLevelStatus(clusterName, level) {
		return
	}
	cmd := " ps -ef | grep -v grep | grep -v tail |grep '" + clusterName + "'|grep " + level + " |sed -n '1P'|awk '{print $2}' |xargs kill -9"
	_, err := shellUtils.Shell(cmd)
	if err != nil {
		// TODO 强制杀掉
		log.Panicln("kill " + clusterName + " " + level + " error: " + err.Error())
	}
}

func (s *GameService) launchSpecifiedLevel(clusterName, level string) {

	cluster := clusterUtils.GetCluster(clusterName)
	dst_install_dir := cluster.ForceInstallDir

	cmd := "cd " + dst_install_dir + "/bin ; screen -d -m -S \"" + screenKey.Key(clusterName, level) + "\"  ./dontstarve_dedicated_server_nullrenderer -console -cluster " + clusterName + " -shard " + level + "  ;"

	_, err := shellUtils.Shell(cmd)
	if err != nil {
		log.Panicln("launch " + clusterName + " " + level + " error: " + err.Error())
	}

}

func (s *GameService) stopSpecifiedMaster(clusterName string) {
	level := "Master"
	s.stopSpecifiedLevel(clusterName, level)
}

func (s *GameService) stopSpecifiedCaves(clusterName string) {

	level := "Caves"
	s.stopSpecifiedLevel(clusterName, level)
}

func (s *GameService) stopSpecifiedLevel(clusterName, level string) {
	s.shutdownSpecifiedLevel(clusterName, level)

	time.Sleep(3 * time.Second)

	if s.GetSpecifiedLevelStatus(clusterName, level) {
		var i uint8 = 1
		for {
			if s.GetSpecifiedLevelStatus(clusterName, level) {
				break
			}
			s.shutdownSpecifiedLevel(clusterName, level)
			time.Sleep(1 * time.Second)
			i++
			if i > 3 {
				break
			}
		}
	}
	s.killSpecifiedLevel(clusterName, level)
}

func (s *GameService) StopSpecifiedGame(clusterName string, opType int) {
	if opType == dst.START_GAME {
		s.stopSpecifiedMaster(clusterName)
		s.stopSpecifiedCaves(clusterName)
	}

	if opType == dst.START_MASTER {
		s.stopSpecifiedMaster(clusterName)
	}

	if opType == dst.START_CAVES {
		s.stopSpecifiedCaves(clusterName)
	}
}

func (s *GameService) launchSpecifiedMaster(clusterName string) {
	level := "Master"
	s.launchSpecifiedLevel(clusterName, level)
}

func (s *GameService) launchSpecifiedCaves(clusterName string) {
	level := "Caves"
	s.launchSpecifiedLevel(clusterName, level)
}

func (s *GameService) StartSpecifiedGame(clusterName string, opType int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if opType == dst.START_GAME {

		s.stopSpecifiedMaster(clusterName)
		s.stopSpecifiedCaves(clusterName)

		s.launchSpecifiedMaster(clusterName)
		s.launchSpecifiedCaves(clusterName)
	}

	if opType == dst.START_MASTER {
		s.stopSpecifiedMaster(clusterName)
		s.launchSpecifiedMaster(clusterName)
	}

	if opType == dst.START_CAVES {
		s.stopSpecifiedCaves(clusterName)
		s.launchSpecifiedCaves(clusterName)
	}

	ClearScreen()
}

func (s *GameService) GetSpecifiedClusterDashboard(clusterName string) vo.DashboardVO {
	var wg sync.WaitGroup
	wg.Add(10)

	dashboardVO := vo.NewDashboardVO(clusterName)
	start := time.Now()
	go func() {
		defer wg.Done()
		s1 := time.Now()
		dashboardVO.MasterStatus = s.GetSpecifiedLevelStatus(clusterName, "Master")
		elapsed := time.Since(s1)
		fmt.Println("master =", elapsed)
	}()

	go func() {
		defer wg.Done()
		s1 := time.Now()
		dashboardVO.CavesStatus = s.GetSpecifiedLevelStatus(clusterName, "Caves")
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
		defer wg.Done()
		dashboardVO.Version = ""
	}()

	// 获取master进程占用情况
	go func() {
		defer wg.Done()
		dashboardVO.MasterPs = s.PsAuxSpecified(clusterName, "Master")
	}()
	// 获取caves进程占用情况
	go func() {
		defer wg.Done()
		dashboardVO.CavesPs = s.PsAuxSpecified(clusterName, "Caves")
	}()

	wg.Wait()
	elapsed := time.Since(start)
	fmt.Println("Elapsed =", elapsed)

	return *dashboardVO
}

func (s *GameService) PsAuxSpecified(clusterName, level string) *vo.DstPsVo {
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
