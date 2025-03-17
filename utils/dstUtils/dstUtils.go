package dstUtils

import (
	"bytes"
	"dst-admin-go/constant"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/utils/dstConfigUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/utils/shellUtils"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

func GetBlacklistPath(clusterName string) string {
	return path.Join(constant.HOME_PATH, ".klei", "DoNotStarveTogether", clusterName, "blocklist.txt")
}

func GetWhitelistPath(clusterName string) string {
	return path.Join(constant.HOME_PATH, ".klei", "DoNotStarveTogether", clusterName, "whitelist.txt")
}

func GetLevelLeveldataoverridePath(clusterName string, levelName string) string {
	return path.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether", clusterName, levelName, "leveldataoverride.lua")
}

func GetLevelModoverridesPath(clusterName string, levelName string) string {
	return path.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether", clusterName, levelName, "modoverrides.lua")
}

func GetLevelServerIniPath(clusterName string, levelName string) string {
	return path.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether", clusterName, levelName, "server.ini")
}

func GetLevelServerLogPath(clusterName string, levelName string) string {
	return path.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether", clusterName, levelName, "server_log.txt")
}

func GetLevelServerChatLogPath(clusterName string, levelName string) string {
	return path.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether", clusterName, levelName, "server_chat_log.txt")
}

func GetClusterBasePath(clusterName string) string {
	return path.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether", clusterName)
}
func GetKleiDstPath() string {
	return path.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether")
}

func GetDoNotStarveTogetherPath() string {
	return path.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether")
}

func GetClusterIniPath(clusterName string) string {
	return path.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether", clusterName, "cluster.ini")
}

func GetClusterTokenPath(clusterName string) string {
	return path.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether", clusterName, "cluster_token.txt")
}

func GetMasterModoverridesPath(clusterName string) string {
	return path.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether", clusterName, "Master", "modoverrides.lua")
}

func GetCavesModoverridesPath(clusterName string) string {
	return path.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether", clusterName, "Caves", "modoverrides.lua")
}

func GetMasterLeveldataoverridePath(clusterName string) string {
	return path.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether", clusterName, "Master", "leveldataoverride.lua")
}
func GetCavesLeveldataoverridePath(clusterName string) string {
	return path.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether", clusterName, "Caves", "leveldataoverride.lua")
}

func GetMasterServerIniPath(clusterName string) string {
	return path.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether", clusterName, "Master", "server.ini")
}

func GetCavesServerIniPath(clusterName string) string {
	return path.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether", clusterName, "Master", "server.ini")
}

func GetAdminlistPath(clusterName string) string {
	return path.Join(constant.HOME_PATH, ".klei", "DoNotStarveTogether", clusterName, "adminlist.txt")
}

func GetBlocklistPath(clusterName string) string {
	return path.Join(constant.HOME_PATH, ".klei", "DoNotStarveTogether", clusterName, "blocklist.txt")
}

func GetModSetup(clusterName string) string {
	cluster := dstConfigUtils.GetDstConfig()
	return path.Join(cluster.Force_install_dir, "mods", "dedicated_server_mods_setup.lua")
}

func GetDstUpdateCmd(clusterName string) string {
	cluster := dstConfigUtils.GetDstConfig()
	steamcmd := cluster.Steamcmd
	dst_install_dir := cluster.Force_install_dir
	return "cd " + steamcmd + " ; ./steamcmd.sh +login anonymous +force_install_dir " + dst_install_dir + " +app_update 343050 validate +quit"
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
	logPath := path.Join(GetClusterBasePath(clusterName), "Master", "server_log.txt")
	logs, err := fileUtils.ReverseRead(logPath, lineNum)
	if err != nil {
		log.Panicln("read dstUtils2 master log error:", err)
	}
	return logs
}

func ReadLevelLog(clusterName string, levelName string, lineNum uint) []string {
	logPath := path.Join(GetClusterBasePath(clusterName), levelName, "server_log.txt")
	logs, err := fileUtils.ReverseRead(logPath, lineNum)
	if err != nil {
		log.Panicln("read dstUtils2 master log error:", err)
	}
	return logs
}

func ReadCavesLog(clusterName string, lineNum uint) []string {

	logPath := path.Join(GetClusterBasePath(clusterName), "Caves", "server_log.txt")
	logs, err := fileUtils.ReverseRead(logPath, lineNum)
	if err != nil {
		log.Panicln("read dstUtils2 caves log error:", err)
	}
	return logs
}

func ClearScreen() bool {
	result, err := shellUtils.Shell(constant.CLEAR_SCREEN_CMD)
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

func DedicatedServerModsSetup2(clusterName string, modConfig string) {
	if modConfig != "" {
		var serverModSetup []string
		workshopIds := WorkshopIds(modConfig)
		for _, workshopId := range workshopIds {
			serverModSetup = append(serverModSetup, "ServerModSetup(\""+workshopId+"\")")
		}

		modSetupPath := GetModSetup2(clusterName)
		mods, err := fileUtils.ReadLnFile(modSetupPath)
		if err != nil {
			log.Panicln("读取 dedicated_server_mods_setup.lua 失败", err)
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
			log.Panicln("写入 dedicated_server_mods_setup.lua 失败", err)
		}
	}
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
	return path.Join(cluster.ForceInstallDir, "mods", "dedicated_server_mods_setup.lua")
}

func EscapePath(path string) string {
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
