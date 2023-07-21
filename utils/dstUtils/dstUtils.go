package dstUtils

import (
	"bytes"
	"dst-admin-go/constant"
	"dst-admin-go/constant/dst"
	"dst-admin-go/utils/dstConfigUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/utils/shellUtils"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"regexp"
	"strings"
)

func GetClusterBasePath(clusterName string) string {
	return path.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether", clusterName)
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
		log.Panicln("read dst master log error:", err)
	}
	return logs
}

func ReadCavesLog(clusterName string, lineNum uint) []string {

	logPath := path.Join(GetClusterBasePath(clusterName), "Caves", "server_log.txt")
	logs, err := fileUtils.ReverseRead(logPath, lineNum)
	if err != nil {
		log.Panicln("read dst caves log error:", err)
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
		fileUtils.WriterTXT(dst.GetModSetup(clusterName), serverModSetup)
	}
}
