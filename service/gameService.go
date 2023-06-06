package service

import (
	"dst-admin-go/constant"
	"dst-admin-go/constant/dst"
	"dst-admin-go/constant/screenKey"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/utils/shellUtils"
	"fmt"
	"log"
	"path"
	"regexp"
	"strings"
	"time"
)

type GameService struct {
	specifiedGameService SpecifiedGameService
}

func (s *GameService) UpdateGame(clusterName string) {

	SentBroadcast(clusterName, ":pig 正在更新游戏......")
	time.Sleep(3 * time.Second)

	s.specifiedGameService.stopSpecifiedMaster(clusterName)
	s.specifiedGameService.stopSpecifiedCaves(clusterName)
	updateGameCMd := dst.GetDstUpdateCmd(clusterName)
	log.Println(updateGameCMd)
	_, err := shellUtils.Shell(updateGameCMd)
	if err != nil {
		log.Panicln("update game error: " + err.Error())
	}
}

func ClearScreen() bool {
	result, err := shellUtils.Shell(constant.CLEAR_SCREEN_CMD)
	if err != nil {
		return false
	}
	res := strings.Split(result, "\n")[0]
	return res != ""
}

func SentBroadcast(clusterName string, message string) {

	broadcast := "screen -S \"" + screenKey.Key(clusterName, "Master") + "\" -p 0 -X stuff \"c_announce(\\\""
	broadcast += message
	broadcast += "\\\")\\n\""

	shellUtils.Shell(broadcast)
}

func KickPlayer(clusterName, KuId string) {

	masterCMD := "screen -S \"" + screenKey.Key(clusterName, "Master") + "\" -p 0 -X stuff \"TheNet:Kick(\\\"" + KuId + "\\\")\\n\""
	cavesCMD := "screen -S \"" + screenKey.Key(clusterName, "Caves") + "\" -p 0 -X stuff \"TheNet:Kick(\\\"" + KuId + "\\\")\\n\""

	shellUtils.Shell(masterCMD)
	shellUtils.Shell(cavesCMD)
}

func KillPlayer(clusterName, KuId string) {
	masterCMD := "screen -S \"" + screenKey.Key(clusterName, "Master") + "\" -p 0 -X stuff \"UserToPlayer(\\\"" + KuId + "\\\"):PushEvent('death')\\n\""
	cavesCMD := "screen -S \"" + screenKey.Key(clusterName, "Caves") + "\" -p 0 -X stuff \"UserToPlayer(\\\"" + KuId + "\\\"):PushEvent('death')\\n\""

	shellUtils.Shell(masterCMD)
	shellUtils.Shell(cavesCMD)
}

func RespawnPlayer(clusterName string, KuId string) {

	masterCMD := "screen -S \"" + screenKey.Key(clusterName, "Master") + "\" -p 0 -X stuff \"UserToPlayer(\\\"" + KuId + "\\\"):PushEvent('respawnfromghost')\\n\""
	cavesCMD := "screen -S \"" + screenKey.Key(clusterName, "Caves") + "\" -p 0 -X stuff \"UserToPlayer(\\\"" + KuId + "\\\"):PushEvent('respawnfromghost')\\n\""

	shellUtils.Shell(masterCMD)
	shellUtils.Shell(cavesCMD)
}

func RollBack(clusterName string, dayNum int) {
	days := fmt.Sprint(dayNum)
	SentBroadcast(clusterName, ":pig 正在回档"+days+"天")

	masterCMD := "screen -S \"" + screenKey.Key(clusterName, "Master") + "\" -p 0 -X stuff \"c_rollback(" + days + ")\\n\""
	cavesCMD := "screen -S \"" + screenKey.Key(clusterName, "Caves") + "\" -p 0 -X stuff \"c_rollback(" + days + ")\\n\""

	shellUtils.Shell(masterCMD)
	shellUtils.Shell(cavesCMD)
}

func CleanWorld(clusterName string) {

	basePath := dst.GetClusterBasePath(clusterName)

	fileUtils.DeleteDir(path.Join(basePath, "Master", "backup"))
	fileUtils.DeleteDir(path.Join(basePath, "Master", "save"))

	fileUtils.DeleteDir(path.Join(basePath, "Caves", "backup"))
	fileUtils.DeleteDir(path.Join(basePath, "Caves", "save"))
}

func Regenerateworld(clusterName string) {

	SentBroadcast(clusterName, ":pig 即将重置世界！！！")

	masterCMD := "screen -S \"" + screenKey.Key(clusterName, "Master") + "\" -p 0 -X stuff \"c_regenerateworld()\\n\""
	cavesCMD := "screen -S \"" + screenKey.Key(clusterName, "Caves") + "\" -p 0 -X stuff \"c_regenerateworld()\\n\""
	shellUtils.Shell(masterCMD)
	shellUtils.Shell(cavesCMD)
}

func MasterConsole(clusterName string, command string) {

	cmd := "screen -S \"" + screenKey.Key(clusterName, "Master") + "\" -p 0 -X stuff \"" + command + "\\n\""
	shellUtils.Shell(cmd)
}

func CavesConsole(clusterName string, command string) {

	cmd := "screen -S \"" + screenKey.Key(clusterName, "Master") + "\" -p 0 -X stuff \"" + command + "\\n\""
	shellUtils.Shell(cmd)
}

func OperatePlayer(clusterName string, otype, kuId string) {
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
	MasterConsole(clusterName, command)
	CavesConsole(clusterName, command)
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
	res, err := shellUtils.Shell(cmd)
	if err != nil {
		log.Println("ps -aux |grep " + processName + " error: " + err.Error())
		return ""
	}
	return res
}
