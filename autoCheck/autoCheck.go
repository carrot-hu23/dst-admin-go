package autoCheck

import (
	"dst-admin-go/constant/consts"
	"log"
	"time"
)

var AutoCheckObject *AutoCheckConfig

func NewAutoCheckConfig(clusterName string, bin, beta int) *AutoCheckConfig {
	a := AutoCheckConfig{}
	a.InitAutoCheck(clusterName, bin, beta)
	return &a
}

type AutoCheckConfig struct {
	GameRunningM       *Monitor
	UpdateGameVersionM *Monitor
	UpdateGameModM     *Monitor
}

func (a *AutoCheckConfig) InitAutoCheck(clusterName string, bin, beta int) {
	log.Println("开始初始化循检", clusterName)
	if clusterName != "" {
		a.GameRunningM = NewMonitor(clusterName, consts.GameRunning, bin, beta, IsGameRunning, 5*time.Minute, StartGameProcess)
		a.UpdateGameVersionM = NewMonitor(clusterName, consts.UpdateGameVersion, bin, beta, IsGameUpdateVersionProcess, 10*time.Minute, UpdateGameVersionProcess)
		a.UpdateGameModM = NewMonitor(clusterName, consts.UpdateGameMod, bin, beta, IsGameModUpdateProcess, 10*time.Minute, UpdateGameVersionProcess)

		go a.GameRunningM.Start()
		go a.UpdateGameVersionM.Start()
		go a.UpdateGameModM.Start()
	}

}

func (a *AutoCheckConfig) RestartAutoCheck(clusterName string, bin, beta int) {
	log.Println("停止自动巡检")
	a.GameRunningM.Stop()
	a.UpdateGameVersionM.Stop()
	a.UpdateGameModM.Stop()
	a.InitAutoCheck(clusterName, bin, beta)
	log.Println("重新自动巡检")
}
