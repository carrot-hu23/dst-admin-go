package service

import (
	"dst-admin-go/constant/dst"
	"dst-admin-go/utils/dstConfigUtils"
	"dst-admin-go/utils/systemUtils"
	"dst-admin-go/vo"
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"
)

func getscreenKey(clusterName, level string) string {
	return "DST_" + level + "_" + clusterName
}

/*
ps -ef | grep -v grep | grep MyDediServer | grep Master | sed -n '1P' | awk '{print $2}'
*/
func GetSpecifiedLevelStatus(clusterName, level string) bool {
	cmd := " ps -ef | grep -v grep | grep -v tail |grep '" + clusterName + "'|grep " + level + " |sed -n '1P'|awk '{print $2}' "
	result, err := Shell(cmd)
	if err != nil {
		return false
	}
	res := strings.Split(result, "\n")[0]
	return res != ""
}

func shutdownSpecifiedLevel(clusterName, level string) {
	if !GetSpecifiedLevelStatus(clusterName, level) {
		return
	}
	screenKey := getscreenKey(clusterName, level)
	shell := "screen -S \"" + screenKey + "\" -p 0 -X stuff \"c_shutdown(true)\\n\""
	_, err := Shell(shell)
	if err != nil {
		log.Panicln("shut down " + clusterName + " " + level + " error: " + err.Error())
	}
}

/*
STOP_CAVES_CMD = "ps -ef | grep -v grep |grep '" + DST_CAVES + "' |sed -n '1P'|awk '{print $2}' |xargs kill -9"
*/
func killSpecifiedLevel(clusterName, level string) {

	if !GetSpecifiedLevelStatus(clusterName, level) {
		return
	}
	cmd := " ps -ef | grep -v grep | grep -v tail |grep '" + clusterName + "'|grep " + level + " |sed -n '1P'|awk '{print $2}' |xargs kill -9"
	_, err := Shell(cmd)
	if err != nil {
		log.Panicln("kill " + clusterName + " " + level + " error: " + err.Error())
	}
}

func launchSpecifiedLevel(clusterName, level string) {

	dstConfig := dstConfigUtils.GetDstConfig()
	cluster := dstConfig.Cluster
	dst_install_dir := dstConfig.Force_install_dir

	screenKey := getscreenKey(clusterName, level)

	cmd := "cd " + dst_install_dir + "/bin ; screen -d -m -S \"" + screenKey + "\"  ./dontstarve_dedicated_server_nullrenderer -console -cluster " + cluster + " -shard " + level + "  ;"

	_, err := Shell(cmd)
	if err != nil {
		log.Panicln("launch " + cluster + " " + level + " error: " + err.Error())
	}

}

func stopSpecifiedMaster(clusterName string) {
	level := "Master"
	stopSpecifiedLevel(clusterName, level)
}

func stopSpecifiedCaves(clusterName string) {

	level := "Caves"
	stopSpecifiedLevel(clusterName, level)
}

func stopSpecifiedLevel(clusterName, level string) {
	shutdownSpecifiedLevel(clusterName, level)

	time.Sleep(3 * time.Second)

	if GetSpecifiedLevelStatus(clusterName, level) {
		var i uint8 = 1
		for {
			if GetSpecifiedLevelStatus(clusterName, level) {
				break
			}
			time.Sleep(1 * time.Second)
			i++
			if i > 3 {
				break
			}
		}
	}
	killSpecifiedLevel(clusterName, level)
}

func StopSpecifiedGame(clusterName string, opType int) {
	if opType == dst.START_GAME {
		stopSpecifiedMaster(clusterName)
		stopSpecifiedCaves(clusterName)
	}

	if opType == dst.START_MASTER {
		stopSpecifiedMaster(clusterName)
	}

	if opType == dst.START_CAVES {
		stopSpecifiedCaves(clusterName)
	}
}

func launchSpecifiedMaster(clusterName string) {
	level := "Master"
	launchSpecifiedLevel(clusterName, level)
}

func launchSpecifiedCaves(clusterName string) {
	level := "Caves"
	launchSpecifiedLevel(clusterName, level)
}

func StartSpecifiedGame(clusterName string, opType int) {
	if opType == dst.START_GAME {

		stopSpecifiedMaster(clusterName)
		stopSpecifiedCaves(clusterName)

		launchSpecifiedMaster(clusterName)
		launchSpecifiedCaves(clusterName)
	}

	if opType == dst.START_MASTER {
		stopSpecifiedMaster(clusterName)
		launchSpecifiedMaster(clusterName)
	}

	if opType == dst.START_CAVES {
		stopSpecifiedCaves(clusterName)
		launchSpecifiedCaves(clusterName)
	}

	ClearScreen()
}

func GetSpecifiedClusterDashboard(clusterName string) vo.DashboardVO {
	wg.Add(N)
	dashboardVO := vo.NewDashboardVO()
	start := time.Now()
	go func() {
		defer wg.Done()
		s1 := time.Now()
		dashboardVO.MasterStatus = GetSpecifiedLevelStatus(clusterName, "Master")
		elapsed := time.Since(s1)
		fmt.Println("master =", elapsed)
	}()

	go func() {
		defer wg.Done()
		s1 := time.Now()
		dashboardVO.CavesStatus = GetSpecifiedLevelStatus(clusterName, "Caves")
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
		dashboardVO.Version = GetDstVersion()
	}()

	// 获取直连ip
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Error reading caves file: %v\n", r)
				dashboardVO.IpConnect = ""
			}
			wg.Done()
		}()
		ip, _ := GetPublicIP()
		dashboardVO.IpConnect = ip
	}()

	// 获取master进程占用情况
	go func() {
		defer wg.Done()
		dashboardVO.MasterPs = PsAuxSpecified(clusterName, "Master")
	}()
	// 获取caves进程占用情况
	go func() {
		defer wg.Done()
		dashboardVO.CavesPs = PsAuxSpecified(clusterName, "Caves")
	}()

	wg.Wait()
	elapsed := time.Since(start)
	fmt.Println("Elapsed =", elapsed)

	return *dashboardVO
}

/*
 */
func PsAuxSpecified(clusterName, level string) *vo.DstPsVo {
	dstPsVo := vo.NewDstPsVo()
	cmd := "ps -aux | grep -v grep | grep -v tail | grep " + clusterName + "  | grep " + level + " | sed -n '2P' |awk '{print $3, $4, $5, $6}'"

	info, err := Shell(cmd)
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
