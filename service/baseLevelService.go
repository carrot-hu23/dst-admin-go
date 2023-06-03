package service

import (
	"dst-admin-go/constant"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo/cluster"
	"log"
	"path/filepath"
)

// TODO 多开
func GetLevellist() []string {
	klei_path := filepath.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether")
	levellist, err := fileUtils.ListDirectories(klei_path)
	if err != nil {
		log.Println(err)
		return []string{}
	}
	return levellist
}

func CreateBaseLevel(baseLevel *cluster.BaseLevel) {

	clusterName := baseLevel.ClusterName
	klei_path := filepath.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether")
	baseLevelPath := filepath.Join(klei_path, clusterName)

	createFileIfNotExsists(baseLevelPath)

	cluster_token_path := filepath.Join(baseLevelPath, "cluster_token.txt")
	cluster_ini_path := filepath.Join(baseLevelPath, "cluster.ini")

	adminlist_path := filepath.Join(baseLevelPath, "adminlist.txt")
	blocklist_path := filepath.Join(baseLevelPath, "blocklist.txt")

	// master_server_ini_path := filepath.Join(baseLevelPath, "Master", "server.ini")
	// master_leveldataoverride_path := filepath.Join(baseLevelPath, "Master", "leveldataoverride")
	// master_modoverrides_path := filepath.Join(baseLevelPath, "Master", "modoverrides")
	// caves_server_ini_path := filepath.Join(baseLevelPath, "Caves", "server.ini")
	// caves_leveldataoverride_path := filepath.Join(baseLevelPath, "Caves", "leveldataoverride")
	// caves_modoverrides_path := filepath.Join(baseLevelPath, "Caves", "modoverrides")
	// log.Println(cluster_token_path, cluster_ini_path, adminlist_path, blocklist_path)
	// log.Println(master_server_ini_path, master_leveldataoverride_path, master_modoverrides_path)
	// log.Println(caves_server_ini_path, caves_leveldataoverride_path, caves_modoverrides_path)

	SaveBaseClusterToken(baseLevel.ClusterToken, cluster_token_path)
	SaveBaseClusterIni(baseLevel.Cluster, cluster_ini_path)
	SaveBaseAdminlist(baseLevel.Adminlist, adminlist_path)
	SaveBaseBlocklist(baseLevel.Blocklist, blocklist_path)
	SaveBaseMaster(baseLevel.Master, baseLevelPath)
	SaveBaseCavesr(baseLevel.Caves, baseLevelPath)

}

func SaveBaseClusterIni(cluster *cluster.Cluster, cluster_ini_path string) {
	createFileIfNotExsists(cluster_ini_path)
	fileUtils.WriterTXT(cluster_ini_path, pareseTemplate(CLUSTER_INI_TEMPLATE, cluster))
}

func SaveBaseClusterToken(token string, cluster_token_path string) {
	createFileIfNotExsists(cluster_token_path)
	fileUtils.WriterTXT(cluster_token_path, token)
}

func SaveBaseAdminlist(str []string, adminlist_path string) {
	createFileIfNotExsists(adminlist_path)
	fileUtils.WriterLnFile(adminlist_path, str)
}

func SaveBaseBlocklist(str []string, blocklist_path string) {
	createFileIfNotExsists(blocklist_path)
	fileUtils.WriterLnFile(blocklist_path, str)
}

func SaveBaseMaster(world *cluster.World, baseLevelPath string) {

	l_path := filepath.Join(baseLevelPath, "Master", "leveldataoverride.lua")
	m_path := filepath.Join(baseLevelPath, "Master", "modoverrides.lua")
	s_path := filepath.Join(baseLevelPath, "Master", "server.ini")

	createFileIfNotExsists(l_path)
	createFileIfNotExsists(m_path)
	createFileIfNotExsists(s_path)

	fileUtils.WriterTXT(l_path, world.Leveldataoverride)
	fileUtils.WriterTXT(m_path, world.Modoverrides)

	serverBuf := pareseTemplate(MASTER_SERVER_INI_TEMPLATE, world.ServerIni)

	fileUtils.WriterTXT(s_path, serverBuf)
}

func SaveBaseCavesr(world *cluster.World, baseLevelPath string) {

	l_path := filepath.Join(baseLevelPath, "Caves", "leveldataoverride.lua")
	m_path := filepath.Join(baseLevelPath, "Caves", "modoverrides.lua")
	s_path := filepath.Join(baseLevelPath, "Caves", "server.ini")

	createFileIfNotExsists(l_path)
	createFileIfNotExsists(m_path)
	createFileIfNotExsists(s_path)

	fileUtils.WriterTXT(l_path, world.Leveldataoverride)
	fileUtils.WriterTXT(m_path, world.Modoverrides)

	serverBuf := pareseTemplate(CAVES_SERVER_INI_TEMPLATE, world.ServerIni)

	fileUtils.WriterTXT(s_path, serverBuf)
}
