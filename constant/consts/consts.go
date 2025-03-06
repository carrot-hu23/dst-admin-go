package consts

import (
	"dst-admin-go/utils/systemUtils"
	"fmt"
	"runtime"
)

const (

	// ClearScreenCmd 检查目前所有的screen作业，并删除已经无法使用的screen作业
	ClearScreenCmd = "screen -wipe "

	MasterRunning     = "masterRunning"
	CavesRunning      = "cavesRunning"
	UpdateGameVersion = "updateGameVersion"
	UpdateGameMod     = "updateGameMod"
	UpdateMasterMod   = "updateMasterMod"
	UpdateCavesMod    = "updateCavesMod"

	UPDATE_GAME = "UPDATE_GAME"
	LEVEL_MOD   = "LEVEL_MOD"
	LEVEL_DOWN  = "LEVEL_DOWN"
)

var HomePath string

const PasswordPath = "./password.txt"

func init() {
	home, err := systemUtils.Home()
	if err != nil {
		panic("Home path error: " + err.Error())
	}
	HomePath = home
	fmt.Println("home path: " + HomePath)

	if runtime.GOOS == "windows" {
		//klei_path = filepath.Join(consts.HomePath, "Documents", "klei", "DoNotStarveTogether")
		// DefaultKleiDstPath = filepath.Join(home, "Documents", "klei", "DoNotStarveTogether")
	} else {
		// klei_path = filepath.Join(consts.HomePath, ".klei", "DoNotStarveTogether")
		// DefaultKleiDstPath = filepath.Join(home, ".klei", "DoNotStarveTogether")
	}

}
