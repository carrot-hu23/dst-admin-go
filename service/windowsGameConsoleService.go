package service

import (
	"dst-admin-go/constant/screenKey"
	dst_cli_window "dst-admin-go/dst-cli-window"
	"dst-admin-go/utils/dstUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/utils/shellUtils"
	"fmt"
	"log"
	"path"
	"time"
)

type WindowsGameConsoleService struct {
	WindowsGameService
}

func (c *WindowsGameConsoleService) SentBroadcast2(clusterName string, levelName string, message string) {

	if c.GetLevelStatus(clusterName, levelName) {
		broadcast := "c_announce(\\\""
		broadcast += message
		broadcast += "\\\")"
		log.Println(broadcast)
		dst_cli_window.DstCliClient.Command(clusterName, levelName, broadcast)
	}

}

func (c *WindowsGameConsoleService) SentBroadcast(clusterName string, message string) {

	if c.GetLevelStatus(clusterName, "Master") {
		broadcast := "c_announce(\\\""
		broadcast += message
		broadcast += "\\\")"
		log.Println(broadcast)
		shellUtils.Shell(broadcast)

		dst_cli_window.DstCliClient.Command(clusterName, "Master", broadcast)
	}

}

func (c *WindowsGameConsoleService) KickPlayer(clusterName, KuId string) {

	masterCMD := "TheNet:Kick(\\\"" + KuId + "\\\")"
	dst_cli_window.DstCliClient.Command(clusterName, "Master", masterCMD)

}

func (c *WindowsGameConsoleService) KillPlayer(clusterName, KuId string) {
	masterCMD := "UserToPlayer(\\\"" + KuId + "\\\"):PushEvent('death')"
	dst_cli_window.DstCliClient.Command(clusterName, "Master", masterCMD)
}

func (c *WindowsGameConsoleService) RespawnPlayer(clusterName string, KuId string) {

	masterCMD := "UserToPlayer(\\\"" + KuId + "\\\"):PushEvent('respawnfromghost')"
	dst_cli_window.DstCliClient.Command(clusterName, "Master", masterCMD)
}

func (c *WindowsGameConsoleService) RollBack(clusterName string, dayNum int) {
	days := fmt.Sprint(dayNum)

	masterCMD := "c_rollback(" + days + ")"
	dst_cli_window.DstCliClient.Command(clusterName, "Master", masterCMD)
}

func (c *WindowsGameConsoleService) CleanWorld(clusterName string) {

	basePath := dstUtils.GetClusterBasePath(clusterName)

	fileUtils.DeleteDir(path.Join(basePath, "Master", "backup"))
	fileUtils.DeleteDir(path.Join(basePath, "Master", "save"))

	fileUtils.DeleteDir(path.Join(basePath, "Caves", "backup"))
	fileUtils.DeleteDir(path.Join(basePath, "Caves", "save"))
}

func (c *WindowsGameConsoleService) Regenerateworld(clusterName string) {

	c.SentBroadcast(clusterName, ":pig 即将重置世界！！！")
	masterCMD := "c_regenerateworld()"
	dst_cli_window.DstCliClient.Command(clusterName, "Master", masterCMD)
}

func (c *WindowsGameConsoleService) MasterConsole(clusterName string, command string) {

	cmd := "screen -S \"" + screenKey.Key(clusterName, "Master") + "\" -p 0 -X stuff \"" + command + "\\n\""
	shellUtils.Shell(cmd)
}

func (c *WindowsGameConsoleService) CavesConsole(clusterName string, command string) {

	cmd := "screen -S \"" + screenKey.Key(clusterName, "Master") + "\" -p 0 -X stuff \"" + command + "\\n\""
	shellUtils.Shell(cmd)
}

func (c *WindowsGameConsoleService) OperatePlayer(clusterName string, otype, kuId string) {
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

func (c *WindowsGameConsoleService) ReadLevelServerLog(clusterName, levelName string, length uint) []string {
	// levelServerIniPath := dstUtils2.GetLevelServerIniPath(clusterName, levelName)
	serverLogPath := dstUtils.GetLevelServerLogPath(clusterName, levelName)
	lines, err := fileUtils.ReverseRead(serverLogPath, length)
	if err != nil {
		log.Panicln("读取日志server_log失败")
	}
	return lines
}

func (c *WindowsGameConsoleService) ReadLevelServerChatLog(clusterName, levelName string, length uint) []string {
	// levelServerIniPath := dstUtils2.GetLevelServerIniPath(clusterName, levelName)
	serverChatLogPath := dstUtils.GetLevelServerChatLogPath(clusterName, levelName)
	lines, err := fileUtils.ReverseRead(serverChatLogPath, length)
	if err != nil {
		log.Panicln("读取日志server_chat_log失败")
	}
	return lines
}

func (c *WindowsGameConsoleService) SendCommand(clusterName string, levelName string, command string) {
	dst_cli_window.DstCliClient.Command(clusterName, levelName, command)
}

func (c *WindowsGameConsoleService) CSave(clusterName string, levelName string) {
	log.Println("正在 s_save() 存档", clusterName, levelName)
	command := "c_save()"
	cmd := "screen -S \"" + screenKey.Key(clusterName, levelName) + "\" -p 0 -X stuff \"" + command + "\\n\""
	shellUtils.Shell(cmd)

	time.Sleep(5 * time.Second)
}

func (c *WindowsGameConsoleService) CSaveMaster(clusterName string) {
	c.CSave(clusterName, "Master")
}
