package dst

import (
	"dst-admin-go/constant"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/utils/shellUtils"
	"path"
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

func Status(clusterName, level string) bool {
	cmd := " ps -ef | grep -v grep | grep -v tail |grep '" + clusterName + "'|grep " + level + " |sed -n '1P'|awk '{print $2}' "
	result, err := shellUtils.Shell(cmd)
	if err != nil {
		return false
	}
	res := strings.Split(result, "\n")[0]
	return res != ""
}
