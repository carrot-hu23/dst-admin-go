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
	MasterRunning      *Monitor
	CavesRunning       *Monitor
	UpdateGameVersionM *Monitor
	UpdateGameModM     *Monitor
}

func (a *AutoCheckConfig) InitAutoCheck(clusterName string, bin, beta int) {
	log.Println("开始初始化循检", clusterName)
	if clusterName != "" {
		a.MasterRunning = NewMonitor(clusterName, consts.MasterRunning, bin, beta, IsMasterRunning, 5*time.Minute, StartMasterProcess)
		a.CavesRunning = NewMonitor(clusterName, consts.CavesRunning, bin, beta, IsCavesRunning, 5*time.Minute, StartCavesProcess)
		a.UpdateGameVersionM = NewMonitor(clusterName, consts.UpdateGameVersion, bin, beta, IsGameUpdateVersionProcess, 10*time.Minute, UpdateGameVersionProcess)
		a.UpdateGameModM = NewMonitor(clusterName, consts.UpdateGameMod, bin, beta, IsGameModUpdateProcess, 10*time.Minute, UpdateGameVersionProcess)

		go a.MasterRunning.Start()
		go a.CavesRunning.Start()
		go a.UpdateGameVersionM.Start()
		go a.UpdateGameModM.Start()
	}

}

func (a *AutoCheckConfig) RestartAutoCheck(clusterName string, bin, beta int) {
	log.Println("停止自动巡检")
	a.MasterRunning.Stop()
	a.CavesRunning.Stop()
	a.UpdateGameVersionM.Stop()
	a.UpdateGameModM.Stop()
	a.InitAutoCheck(clusterName, bin, beta)
	log.Println("重新自动巡检")
}
