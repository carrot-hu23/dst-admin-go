package schedule

import (
	"dst-admin-go/constant/consts"
	"dst-admin-go/service"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/utils/zip"
	"log"
	"path/filepath"
)

var backupService = service.BackupService{}
var gameService = service.GameService{}

type Strategy interface {
	Execute(string)
}

type BackupStrategy struct{}

func (b *BackupStrategy) Execute(clusterName string) {
	cluster := clusterUtils.GetCluster(clusterName)
	src := filepath.Join(consts.KleiDstPath, cluster.ClusterName)

	dst := filepath.Join(cluster.Backup, backupService.GenGameBackUpName(clusterName))
	log.Println("正在定时创建游戏备份", "src: ", src, "dst: ", dst)
	err := zip.Zip(src, dst)
	if err != nil {
		log.Panicln("create backup error", err)
	}

}

type UpdateStrategy struct{}

func (u *UpdateStrategy) Execute(clusterName string) {
	log.Println("正在定时更新游戏 clusterName: ", clusterName)
	gameService.UpdateGame(clusterName)
}
