package schedule

import "log"

type Strategy interface {
	Execute(string)
}

type BackupStrategy struct{}

func (b *BackupStrategy) Execute(clusterName string) {
	log.Println("正在创建备份 clusterName: ", clusterName)
}

type UpdateStrategy struct{}

func (u *UpdateStrategy) Execute(clusterName string) {
	log.Println("正在更新游戏 clusterName: ", clusterName)
}

type LobbyServerStrategy struct{}

func (l *LobbyServerStrategy) Execute(clusterName string) {

}
