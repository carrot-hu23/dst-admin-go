package autoCheck

import (
	"dst-admin-go/config/database"
	"dst-admin-go/constant/consts"
	"dst-admin-go/model"
	"dst-admin-go/service"
	"log"
	"time"
)

var gameService = service.GameService{}

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
	log.Println("bin: ", bin, "beta: ", beta)
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

func IsGameRunning(clusterName string, bin, beta int) bool {
	db := database.DB
	autoCheck := model.AutoCheck{}
	db.Where("name = ?", consts.GameRunning).Find(&autoCheck)
	if autoCheck.Enable == 1 {
		return gameService.GetLevelStatus(clusterName, "Master")
	}
	return true
}

func StartGameProcess(clusterName string, bin, beta int) error {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	gameService.StartGame(clusterName, bin, beta, consts.StartGame)
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
	err := gameService.UpdateGame(clusterName)
	if err != nil {
		return err
	}
	gameService.StartGame(clusterName, consts.StartGame, bin, beta)
	return nil
}

// TODO
func IsGameModUpdateProcess(clusterName string, bin, beta int) bool {

	return true
}

func UpdateGameModUpdateProcess(clusterName string, bin, beta int) error {
	return nil
}
