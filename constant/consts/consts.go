package consts

import (
	"dst-admin-go/utils/systemUtils"
	"fmt"
	"path/filepath"
)

const (
	StartGame   = 0
	StartMaster = 1
	StartCaves  = 2

	StopGame   = 0
	StopMaster = 1
	StopCaves  = 2

	// ClearScreenCmd 检查目前所有的screen作业，并删除已经无法使用的screen作业
	ClearScreenCmd = "screen -wipe "

	GameRunning       = "gameRunning"
	UpdateGameVersion = "updateGameVersion"
	UpdateGameMod     = "updateGameMod"

	ServerIniTemplate = "./static/template/server.ini"

	MasterLevelType = "MASTER"
	CaveLevelType   = "CAVES"
)

var HomePath string
var KleiDstPath string

func init() {
	home, err := systemUtils.Home()
	if err != nil {
		panic("Home path error: " + err.Error())
	}
	HomePath = home
	fmt.Println("home path: " + HomePath)

	KleiDstPath = filepath.Join(home, ".klei/DoNotStarveTogether")

}
