package service

import (
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/utils/systemUtils"
	"dst-admin-go/vo"
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	wg sync.WaitGroup
)

const N = 11

func GetBashboard() vo.DashboardVO {
	wg.Add(N)
	dashboardVO := vo.NewDashboardVO()
	start := time.Now()
	go func() {
		defer wg.Done()
		s1 := time.Now()
		dashboardVO.MasterStatus = getMasterStatus()
		elapsed := time.Since(s1)
		fmt.Println("master =", elapsed)
	}()

	go func() {
		defer wg.Done()
		s1 := time.Now()
		dashboardVO.CavesStatus = getCavesStatus()
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

	// 获取master进程占用情况
	go func() {
		defer wg.Done()
		dashboardVO.MasterPs = DstPs("Master")
	}()
	// 获取caves进程占用情况
	go func() {
		defer wg.Done()
		dashboardVO.CavesPs = DstPs("Caves")
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

	wg.Wait()
	elapsed := time.Since(start)
	fmt.Println("Elapsed =", elapsed)

	return *dashboardVO
}

func DstPs(psName string) *vo.DstPsVo {
	dstPsVo := vo.NewDstPsVo()
	info := PsAux(psName)

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

const version_path = "./version"

func GetDstVersion() string {

	exists := fileUtils.Exists(version_path)
	if !exists {
		if err := fileUtils.CreateFile(version_path); err != nil {
			log.Panicln("create version file error: " + err.Error())
		}
	}
	str, err := fileUtils.ReadFile(version_path)
	if err != nil {
		log.Panicln("read version file error: " + err.Error())
	}
	return str
}

func SetDstVersion(version string) {
	fileUtils.WriterTXT(version_path, version)
}
