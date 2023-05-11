package service

import (
	"bytes"
	"dst-admin-go/constant"
	optype "dst-admin-go/constant/opType"
	"dst-admin-go/vo"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
)

func UpdateGame() {
	SentBroadcast(":pig 正在更新游戏......")
	ElegantShutdownMaster()
	ElegantShutdownCaves()
	log.Println(constant.GET_UPDATE_GAME_CMD())
	_, err := Shell(constant.GET_UPDATE_GAME_CMD())
	if err != nil {
		log.Panicln("update game error: " + err.Error())
	}
}

func StartGame(opType int) {

	if opType == optype.START_GAME {
		SentBroadcast(":pig 正在重启游戏......")
		stopMaster()
		stopCaves()
		startMaster()
		startCaves()
	}

	if opType == optype.START_MASTER {
		SentBroadcast(":pig 正在重启世界......")
		stopMaster()
		startMaster()
	}

	if opType == optype.START_CAVES {
		SentBroadcast(":pig 正在重启洞穴......")
		stopCaves()
		startCaves()
	}
	ClearScreen()
}

func StopGame(opType int) {

	if opType == optype.START_GAME {
		SentBroadcast(":pig 正在停止游戏......")
		ElegantShutdownMaster()
		ElegantShutdownCaves()
	}

	if opType == optype.START_MASTER {
		SentBroadcast(":pig 正在停止世界......")
		ElegantShutdownMaster()
	}

	if opType == optype.START_CAVES {
		SentBroadcast(":pig 正在停止洞穴......")
		ElegantShutdownCaves()
	}

}

func StartMaster() {
	SentBroadcast(":pig 正在重启世界......")
	stopMaster()
	startMaster()
}

func StartCaves() {
	SentBroadcast(":pig 正在重启洞穴......")
	stopCaves()
	startCaves()
}

const max_waitting = 10

func ElegantShutdownMaster() {

	if getMasterStatus() {
		var i uint8 = 1
		for {
			if getMasterStatus() {
				break
			}
			time.Sleep(1 * time.Second)
			i++
			if i >= max_waitting {
				break
			}
		}
	}
	SentBroadcast(":pig 正在关闭世界......")
	SentBroadcast(":pig 正在关闭世界......")
	SentBroadcast(":pig 正在关闭世界......")
	ShutdownMaster()
	stopMaster()
}

func ElegantShutdownCaves() {

	if getCavesStatus() {
		var i uint8 = 1
		for {
			if getCavesStatus() {
				break
			}
			time.Sleep(1 * time.Second)
			i++
			if i >= max_waitting {
				break
			}
		}
	}
	SentBroadcast(":pig 正在关闭洞穴......")
	SentBroadcast(":pig 正在关闭洞穴......")
	SentBroadcast(":pig 正在关闭洞穴......")
	shutdownCaves()
	stopCaves()
}

func ShutdownMaster() {
	shell := "screen -S \"" + constant.SCREEN_WORK_MASTER_NAME + "\" -p 0 -X stuff \"c_shutdown(true)\\n\""
	_, err := Shell(shell)
	if err != nil {
		log.Panicln("shut down master error: " + err.Error())
	}
}

func shutdownCaves() {
	shell := "screen -S \"" + constant.SCREEN_WORK_CAVES_NAME + "\" -p 0 -X stuff \"c_shutdown(true)\\n\""
	_, err := Shell(shell)
	if err != nil {
		log.Panicln("shut down caves error: " + err.Error())
	}
}

func stopMaster() {
	//TODO 写入mod安装文件
	check_cmd := "ps -ef | grep -v grep |grep '" + constant.DST_MASTER + "' |sed -n '1P'|awk '{print $2}'"
	check, error := Shell(check_cmd)

	if error != nil || check == "" {
		log.Printf("grep Master pId error")
	} else {
		_, err := Shell(constant.STOP_MASTER_CMD)
		if err != nil {
			log.Panicln("shut down caves error: " + err.Error())
		}
	}
}

func startMaster() {
	// TODO 有些变量还没有修改完成
	_, err := Shell(constant.GET_START_MASTER_CMD())
	if err != nil {
		log.Panicln("start master error: " + err.Error())
	}
}

func stopCaves() {
	check_cmd := "ps -ef | grep -v grep |grep '" + constant.DST_CAVES + "' |sed -n '1P'|awk '{print $2}'"
	check, error := Shell(check_cmd)

	if error != nil || check == "" {
		log.Println("grep Caves pId error")
	} else {
		_, err := Shell(constant.STOP_CAVES_CMD)
		if err != nil {
			log.Panicln("stop caves error: " + err.Error())
		}
	}

}
func startCaves() {
	_, err := Shell(constant.GET_START_CAVES_CMD())
	if err != nil {
		log.Panicln("start caves error: " + err.Error())
	}
}

func getMasterStatus() bool {

	result, err := Shell(constant.FIND_MASTER_CMD)
	if err != nil {
		return false
	}
	res := strings.Split(result, "\n")[0]
	return res != ""
}

func getCavesStatus() bool {
	result, err := Shell(constant.FIND_CAVES_CMD)
	if err != nil {
		return false
	}
	res := strings.Split(result, "\n")[0]
	return res != ""
}

/**
* 检查目前所有的screen作业，并删除已经无法使用的screen作业
 */
func ClearScreen() bool {
	result, err := Shell(constant.CLEAR_SCREEN_CMD)
	if err != nil {
		return false
	}
	res := strings.Split(result, "\n")[0]
	return res != ""
}

func DeleteGameRecord() {
	DeleteCavesRecord()
	DeleteMasterRecord()
}

func DeleteMasterRecord() {
	stopMaster()
	Shell(constant.DEL_RECORD_MASTER_CMD)
}

func DeleteCavesRecord() {
	stopCaves()
	Shell(constant.DEL_RECORD_CAVES_CMD)
}

func SentBroadcast(message string) {

	broadcast := "screen -S \"" + constant.SCREEN_WORK_MASTER_NAME + "\" -p 0 -X stuff \"c_announce(\\\""
	broadcast += message
	broadcast += "\\\")\\n\""

	Shell(broadcast)
}

func KickPlayer(KuId string) {

	masterCMD := "screen -S \"" + constant.SCREEN_WORK_MASTER_NAME + "\" -p 0 -X stuff \"TheNet:Kick(\\\"" + KuId + "\\\")\\n\""
	cavesCMD := "screen -S \"" + constant.SCREEN_WORK_CAVES_NAME + "\" -p 0 -X stuff \"TheNet:Kick(\\\"" + KuId + "\\\")\\n\""

	Shell(masterCMD)
	Shell(cavesCMD)
}

func KillPlayer(KuId string) {

	masterCMD := "screen -S \"" + constant.SCREEN_WORK_MASTER_NAME + "\" -p 0 -X stuff \"UserToPlayer(\\\"" + KuId + "\\\"):PushEvent('death')\\n\""
	cavesCMD := "screen -S \"" + constant.SCREEN_WORK_CAVES_NAME + "\" -p 0 -X stuff \"UserToPlayer(\\\"" + KuId + "\\\"):PushEvent('death')\\n\""

	Shell(masterCMD)
	Shell(cavesCMD)
}

func RespawnPlayer(KuId string) {

	masterCMD := "screen -S \"" + constant.SCREEN_WORK_MASTER_NAME + "\" -p 0 -X stuff \"UserToPlayer(\\\"" + KuId + "\\\"):PushEvent('respawnfromghost')\\n\""
	cavesCMD := "screen -S \"" + constant.SCREEN_WORK_CAVES_NAME + "\" -p 0 -X stuff \"UserToPlayer(\\\"" + KuId + "\\\"):PushEvent('respawnfromghost')\\n\""

	Shell(masterCMD)
	Shell(cavesCMD)
}

func RollBack(dayNum int) {

	days := fmt.Sprint(dayNum)
	SentBroadcast(":pig 正在回档" + days + "天")

	masterCMD := "screen -S \"" + constant.SCREEN_WORK_MASTER_NAME + "\" -p 0 -X stuff \"c_rollback(" + days + ")\\n\""
	cavesCMD := "screen -S \"" + constant.SCREEN_WORK_CAVES_NAME + "\" -p 0 -X stuff \"c_rollback(" + days + ")\\n\""

	Shell(masterCMD)
	Shell(cavesCMD)
}

func Regenerateworld() {
	SentBroadcast(":pig 即将重置世界！！！")
	masterCMD := "screen -S \"" + constant.SCREEN_WORK_MASTER_NAME + "\" -p 0 -X stuff \"c_regenerateworld()\\n\""
	cavesCMD := "screen -S \"" + constant.SCREEN_WORK_CAVES_NAME + "\" -p 0 -X stuff \"c_regenerateworld()\\n\""
	Shell(masterCMD)
	Shell(cavesCMD)
}

func MasterConsole(command string) {
	cmd := "screen -S \"" + constant.SCREEN_WORK_MASTER_NAME + "\" -p 0 -X stuff \"" + command + "\\n\""
	Shell(cmd)
}

func CavesConsole(command string) {
	cmd := "screen -S \"" + constant.SCREEN_WORK_CAVES_NAME + "\" -p 0 -X stuff \"" + command + "\\n\""
	Shell(cmd)
}

func OperatePlayer(otype, kuId string) {
	command := ""
	//复活
	if otype == "0" {
		command = "UserToPlayer('%s'):PushEvent('respawnfromghost')"
	}
	//杀死
	if otype == "1" {
		command = "UserToPlayer('%s'):PushEvent('death')"
	}
	//更换角色
	if otype == "2" {
		command = "c_despawn('%s')"
	}
	MasterConsole(command)
	CavesConsole(command)
}

func GetGameArchive() *vo.GameArchiveVO {
	gameArchoeVO := vo.NewGameArchieVO()
	gameConfig := GetConfig()

	gameArchoeVO.ClusterName = gameConfig.ClusterName
	gameArchoeVO.ClusterPassword = gameConfig.ClusterPassword
	gameArchoeVO.GameMode = gameConfig.GameMode
	gameArchoeVO.MaxPlayers = gameConfig.MaxPlayers
	gameArchoeVO.Modoverrides = gameConfig.ModData

	workshopIds := WorkshopIds(gameConfig.ModData)
	gameArchoeVO.WorkshopIds = workshopIds
	gameArchoeVO.TotalModNum = len(workshopIds)

	return gameArchoeVO
}

func WorkshopIds(content string) []string {
	var workshopIds []string

	re := regexp.MustCompile("\"workshop-\\w[-\\w+]*\"")
	workshops := re.FindAllString(content, -1)

	for _, workshop := range workshops {
		workshop = strings.Replace(workshop, "\"", "", -1)
		split := strings.Split(workshop, "-")
		workshopId := strings.TrimSpace(split[1])
		workshopIds = append(workshopIds, workshopId)
	}
	return workshopIds
}

func PsAux(processName string) string {
	cmd := "ps -aux | grep -v grep |grep '" + processName + "' |sed -n '2P'|awk '{print $3,$4,$5,$6}'"
	res, err := Shell(cmd)
	if err != nil {
		log.Println("ps -aux |grep " + processName + " error: " + err.Error())
		return ""
	}
	return res
}

// 执行shell命令
func Shell(cmd string) (res string, err error) {
	var execCmd *exec.Cmd
	if runtime.GOOS == "windows" {
		execCmd = exec.Command("cmd.exe", "/c", cmd)
	} else {
		execCmd = exec.Command("bash", "-c", cmd)
	}
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	execCmd.Stdout = &stdout
	execCmd.Stderr = &stderr
	err = execCmd.Run()
	if err != nil {
		log.Println("error: " + err.Error())
	}

	output := ConvertByte2String(stderr.Bytes(), GB18030)
	errput := ConvertByte2String(stdout.Bytes(), GB18030)
	//res = fmt.Sprintf("Output:\n%s\nError:\n%s", stdout.String(), stderr.String())

	log.Printf("shell exec: %s \nOutput:\n%s\nError:\n%s", cmd, output, errput)

	return stdout.String(), err
}

// func execShell(shell string) {
// 	cmd := exec.Command(shell)
// 	e := cmd.Run()
// 	CheckError(e)
// }

// func execShellBin(shell string) {
// 	cmd := exec.Command("/bin/bash", "-c", shell)

// 	output, err := cmd.Output()
// 	if err != nil {
// 		fmt.Printf("Execute Shell:%s failed with error:%s", shell, err.Error())
// 		return
// 	}
// 	fmt.Printf("Execute Shell:%s finished with output:\n%s", shell, string(output))
// }

// func execShellSTD(shell string) []string {
// 	out, err := exec.Command(shell).Output()
// 	if err != nil {
// 		CheckError(err)
// 	}
// 	str := string(out)
// 	arr := strings.Split(str, "\n")
// 	return arr

// }

//	func CheckError(e error) {
//		if e != nil {
//			fmt.Println(e)
//		}
//	}
type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

func ConvertByte2String(byte []byte, charset Charset) string {

	var str string
	switch charset {
	case GB18030:
		decodeBytes, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}

	return str
}
