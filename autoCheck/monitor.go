package autoCheck

import (
	"dst-admin-go/config/database"
	"dst-admin-go/config/global"
	"dst-admin-go/constant/consts"
	"dst-admin-go/mod"
	"dst-admin-go/model"
	"dst-admin-go/service"
	"dst-admin-go/utils/dstConfigUtils"
	"dst-admin-go/utils/dstUtils"
	"dst-admin-go/utils/fileUtils"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var gameService = service.GameService{}
var gameConsoleService = service.GameConsoleService{}

type Monitor struct {
	t             string
	beta          int
	bin           int
	clusterName   string
	checkFunc     func(clusterName string, int2 int, int3 int) bool
	checkInterval time.Duration
	startFunc     func(clusterName string, int2 int, int3 int) error
	stopCh        chan int
}

func NewMonitor(clusterName string, t string, bin int, beta int, checkFunc func(string2 string, int2 int, int3 int) bool, checkInterval time.Duration, startFunc func(string2 string, int2 int, int3 int) error) *Monitor {
	return &Monitor{
		t:             t,
		beta:          beta,
		bin:           bin,
		clusterName:   clusterName,
		checkFunc:     checkFunc,
		checkInterval: checkInterval,
		startFunc:     startFunc,
		stopCh:        make(chan int, 1),
	}
}

func (m *Monitor) Start() {
	log.Println("Starting monitor cluster: ", m.clusterName, m.t, "bin:", m.bin, "beta:", m.beta)
	for {
		select {
		case <-m.stopCh:
			log.Println("Stopping monitor cluster: ", m.clusterName, m.t)
			return
		default:
			if !m.checkFunc(m.clusterName, m.bin, m.beta) {
				log.Println(m.clusterName, m.t, " is not running, waiting for ", m.checkInterval)
				log.Println("start sleep")
				time.Sleep(m.checkInterval)
				log.Println("end sleep")
				if !m.checkFunc(m.clusterName, m.bin, m.beta) {
					log.Println(m.clusterName, m.t, "has not started, starting it...")
					err := m.startFunc(m.clusterName, m.bin, m.beta)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
			time.Sleep(m.checkInterval)
		}
	}
}

func (m *Monitor) Stop() {
	m.stopCh <- 1
}

func IsMasterRunning(clusterName string, bin, beta int) bool {
	db := database.DB
	autoCheck := model.AutoCheck{}
	db.Where("name = ?", consts.MasterRunning).Find(&autoCheck)
	if autoCheck.Enable == 1 {
		return gameService.GetLevelStatus(clusterName, "Master")
	}
	return true
}

func IsCavesRunning(clusterName string, bin, beta int) bool {
	db := database.DB
	autoCheck := model.AutoCheck{}
	db.Where("name = ?", consts.CavesRunning).Find(&autoCheck)
	if autoCheck.Enable == 1 {
		return gameService.GetLevelStatus(clusterName, "Caves")
	}
	return true
}

func StartMasterProcess(clusterName string, bin, beta int) error {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	gameService.StartGame(clusterName, bin, beta, consts.StartMaster)
	return nil
}

func StartCavesProcess(clusterName string, bin, beta int) error {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	gameService.StartGame(clusterName, bin, beta, consts.StartCaves)
	return nil
}

func IsGameUpdateVersionProcess(clusterName string, bin, beta int) bool {
	db := database.DB
	autoCheck := model.AutoCheck{}
	db.Where("name = ?", consts.UpdateGameVersion).Find(&autoCheck)
	if autoCheck.Enable == 1 {
		// diff dst version
		localVersion := gameService.GetLocalDstVersion(clusterName)
		version := gameService.GetLastDstVersion()
		log.Println("localVersion: ", localVersion, "lastVersion: ", version)
		return localVersion == version
	}
	return true
}

func UpdateGameVersionProcess(clusterName string, bin, beta int) error {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	for i := 0; i < 3; i++ {
		gameConsoleService.SentBroadcast(clusterName, global.Config.AutoCheck.GameUpdatePrompt)
		time.Sleep(3 * time.Second)
	}
	err := gameService.UpdateGame(clusterName)
	if err != nil {
		return err
	}
	gameService.StartGame(clusterName, consts.StartGame, bin, beta)
	return nil
}

func IsMasterModUpdateProcess(clusterName string, bin, beta int) bool {
	return isLevelModUpdateProcess(clusterName, bin, beta, consts.Master)
}

func UpdateMasterModUpdateProcess(clusterName string, bin, beta int) error {
	for i := 0; i < 3; i++ {
		gameConsoleService.SentBroadcast(clusterName, global.Config.AutoCheck.ModUpdatePrompt)
		time.Sleep(3 * time.Second)
	}
	return updateLevelModUpdateProcess(clusterName, bin, beta, consts.StartMaster)
}

func IsCavesModUpdateProcess(clusterName string, bin, beta int) bool {
	return isLevelModUpdateProcess(clusterName, bin, beta, consts.Caves)
}

func UpdateCavesModUpdateProcess(clusterName string, bin, beta int) error {
	for i := 0; i < 3; i++ {
		gameConsoleService.SentBroadcast(clusterName, global.Config.AutoCheck.ModUpdatePrompt)
		time.Sleep(3 * time.Second)
	}
	return updateLevelModUpdateProcess(clusterName, bin, beta, consts.StartMaster)
}

func isLevelModUpdateProcess(clusterName string, bin, beta int, levelName string) bool {

	// 找到当前存档的modId, 然后根据判断当前存档的
	dstConfig := dstConfigUtils.GetDstConfig()
	cluster := dstConfig.Cluster
	masterAcfPath := filepath.Join(dstConfig.Force_install_dir, "ugc_mods", cluster, levelName)
	acfWorkShops := dstUtils.ParseACFFile(masterAcfPath)

	var needUpdate atomic.Bool
	var wg sync.WaitGroup
	for key := range acfWorkShops {
		wg.Add(1)
		go func(key string) {
			defer func() {
				if r := recover(); r != nil {
					log.Println(r)
				}
				wg.Done()
			}()
			acfWorkShop := acfWorkShops[key]
			modInfo, err, _ := mod.GetModInfo(key)
			log.Println(key, acfWorkShop.TimeUpdated, modInfo.LastTime)
			if err == nil {
				if float64(acfWorkShop.TimeUpdated) != modInfo.LastTime {
					needUpdate.Store(true)
				}
			}
		}(key)
	}
	wg.Wait()
	return needUpdate.Load()
}

func updateLevelModUpdateProcess(clusterName string, bin, beta int, startOpt int) error {
	log.Println("开始更新mod", clusterName)
	dstPath := dstConfigUtils.GetDstConfig().Force_install_dir
	modsPath := filepath.Join(dstPath, "mods")
	directories, err := fileUtils.ListDirectories(modsPath)
	if err != nil {
		log.Println("delete dst workshop file error", err)
	}
	var workshopList []string
	for _, directory := range directories {
		if strings.Contains(directory, "workshop") {
			workshopList = append(workshopList, directory)
		}
	}
	for _, workshop := range workshopList {
		err := fileUtils.DeleteDir(workshop)
		if err != nil {
			log.Println("删除mod失败", err)
		}
	}
	gameService.StartGame(clusterName, bin, beta, startOpt)
	return nil
}
