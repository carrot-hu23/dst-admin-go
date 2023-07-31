package autoCheck

import (
	"dst-admin-go/config/global"
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
	MasterModUpdate    *Monitor
	CavesModUpdate     *Monitor
	UpdateGameVersionM *Monitor
}

func (a *AutoCheckConfig) InitAutoCheck(clusterName string, bin, beta int) {
	log.Println("开始初始化循检", clusterName)
	config := global.Config
	if clusterName != "" {
		a.MasterRunning = NewMonitor(clusterName, consts.MasterRunning, bin, beta, IsMasterRunning, time.Duration(config.AutoCheck.MasterInterval)*time.Minute, StartMasterProcess)
		a.CavesRunning = NewMonitor(clusterName, consts.CavesRunning, bin, beta, IsCavesRunning, time.Duration(config.AutoCheck.CavesInterval)*time.Minute, StartCavesProcess)
		a.UpdateGameVersionM = NewMonitor(clusterName, consts.UpdateGameVersion, bin, beta, IsGameUpdateVersionProcess, time.Duration(config.AutoCheck.GameUpdateInterval)*time.Minute, UpdateGameVersionProcess)
		a.MasterModUpdate = NewMonitor(clusterName, consts.UpdateMasterMod, bin, beta, IsMasterModUpdateProcess, time.Duration(config.AutoCheck.MasterModInterval)*time.Minute, UpdateMasterModUpdateProcess)
		a.CavesModUpdate = NewMonitor(clusterName, consts.UpdateCavesMod, bin, beta, IsCavesModUpdateProcess, time.Duration(config.AutoCheck.CavesModInterval)*time.Minute, UpdateCavesModUpdateProcess)

		go a.MasterRunning.Start()
		go a.CavesRunning.Start()
		go a.UpdateGameVersionM.Start()
		go a.MasterModUpdate.Start()
		go a.CavesModUpdate.Start()
	}

}

func (a *AutoCheckConfig) RestartAutoCheck(clusterName string, bin, beta int) {
	log.Println("停止自动巡检")
	a.MasterRunning.Stop()
	a.CavesRunning.Stop()
	a.UpdateGameVersionM.Stop()
	a.MasterModUpdate.Stop()
	a.CavesModUpdate.Stop()
	a.InitAutoCheck(clusterName, bin, beta)
	log.Println("重新自动巡检")
}
