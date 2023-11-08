package dstUtils2

import (
	"dst-admin-go/constant"
	"dst-admin-go/utils/clusterUtils"
	"path"
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
	return path.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether", clusterName, "Caves", "server.ini")
}

func GetAdminlistPath(clusterName string) string {
	return path.Join(constant.HOME_PATH, ".klei", "DoNotStarveTogether", clusterName, "adminlist.txt")
}

func GetBlocklistPath(clusterName string) string {
	return path.Join(constant.HOME_PATH, ".klei", "DoNotStarveTogether", clusterName, "blocklist.txt")
}

func GetBlacklistPath(clusterName string) string {
	return path.Join(constant.HOME_PATH, ".klei", "DoNotStarveTogether", clusterName, "blocklist.txt")
}

func GetWhitelistPath(clusterName string) string {
	return path.Join(constant.HOME_PATH, ".klei", "DoNotStarveTogether", clusterName, "whitelist.txt")
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
