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
	"strings"
)

type GameConsoleService struct {
	GameService
}

func (c *GameConsoleService) ClearScreen() bool {
	result, err := shellUtils.Shell(constant.CLEAR_SCREEN_CMD)
	if err != nil {
		return false
	}
	res := strings.Split(result, "\n")[0]
	return res != ""
}

func (c *GameConsoleService) SentBroadcast(clusterName string, message string) {

	if c.GetLevelStatus(clusterName, "Master") {
		broadcast := "screen -S \"" + screenKey.Key(clusterName, "Master") + "\" -p 0 -X stuff \"c_announce(\\\""
		broadcast += message
		broadcast += "\\\")\\n\""
		log.Println(broadcast)
		shellUtils.Shell(broadcast)
	}

	if c.GetLevelStatus(clusterName, "Caves") {
		broadcast2 := "screen -S \"" + screenKey.Key(clusterName, "Caves") + "\" -p 0 -X stuff \"c_announce(\\\""
		broadcast2 += message
		broadcast2 += "\\\")\\n\""
		log.Println(broadcast2)
		shellUtils.Shell(broadcast2)
	}

}

func (c *GameConsoleService) KickPlayer(clusterName, KuId string) {

	masterCMD := "screen -S \"" + screenKey.Key(clusterName, "Master") + "\" -p 0 -X stuff \"TheNet:Kick(\\\"" + KuId + "\\\")\\n\""
	cavesCMD := "screen -S \"" + screenKey.Key(clusterName, "Caves") + "\" -p 0 -X stuff \"TheNet:Kick(\\\"" + KuId + "\\\")\\n\""

	shellUtils.Shell(masterCMD)
	shellUtils.Shell(cavesCMD)
}

func (c *GameConsoleService) KillPlayer(clusterName, KuId string) {
	masterCMD := "screen -S \"" + screenKey.Key(clusterName, "Master") + "\" -p 0 -X stuff \"UserToPlayer(\\\"" + KuId + "\\\"):PushEvent('death')\\n\""
	cavesCMD := "screen -S \"" + screenKey.Key(clusterName, "Caves") + "\" -p 0 -X stuff \"UserToPlayer(\\\"" + KuId + "\\\"):PushEvent('death')\\n\""

	shellUtils.Shell(masterCMD)
	shellUtils.Shell(cavesCMD)
}

func (c *GameConsoleService) RespawnPlayer(clusterName string, KuId string) {

	masterCMD := "screen -S \"" + screenKey.Key(clusterName, "Master") + "\" -p 0 -X stuff \"UserToPlayer(\\\"" + KuId + "\\\"):PushEvent('respawnfromghost')\\n\""
	cavesCMD := "screen -S \"" + screenKey.Key(clusterName, "Caves") + "\" -p 0 -X stuff \"UserToPlayer(\\\"" + KuId + "\\\"):PushEvent('respawnfromghost')\\n\""

	shellUtils.Shell(masterCMD)
	shellUtils.Shell(cavesCMD)
}

func (c *GameConsoleService) RollBack(clusterName string, dayNum int) {
	days := fmt.Sprint(dayNum)
	c.SentBroadcast(clusterName, ":pig 正在回档"+days+"天")

	masterCMD := "screen -S \"" + screenKey.Key(clusterName, "Master") + "\" -p 0 -X stuff \"c_rollback(" + days + ")\\n\""
	cavesCMD := "screen -S \"" + screenKey.Key(clusterName, "Caves") + "\" -p 0 -X stuff \"c_rollback(" + days + ")\\n\""

	shellUtils.Shell(masterCMD)
	shellUtils.Shell(cavesCMD)
}

func (c *GameConsoleService) CleanWorld(clusterName string) {

	basePath := dst.GetClusterBasePath(clusterName)

	fileUtils.DeleteDir(path.Join(basePath, "Master", "backup"))
	fileUtils.DeleteDir(path.Join(basePath, "Master", "save"))

	fileUtils.DeleteDir(path.Join(basePath, "Caves", "backup"))
	fileUtils.DeleteDir(path.Join(basePath, "Caves", "save"))
}

func (c *GameConsoleService) Regenerateworld(clusterName string) {

	c.SentBroadcast(clusterName, ":pig 即将重置世界！！！")

	masterCMD := "screen -S \"" + screenKey.Key(clusterName, "Master") + "\" -p 0 -X stuff \"c_regenerateworld()\\n\""
	cavesCMD := "screen -S \"" + screenKey.Key(clusterName, "Caves") + "\" -p 0 -X stuff \"c_regenerateworld()\\n\""
	shellUtils.Shell(masterCMD)
	shellUtils.Shell(cavesCMD)
}

func (c *GameConsoleService) MasterConsole(clusterName string, command string) {

	cmd := "screen -S \"" + screenKey.Key(clusterName, "Master") + "\" -p 0 -X stuff \"" + command + "\\n\""
	shellUtils.Shell(cmd)
}

func (c *GameConsoleService) CavesConsole(clusterName string, command string) {

	cmd := "screen -S \"" + screenKey.Key(clusterName, "Master") + "\" -p 0 -X stuff \"" + command + "\\n\""
	shellUtils.Shell(cmd)
}

func (c *GameConsoleService) OperatePlayer(clusterName string, otype, kuId string) {
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
	c.MasterConsole(clusterName, command)
	c.CavesConsole(clusterName, command)
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

func (c *GameConsoleService) ReadLevelServerLog(clusterName, levelName string, length uint) []string {
	// levelServerIniPath := dst.GetLevelServerIniPath(clusterName, levelName)
	serverLogPath := dst.GetLevelServerLogPath(clusterName, levelName)
	lines, err := fileUtils.ReverseRead(serverLogPath, length)
	if err != nil {
		log.Panicln("读取日志server_log失败")
	}
	return lines
}

func (c *GameConsoleService) ReadLevelServerChatLog(clusterName, levelName string, length uint) []string {
	// levelServerIniPath := dst.GetLevelServerIniPath(clusterName, levelName)
	serverChatLogPath := dst.GetLevelServerChatLogPath(clusterName, levelName)
	lines, err := fileUtils.ReverseRead(serverChatLogPath, length)
	if err != nil {
		log.Panicln("读取日志server_chat_log失败")
	}
	return lines
}
