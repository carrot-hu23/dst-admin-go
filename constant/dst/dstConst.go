package dst

import (
	"dst-admin-go/constant"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/utils/shellUtils"
	"log"
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
	cluster := clusterUtils.GetCluster(clusterName)
	return path.Join(cluster.ForceInstallDir, "mods", "dedicated_server_mods_setup.lua")
}

func GetDstUpdateCmd(clusterName string) string {
	cluster := clusterUtils.GetCluster(clusterName)
	steamcmd := cluster.SteamCmd
	dst_install_dir := cluster.ForceInstallDir
	return "cd " + steamcmd + " ; ./steamcmd.sh +login anonymous +force_install_dir " + dst_install_dir + " +app_update 343050 validate +quit"
}

// ============= 工具类 以后放到别的位置 ================//

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
