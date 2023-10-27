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
	Execute(string, string)
}

type BackupStrategy struct{}

func (b *BackupStrategy) Execute(clusterName string, levelName string) {
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

func (u *UpdateStrategy) Execute(clusterName string, levelName string) {
	log.Println("正在定时更新游戏 clusterName: ", clusterName)
	err := gameService.UpdateGame(clusterName)
	if err != nil {
		log.Println("更新游戏失败: ", err)
		return
	}
}

type StartStrategy struct{}

func (s *StartStrategy) Execute(clusterName string, levelName string) {
	log.Println("正在定时启动游戏 clusterName: ", clusterName)
	cluster := clusterUtils.GetCluster(clusterName)
	gameService.LaunchLevel(clusterName, levelName, cluster.Bin, cluster.Beta)
}

type StopStrategy struct{}

func (s *StopStrategy) Execute(clusterName string, levelName string) {
	log.Println("正在定时关闭游戏 clusterName: ", clusterName)
	gameService.StopLevel(clusterName, levelName)
}

type RestartStrategy struct{}

func (s *RestartStrategy) Execute(clusterName string, levelName string) {
	log.Println("正在定时重启游戏 clusterName: ", clusterName)
	gameService.StopLevel(clusterName, levelName)
	cluster := clusterUtils.GetCluster(clusterName)
	gameService.LaunchLevel(clusterName, levelName, cluster.Bin, cluster.Beta)
}
