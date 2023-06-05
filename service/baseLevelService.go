package service

import (
	"dst-admin-go/constant"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo/cluster"
	"log"
	"path/filepath"
)

type BaseLevelService struct {
	DstHelper
}

func (b *BaseLevelService) GetLevelList() []string {
	klei_path := filepath.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether")
	levellist, err := fileUtils.ListDirectories(klei_path)
	if err != nil {
		log.Println(err)
		return []string{}
	}
	return levellist
}

func (b *BaseLevelService) CreateBaseLevel(baseLevel *cluster.BaseLevel) {

	clusterName := baseLevel.ClusterName
	kleiPath := filepath.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether")
	baseLevelPath := filepath.Join(kleiPath, clusterName)

	fileUtils.CreateFileIfNotExists(baseLevelPath)

	clusterTokenPath := filepath.Join(baseLevelPath, "cluster_token.txt")
	clusterIniPath := filepath.Join(baseLevelPath, "cluster.ini")

	adminlistPath := filepath.Join(baseLevelPath, "adminlist.txt")
	blocklistPath := filepath.Join(baseLevelPath, "blocklist.txt")

	b.SaveBaseClusterToken(baseLevel.ClusterToken, clusterTokenPath)
	b.SaveBaseClusterIni(baseLevel.Cluster, clusterIniPath)
	b.SaveBaseAdminlist(baseLevel.Adminlist, adminlistPath)
	b.SaveBaseBlocklist(baseLevel.Blocklist, blocklistPath)
	b.SaveBaseMaster(baseLevel.Master, baseLevelPath)
	b.SaveBaseCavesr(baseLevel.Caves, baseLevelPath)

}

func (b *BaseLevelService) SaveBaseClusterIni(cluster *cluster.ClusterIni, clusterIniPath string) {
	fileUtils.CreateFileIfNotExists(clusterIniPath)
	fileUtils.WriterTXT(clusterIniPath, b.ParseTemplate(CLUSTER_INI_TEMPLATE, cluster))
}

func (b *BaseLevelService) SaveBaseClusterToken(token string, clusterTokenPath string) {
	fileUtils.CreateFileIfNotExists(clusterTokenPath)
	fileUtils.WriterTXT(clusterTokenPath, token)
}

func (b *BaseLevelService) SaveBaseAdminlist(str []string, adminlistPath string) {
	fileUtils.CreateFileIfNotExists(adminlistPath)
	fileUtils.WriterLnFile(adminlistPath, str)
}

func (b *BaseLevelService) SaveBaseBlocklist(str []string, blocklistPath string) {
	fileUtils.CreateFileIfNotExists(blocklistPath)
	fileUtils.WriterLnFile(blocklistPath, str)
}

func (b *BaseLevelService) SaveBaseMaster(world *cluster.World, baseLevelPath string) {

	l_path := filepath.Join(baseLevelPath, "Master", "leveldataoverride.lua")
	m_path := filepath.Join(baseLevelPath, "Master", "modoverrides.lua")
	s_path := filepath.Join(baseLevelPath, "Master", "server.ini")

	fileUtils.CreateFileIfNotExists(l_path)
	fileUtils.CreateFileIfNotExists(m_path)
	fileUtils.CreateFileIfNotExists(s_path)

	fileUtils.WriterTXT(l_path, world.Leveldataoverride)
	fileUtils.WriterTXT(m_path, world.Modoverrides)

	serverBuf := b.ParseTemplate(MASTER_SERVER_INI_TEMPLATE, world.ServerIni)

	fileUtils.WriterTXT(s_path, serverBuf)
}

func (b *BaseLevelService) SaveBaseCavesr(world *cluster.World, baseLevelPath string) {

	lPath := filepath.Join(baseLevelPath, "Caves", "leveldataoverride.lua")
	mPath := filepath.Join(baseLevelPath, "Caves", "modoverrides.lua")
	sPath := filepath.Join(baseLevelPath, "Caves", "server.ini")

	fileUtils.CreateFileIfNotExists(lPath)
	fileUtils.CreateFileIfNotExists(mPath)
	fileUtils.CreateFileIfNotExists(sPath)

	fileUtils.WriterTXT(lPath, world.Leveldataoverride)
	fileUtils.WriterTXT(mPath, world.Modoverrides)

	serverBuf := b.ParseTemplate(CAVES_SERVER_INI_TEMPLATE, world.ServerIni)

	fileUtils.WriterTXT(sPath, serverBuf)
}
