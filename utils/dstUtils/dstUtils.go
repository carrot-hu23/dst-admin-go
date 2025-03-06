package dstUtils

import (
	"bytes"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/utils/dstConfigUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/utils/shellUtils"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	textTemplate "text/template"
)

func GetUgcWorkshopModPath(clusterName, levelName, workshopId string) string {
	dstConfig := dstConfigUtils.GetDstConfig()
	workshopModPath := ""
	if dstConfig.Ugc_directory != "" {
		workshopModPath = filepath.Join(GetUgcModPath(), "content", "322330", workshopId)
	} else {
		workshopModPath = filepath.Join(dstConfig.Force_install_dir, "ugc_mods", clusterName, levelName, "content", "322330", workshopId)
	}
	return workshopModPath
}

func GetUgcModPath() string {
	dstConfig := dstConfigUtils.GetDstConfig()
	ugcModPath := ""
	if dstConfig.Ugc_directory != "" {
		ugcModPath = dstConfig.Ugc_directory
	} else {
		ugcModPath = filepath.Join(dstConfig.Force_install_dir, "ugc_mods")
	}
	return ugcModPath
}

func GetUgcAcfPath(clusterName, levelName string) string {
	ugcModPath := GetUgcModPath()
	dstConfig := dstConfigUtils.GetDstConfig()
	p := ""
	if dstConfig.Ugc_directory == "" {
		p = filepath.Join(ugcModPath, clusterName, levelName, "appworkshop_322330.acf")
	} else {
		p = filepath.Join(ugcModPath, "appworkshop_322330.acf")
	}
	return p
}

func GetKleiDstPath() string {
	dstConfig := dstConfigUtils.GetDstConfig()
	confDir := dstConfig.Conf_dir
	persistentStorageRoot := dstConfig.Persistent_storage_root
	kleiDstPath := ""
	if persistentStorageRoot == "" {
		kleiDstPath = dstConfigUtils.KleiDstPath()
	} else {
		if confDir == "" {
			confDir = "DoNotStarveTogether"
		}
		kleiDstPath = filepath.Join(persistentStorageRoot, confDir)
	}
	return kleiDstPath
}

func GetBlacklistPath(clusterName string) string {
	return filepath.Join(GetKleiDstPath(), clusterName, "blocklist.txt")
}

func GetWhitelistPath(clusterName string) string {
	return filepath.Join(GetKleiDstPath(), clusterName, "whitelist.txt")
}

func GetLevelLeveldataoverridePath(clusterName string, levelName string) string {
	return filepath.Join(GetKleiDstPath(), clusterName, levelName, "leveldataoverride.lua")
}

func GetLevelModoverridesPath(clusterName string, levelName string) string {
	return filepath.Join(GetKleiDstPath(), clusterName, levelName, "modoverrides.lua")
}

func GetLevelServerIniPath(clusterName string, levelName string) string {
	return filepath.Join(GetKleiDstPath(), clusterName, levelName, "server.ini")
}

func GetLevelServerLogPath(clusterName string, levelName string) string {
	return filepath.Join(GetKleiDstPath(), clusterName, levelName, "server_log.txt")
}

func GetLevelServerChatLogPath(clusterName string, levelName string) string {
	return filepath.Join(GetKleiDstPath(), clusterName, levelName, "server_chat_log.txt")
}

func GetClusterBasePath(clusterName string) string {
	return filepath.Join(GetKleiDstPath(), clusterName)
}

func GetClusterIniPath(clusterName string) string {
	return filepath.Join(GetKleiDstPath(), clusterName, "cluster.ini")
}

func GetClusterTokenPath(clusterName string) string {
	return filepath.Join(GetKleiDstPath(), clusterName, "cluster_token.txt")
}

func GetMasterModoverridesPath(clusterName string) string {
	return filepath.Join(GetKleiDstPath(), clusterName, "Master", "modoverrides.lua")
}

func GetCavesModoverridesPath(clusterName string) string {
	return filepath.Join(GetKleiDstPath(), clusterName, "Caves", "modoverrides.lua")
}

func GetMasterLeveldataoverridePath(clusterName string) string {
	return filepath.Join(GetKleiDstPath(), clusterName, "Master", "leveldataoverride.lua")
}
func GetCavesLeveldataoverridePath(clusterName string) string {
	return filepath.Join(GetKleiDstPath(), clusterName, "Caves", "leveldataoverride.lua")
}

func GetMasterServerIniPath(clusterName string) string {
	return filepath.Join(GetKleiDstPath(), clusterName, "Master", "server.ini")
}

func GetCavesServerIniPath(clusterName string) string {
	return filepath.Join(GetKleiDstPath(), clusterName, "Master", "server.ini")
}

func GetAdminlistPath(clusterName string) string {
	return filepath.Join(GetKleiDstPath(), clusterName, "adminlist.txt")
}

func GetBlocklistPath(clusterName string) string {
	return filepath.Join(GetKleiDstPath(), clusterName, "blocklist.txt")
}

func GetModSetup(clusterName string) string {
	cluster := dstConfigUtils.GetDstConfig()
	dstServerPath := cluster.Force_install_dir
	if dstConfigUtils.IsBeta() {
		dstServerPath = dstServerPath + "-beta"
	}
	if runtime.GOOS == "darwin" {
		return filepath.Join(dstServerPath, "dontstarve_dedicated_server_nullrenderer.app", "Contents", "mods", "dedicated_server_mods_setup.lua")
	}
	return filepath.Join(dstServerPath, "mods", "dedicated_server_mods_setup.lua")
}

//func GetDstUpdateCmd(clusterName string) string {
//	cluster := dstConfigUtils.GetDstConfig()
//	steamcmd := cluster.Steamcmd
//	dst_install_dir := cluster.Force_install_dir
//
//	dst_install_dir = EscapePath(dst_install_dir)
//
//	if runtime.GOOS == "windows" {
//		return "cd /d " + steamcmd + " && Start steamcmd.exe +login anonymous +force_install_dir " + dst_install_dir + " +app_update 343050 validate +quit"
//	}
//	if !fileUtils.Exists(filepath.Join(steamcmd, "steamcmd.sh")) {
//		return "cd " + steamcmd + " ; ./steamcmd +login anonymous +force_install_dir " + dst_install_dir + " +app_update 343050 validate +quit"
//	}
//	return "cd " + steamcmd + " ; ./steamcmd.sh +login anonymous +force_install_dir " + dst_install_dir + " +app_update 343050 validate +quit"
//}

func GetDstUpdateCmd(clusterName string) string {
	cluster := dstConfigUtils.GetDstConfig()
	steamCmdPath := cluster.Steamcmd
	dstInstallDir := cluster.Force_install_dir
	if cluster.Beta == 1 {
		dstInstallDir = dstInstallDir + "-beta"
	}
	// 确保路径是跨平台兼容的
	dstInstallDir = filepath.Clean(EscapePath(dstInstallDir))
	steamCmdPath = filepath.Clean(steamCmdPath)

	// 构建基本命令
	baseCmd := "+login anonymous +force_install_dir %s +app_update 343050"
	if cluster.Beta == 1 {
		baseCmd += " -beta updatebeta"
	}
	baseCmd += " validate +quit"
	baseCmd = fmt.Sprintf(baseCmd, dstInstallDir)

	// 根据操作系统生成不同的命令
	var cmd string
	if runtime.GOOS == "windows" {
		cmd = fmt.Sprintf("cd /d %s && Start steamcmd.exe %s", steamCmdPath, baseCmd)
	} else {
		steamCmdScript := filepath.Join(steamCmdPath, "steamcmd.sh")
		if fileUtils.Exists(steamCmdScript) {
			cmd = fmt.Sprintf("cd %s ; ./steamcmd.sh %s", steamCmdPath, baseCmd)
		} else {
			cmd = fmt.Sprintf("cd %s ; ./steamcmd %s", steamCmdPath, baseCmd)
		}
	}
	return cmd
}

func Status(clusterName, level string) bool {
	defer func() {
		if err := recover(); err != nil {

		}
	}()

	cmd := " ps -ef | grep -v grep | grep -v tail |grep '" + clusterName + "'|grep " + level + " |sed -n '1P'|awk '{print $2}' "
	result, err := shellUtils.Shell(cmd)
	if err != nil {
		return false
	}
	res := strings.Split(result, "\n")[0]
	return res != ""
}

func ReadMasterLog(clusterName string, lineNum uint) []string {
	logPath := filepath.Join(GetClusterBasePath(clusterName), "Master", "server_log.txt")
	logs, err := fileUtils.ReverseRead(logPath, lineNum)
	if err != nil {
		log.Panicln("read dstUtils2 master log error:", err)
	}
	return logs
}

func ReadLevelLog(clusterName string, levelName string, lineNum uint) []string {
	logPath := filepath.Join(GetClusterBasePath(clusterName), levelName, "server_log.txt")
	logs, err := fileUtils.ReverseRead(logPath, lineNum)
	if err != nil {
		log.Panicln("read dstUtils2 master log error:", err)
	}
	return logs
}

func ReadCavesLog(clusterName string, lineNum uint) []string {

	logPath := filepath.Join(GetClusterBasePath(clusterName), "Caves", "server_log.txt")
	logs, err := fileUtils.ReverseRead(logPath, lineNum)
	if err != nil {
		log.Panicln("read dstUtils2 caves log error:", err)
	}
	return logs
}

func ClearScreen() bool {
	result, err := shellUtils.Shell("screen -wipe")
	if err != nil {
		return false
	}
	res := strings.Split(result, "\n")[0]
	return res != ""
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

func GetPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org/")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(ip), nil
}

func Key(level, clusterName string) string {
	return "DST_" + level + "_" + clusterName
}

func ParseTemplate(templatePath string, data interface{}) string {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, data)
	if err != nil {
		log.Println(err)
		panic("模板解析错误")
	}
	return buf.String()
}

func ParseTemplate2(templatePath string, data interface{}) string {

	// 读取文件内容
	content, err := ioutil.ReadFile(templatePath)
	if err != nil {
		panic(err)
	}

	// 创建模板对象
	tmpl, err := textTemplate.New("myTemplate").Parse(string(content))
	if err != nil {
		panic(err)
	}

	// 执行模板并保存结果到字符串
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, data)
	if err != nil {
		panic(err)
	}
	return buf.String()

}

func DedicatedServerModsSetup(clusterName string, modConfig string) {
	if modConfig != "" {
		var serverModSetup = ""
		workshopIds := WorkshopIds(modConfig)
		for _, workshopId := range workshopIds {
			serverModSetup += "ServerModSetup(\"" + workshopId + "\")\n"
		}
		fileUtils.WriterTXT(GetModSetup2(clusterName), serverModSetup)
	}
}

func DedicatedServerModsSetup2(clusterName string, modConfig string) error {
	if modConfig != "" {
		var serverModSetup []string
		workshopIds := WorkshopIds(modConfig)
		for _, workshopId := range workshopIds {
			serverModSetup = append(serverModSetup, "ServerModSetup(\""+workshopId+"\")")
		}

		modSetupPath := GetModSetup2(clusterName)
		mods, err := fileUtils.ReadLnFile(modSetupPath)
		if err != nil {
			return err
		}
		var newServerModSetup []string
		for i := range serverModSetup {
			var notFind = true
			for j := range mods {
				if serverModSetup[i] == mods[j] {
					notFind = false
					break
				}
			}
			if notFind {
				newServerModSetup = append(newServerModSetup, serverModSetup[i])
			}
		}
		newServerModSetup = append(newServerModSetup, mods...)
		err = fileUtils.WriterLnFile(modSetupPath, newServerModSetup)
		if err != nil {
			return err
		}
	}
	return nil
}

type WorkshopItem struct {
	TimeUpdated int64
	Manifest    string
	Ugchandle   string
}

func ParseACFFile(filePath string) map[string]WorkshopItem {

	lines, err := fileUtils.ReadLnFile(filePath)
	if err != nil {
		log.Println(err)
		return nil
	}
	parsingWorkshopItemsInstalled := false
	workshopItems := make(map[string]WorkshopItem)
	var currentItemID string
	var currentItem WorkshopItem
	for _, line := range lines {
		// log.Println(line)
		if strings.Contains(line, "WorkshopItemsInstalled") {
			parsingWorkshopItemsInstalled = true
			continue
		}

		if strings.Contains(line, "{") && parsingWorkshopItemsInstalled {
			continue
		}

		if strings.Contains(line, "}") {
			continue
		}

		if parsingWorkshopItemsInstalled {
			replace := strings.Replace(line, "\t\t", "", -1)
			replace = strings.Replace(replace, "\"", "", -1)
			if _, err := strconv.Atoi(replace); err == nil {
				// This line contains the Workshop Item ID
				// currentItemID = line
				fields := strings.Fields(line)
				value := strings.Replace(fields[0], "\"", "", -1)
				currentItemID = value
			} else {
				// This line contains the Workshop Item details
				fields := strings.Fields(line)
				if len(fields) == 2 {
					key := strings.Replace(fields[0], "\"", "", -1)
					value := strings.Replace(fields[1], "\"", "", -1)
					// Remove double quotes from keys
					key = strings.ReplaceAll(key, "\"", "")
					switch key {
					case "timeupdated":
						currentItem.TimeUpdated, _ = strconv.ParseInt(value, 10, 64)
					case "manifest":
						currentItem.Manifest = strings.ReplaceAll(value, "\"", "")
					case "ugchandle":
						currentItem.Ugchandle = strings.ReplaceAll(value, "\"", "")
					}
				}
			}

			if currentItemID != "" && currentItem.TimeUpdated != 0 {
				workshopItems[currentItemID] = currentItem
				currentItemID = ""
				currentItem = WorkshopItem{}
			}
		}
	}

	return workshopItems
}

func GetModSetup2(clusterName string) string {
	cluster := clusterUtils.GetCluster(clusterName)
	return filepath.Join(cluster.ForceInstallDir, "mods", "dedicated_server_mods_setup.lua")
}

func EscapePath(path string) string {

	// windows 就跳过
	//if runtime.GOOS == "windows" {
	//	return path
	//}
	if runtime.GOOS == "windows" {
		return path
	}
	// 在这里添加需要转义的特殊字符
	escapedChars := []string{" ", "'", "(", ")"}

	for _, char := range escapedChars {
		path = strings.ReplaceAll(path, char, "\\"+char)
	}

	return path
}
