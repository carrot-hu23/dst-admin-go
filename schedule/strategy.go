package schedule

import (
	"dst-admin-go/service"
	"dst-admin-go/utils/clusterUtils"
	"log"
	"time"
)

var backupService = service.BackupService{}
var gameService = service.GameService{}

type Strategy interface {
	Execute(string, string)
}

type BackupStrategy struct{}

func (b *BackupStrategy) Execute(clusterName string, levelName string) {

	//consoleService.CSave(clusterName, "Master")
	//
	//// 保存存档
	//cluster := clusterUtils.GetCluster(clusterName)
	//src := filepath.Join(consts.KleiDstPath, cluster.ClusterName)
	//
	//dst := filepath.Join(cluster.Backup, backupService.GenGameBackUpName(clusterName))
	//log.Println("正在定时创建游戏备份", "src: ", src, "dst: ", dst)
	//err := zip.Zip(src, dst)
	//if err != nil {
	//	log.Println("create backup error", err)
	//}

	backupService.CreateBackup(clusterName, backupService.GenGameBackUpName(clusterName))

}

type UpdateStrategy struct{}

func (u *UpdateStrategy) Execute(clusterName string, levelName string) {
	log.Println("正在定时更新游戏 clusterName: ", clusterName)
	err := gameService.UpdateGame(clusterName)
	if err != nil {
		log.Println("更新游戏失败: ", err)
		return
	}
	time.Sleep(1 * time.Minute)
	gameService.StartGame(clusterName)
}

type StartStrategy struct{}

func (s *StartStrategy) Execute(clusterName string, levelName string) {
	log.Println("正在定时启动游戏 clusterName: ", clusterName)
	cluster := clusterUtils.GetCluster(clusterName)
	gameService.StartLevel(clusterName, levelName, cluster.Bin, cluster.Beta)
}

type StopStrategy struct{}

func (s *StopStrategy) Execute(clusterName string, levelName string) {
	log.Println("正在定时关闭游戏 clusterName: ", clusterName)
	gameService.StopLevel(clusterName, levelName)
}

type RestartStrategy struct{}

func (s *RestartStrategy) Execute(clusterName string, levelName string) {
	log.Println("正在定时重启游戏 clusterName: ", clusterName)
	cluster := clusterUtils.GetCluster(clusterName)
	gameService.StartLevel(clusterName, levelName, cluster.Bin, cluster.Beta)
}

type RegenerateStrategy struct{}

func (s *RegenerateStrategy) Execute(clusterName string, levelName string) {
	log.Println("正在定时重置游戏 clusterName: ", clusterName)
	gameConsoleService.Regenerateworld(clusterName)
}

type StartGameStrategy struct{}

func (s *StartGameStrategy) Execute(clusterName string, levelName string) {
	log.Println("正在定时启动游戏(所有) clusterName: ", clusterName)
	gameService.StartGame(clusterName)
}

type StopGameStrategy struct{}

func (s *StopGameStrategy) Execute(clusterName string, levelName string) {
	log.Println("正在定时关闭游戏(所有) clusterName: ", clusterName)
	gameService.StopGame(clusterName)
}
