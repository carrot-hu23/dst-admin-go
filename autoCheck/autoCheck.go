package autoCheck

import (
	"dst-admin-go/constant/consts"
	"log"
)

var AutoCheckObject *AutoCheckConfig

func NewAutoCheckConfig(clusterName string, bin, beta int) *AutoCheckConfig {
	a := AutoCheckConfig{}
	a.InitAutoCheck(clusterName, bin, beta)
	return &a
}

type AutoCheckConfig struct {
	MasterRunning *Monitor
	Slave1Running *Monitor
	Slave2Running *Monitor
	Slave3Running *Monitor
	Slave4Running *Monitor
	Slave5Running *Monitor
	Slave6Running *Monitor
	Slave7Running *Monitor

	MasterModUpdate *Monitor
	Slave1ModUpdate *Monitor
	Slave2ModUpdate *Monitor
	Slave3ModUpdate *Monitor
	Slave4ModUpdate *Monitor
	Slave5ModUpdate *Monitor
	Slave6ModUpdate *Monitor
	Slave7ModUpdate *Monitor

	UpdateGameVersionM *Monitor
}

func (a *AutoCheckConfig) InitAutoCheck(clusterName string, bin, beta int) {
	log.Println("开始初始化循检", clusterName)
	if clusterName != "" {

		a.MasterRunning = NewMonitor(clusterName, "Master", "MasterRun", bin, beta, IsLevelRunning, StartLevelProcess)
		a.Slave1Running = NewMonitor(clusterName, "Slave1", "Slave1Run", bin, beta, IsLevelRunning, StartLevelProcess)
		a.Slave2Running = NewMonitor(clusterName, "Slave2", "Slave2Run", bin, beta, IsLevelRunning, StartLevelProcess)
		a.Slave3Running = NewMonitor(clusterName, "Slave3", "Slave3Run", bin, beta, IsLevelRunning, StartLevelProcess)
		a.Slave4Running = NewMonitor(clusterName, "Slave4", "Slave4Run", bin, beta, IsLevelRunning, StartLevelProcess)
		a.Slave5Running = NewMonitor(clusterName, "Slave5", "Slave5Run", bin, beta, IsLevelRunning, StartLevelProcess)
		a.Slave6Running = NewMonitor(clusterName, "Slave6", "Slave6Run", bin, beta, IsLevelRunning, StartLevelProcess)
		a.Slave7Running = NewMonitor(clusterName, "Slave7", "Slave7Run", bin, beta, IsLevelRunning, StartLevelProcess)

		a.UpdateGameVersionM = NewMonitor(clusterName, "Master", consts.UpdateGameVersion, bin, beta, IsGameUpdateVersionProcess, UpdateGameVersionProcess)

		a.MasterModUpdate = NewMonitor(clusterName, "Master", "MasterMod", bin, beta, IsLevelModUpdateProcess, UpdateModProcess)
		a.Slave1ModUpdate = NewMonitor(clusterName, "Slave1", "Slave1Mod", bin, beta, IsLevelModUpdateProcess, UpdateModProcess)
		a.Slave2ModUpdate = NewMonitor(clusterName, "Slave2", "Slave2Mod", bin, beta, IsLevelModUpdateProcess, UpdateModProcess)
		a.Slave3ModUpdate = NewMonitor(clusterName, "Slave3", "Slave3Mod", bin, beta, IsLevelModUpdateProcess, UpdateModProcess)
		a.Slave4ModUpdate = NewMonitor(clusterName, "Slave4", "Slave4Mod", bin, beta, IsLevelModUpdateProcess, UpdateModProcess)
		a.Slave5ModUpdate = NewMonitor(clusterName, "Slave5", "Slave5Mod", bin, beta, IsLevelModUpdateProcess, UpdateModProcess)
		a.Slave6ModUpdate = NewMonitor(clusterName, "Slave6", "Slave6Mod", bin, beta, IsLevelModUpdateProcess, UpdateModProcess)
		a.Slave7ModUpdate = NewMonitor(clusterName, "Slave7", "Slave7Mod", bin, beta, IsLevelModUpdateProcess, UpdateModProcess)

		go a.MasterRunning.Start()
		go a.Slave1Running.Start()
		go a.Slave2Running.Start()
		go a.Slave3Running.Start()
		go a.Slave4Running.Start()
		go a.Slave5Running.Start()
		go a.Slave6Running.Start()
		go a.Slave7Running.Start()

		go a.UpdateGameVersionM.Start()

		go a.MasterModUpdate.Start()
		go a.Slave1ModUpdate.Start()
		go a.Slave2ModUpdate.Start()
		go a.Slave3ModUpdate.Start()
		go a.Slave4ModUpdate.Start()
		go a.Slave5ModUpdate.Start()
		go a.Slave6ModUpdate.Start()
		go a.Slave7ModUpdate.Start()
	}

}

func (a *AutoCheckConfig) RestartAutoCheck(clusterName string, bin, beta int) {
	log.Println("停止自动巡检")
	a.MasterRunning.Stop()
	a.Slave1Running.Stop()
	a.Slave2Running.Stop()
	a.Slave3Running.Stop()
	a.Slave4Running.Stop()
	a.Slave5Running.Stop()
	a.Slave6Running.Stop()
	a.Slave7Running.Stop()

	a.UpdateGameVersionM.Stop()

	a.MasterModUpdate.Stop()
	a.Slave1ModUpdate.Stop()
	a.Slave2ModUpdate.Stop()
	a.Slave3ModUpdate.Stop()
	a.Slave4ModUpdate.Stop()
	a.Slave5ModUpdate.Stop()
	a.Slave6ModUpdate.Stop()
	a.Slave7ModUpdate.Stop()

	a.InitAutoCheck(clusterName, bin, beta)
	log.Println("重新自动巡检")
}
