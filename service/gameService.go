package service

import (
	"dst-admin-go/constant/consts"
	"dst-admin-go/model"
	"dst-admin-go/utils/dstUtils"
	"dst-admin-go/utils/levelConfigUtils"
	"dst-admin-go/utils/systemUtils"
	"io"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"

	"dst-admin-go/constant/screenKey"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/utils/shellUtils"
	"dst-admin-go/vo"
	"log"
	"strings"
	"sync"
	"time"
)

var launchLock = sync.Mutex{}

type GameService struct {
	lock sync.Mutex
	c    HomeService

	logRecord LogRecordService
}

func (g *GameService) GetLastDstVersion() int64 {
	if isWindows() {
		return WindowService.GetLastDstVersion()
	}

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
	if isWindows() {
		return WindowService.GetLocalDstVersion(clusterName)
	}

	cluster := clusterUtils.GetCluster(clusterName)
	versionTextPath := filepath.Join(cluster.ForceInstallDir, "version.txt")

	// 使用filepath.Clean确保路径格式正确
	cleanPath := filepath.Clean(versionTextPath)

	// 使用filepath.Abs获取绝对路径
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		log.Println("Error getting absolute path:", err)
		return 0
	}

	// 打印绝对路径
	log.Println("Absolute Path:", absPath)

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
	if isWindows() {
		return true
	}

	result, err := shellUtils.Shell(consts.ClearScreenCmd)
	if err != nil {
		return false
	}
	res := strings.Split(result, "\n")[0]
	return res != ""
}

func (g *GameService) UpdateGame(clusterName string) error {
	if isWindows() {
		return WindowService.UpdateGame(clusterName)
	}

	g.lock.Lock()
	defer g.lock.Unlock()
	// TODO 关闭相应的世界
	g.StopGame(clusterName)

	updateGameCMd := dstUtils.GetDstUpdateCmd(clusterName)
	log.Println("正在更新游戏", "cluster: ", clusterName, "command: ", updateGameCMd)
	_, err := shellUtils.Shell(updateGameCMd)

	// TODO 写入 DedicatedServerModsSetup.lua
	levelConfig, _ := levelConfigUtils.GetLevelConfig(clusterName)
	for i := range levelConfig.LevelList {
		level := homeServe.GetLevel(clusterName, levelConfig.LevelList[i].File)
		modoverrides := level.Modoverrides
		dstUtils.DedicatedServerModsSetup2(clusterName, modoverrides)
	}

	if err != nil {
		return err
	}
	return nil
}

func (g *GameService) GetLevelStatus(clusterName, level string) bool {

	if isWindows() {
		return WindowService.GetLevelStatus(clusterName, level)
	}

	cmd := " ps -ef | grep -v grep | grep -v tail |grep '" + clusterName + "'|grep " + level + " |sed -n '1P'|awk '{print $2}' "
	result, err := shellUtils.Shell(cmd)
	if err != nil {
		return false
	}
	res := strings.Split(result, "\n")[0]
	// log.Println("查询世界状态", cmd, result, res, res != "")
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

func (g *GameService) StartLevel(clusterName, level string, bin, beta int) {
	if isWindows() {
		WindowService.StartLevel(clusterName, level, bin, beta)
		return
	}
	g.StopLevel(clusterName, level)
	g.LaunchLevel(clusterName, level, bin, beta)
	ClearScreen()
}

func (g *GameService) LaunchLevel(clusterName, level string, bin, beta int) {
	if isWindows() {
		WindowService.LaunchLevel(clusterName, level, bin, beta)
		return
	}
	launchLock.Lock()
	defer func() {
		launchLock.Unlock()
		if r := recover(); r != nil {
		}
	}()

	g.logRecord.RecordLog(clusterName, level, model.RUN)

	cluster := clusterUtils.GetCluster(clusterName)
	dstInstallDir := cluster.ForceInstallDir
	ugcDirectory := cluster.Ugc_directory
	persistent_storage_root := cluster.Persistent_storage_root
	conf_dir := cluster.Conf_dir
	var startCmd = ""

	dstInstallDir = dstUtils.EscapePath(dstInstallDir)
	log.Println(dstInstallDir)

	if bin == 64 {
		startCmd = "cd " + dstInstallDir + "/bin64 ; screen -d -m -S \"" + screenKey.Key(clusterName, level) + "\"  ./dontstarve_dedicated_server_nullrenderer_x64 -console -cluster " + clusterName + " -shard " + level
	} else if bin == 100 {
		startCmd = "cd " + dstInstallDir + "/bin64 ; screen -d -m -S \"" + screenKey.Key(clusterName, level) + "\"  ./dontstarve_dedicated_server_nullrenderer_x64_luajit -console -cluster " + clusterName + " -shard " + level
	} else {
		startCmd = "cd " + dstInstallDir + "/bin ; screen -d -m -S \"" + screenKey.Key(clusterName, level) + "\"  ./dontstarve_dedicated_server_nullrenderer -console -cluster " + clusterName + " -shard " + level
	}

	if ugcDirectory != "" {
		startCmd += " -ugc_directory " + ugcDirectory
	}
	if persistent_storage_root != "" {
		startCmd += " -persistent_storage_root " + persistent_storage_root
	}
	if conf_dir != "" {
		startCmd += " -conf_dir " + conf_dir
	}

	startCmd += "  ;"

	log.Println("正在启动世界", "cluster: ", clusterName, "level: ", level, "command: ", startCmd)
	_, err := shellUtils.Shell(startCmd)
	if err != nil {
		log.Panicln("启动 "+clusterName+" "+level+" error,", err)
	}

}

func (g *GameService) StopLevel(clusterName, level string) {
	if isWindows() {
		WindowService.StopLevel(clusterName, level)
		return
	}
	launchLock.Lock()
	defer func() {
		launchLock.Unlock()
		if r := recover(); r != nil {
		}
	}()

	g.logRecord.RecordLog(clusterName, level, model.STOP)

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

func (g *GameService) StopGame(clusterName string) {

	if isWindows() {
		WindowService.StopGame(clusterName)
		return
	}
	config, err := levelConfigUtils.GetLevelConfig(clusterName)
	if err != nil {
		log.Panicln(err)
	}
	var wg sync.WaitGroup
	wg.Add(len(config.LevelList))
	for i := range config.LevelList {
		go func(i int) {
			defer func() {
				wg.Done()
				if r := recover(); r != nil {
					log.Println(r)
				}
			}()
			levelName := config.LevelList[i].File
			g.StopLevel(clusterName, levelName)
		}(i)
	}
	wg.Wait()
}

func (g *GameService) StartGame(clusterName string) {
	if isWindows() {
		WindowService.StartGame(clusterName)
		return
	}
	g.lock.Lock()
	defer g.lock.Unlock()
	g.StopGame(clusterName)
	cluster := clusterUtils.GetCluster(clusterName)
	bin := cluster.Bin
	beta := cluster.Beta

	config, err := levelConfigUtils.GetLevelConfig(clusterName)
	if err != nil {
		log.Panicln(err)
	}
	var wg sync.WaitGroup
	wg.Add(len(config.LevelList))
	for i := range config.LevelList {
		go func(i int) {
			defer func() {
				wg.Done()
				if r := recover(); r != nil {
					log.Println(r)
				}
			}()
			levelName := config.LevelList[i].File
			g.LaunchLevel(clusterName, levelName, bin, beta)
		}(i)
	}
	ClearScreen()
	wg.Wait()
}

func (g *GameService) PsAuxSpecified(clusterName, level string) *vo.DstPsVo {
	if isWindows() {
		return vo.NewDstPsVo()
	}
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

type SystemInfo struct {
	HostInfo      *systemUtils.HostInfo `json:"host"`
	CpuInfo       *systemUtils.CpuInfo  `json:"cpu"`
	MemInfo       *systemUtils.MemInfo  `json:"mem"`
	DiskInfo      *systemUtils.DiskInfo `json:"disk"`
	PanelMemUsage uint64                `json:"panelMemUsage"`
	PanelCpuUsage float64               `json:"panelCpuUsage"`
}

func (g *GameService) GetSystemInfo(clusterName string) *SystemInfo {
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

		// 获取当前程序使用的CPU信息
		//startCPU, _ := cpu.Percent(time.Second, false)
		//time.Sleep(1 * time.Second) // 假设程序运行1秒
		//endCPU, _ := cpu.Percent(time.Second, false)
		//cpuUsage := endCPU[0] - startCPU[0]
		//dashboardVO.PanelCpuUsage = cpuUsage

	}()

	wg.Wait()
	return &dashboardVO
}
