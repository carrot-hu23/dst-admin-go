package update

import (
	"dst-admin-go/internal/pkg/utils/shellUtils"
	"dst-admin-go/internal/service/dstConfig"
	"log"
)

type WindowUpdate struct {
	dstConfig dstConfig.Config
}

func NewWindowUpdate(dstConfig dstConfig.Config) *WindowUpdate {
	return &WindowUpdate{
		dstConfig: dstConfig,
	}
}

func (u WindowUpdate) Update(clusterName string) error {
	config, err := u.dstConfig.GetDstConfig(clusterName)
	if err != nil {
		return err
	}
	updateCommand, err := WindowUpdateCommand(config)
	if err != nil {
		return err
	}
	log.Println("正在更新游戏", "cluster: ", clusterName, "command: ", updateCommand)
	result, err := shellUtils.ExecuteCommandInWin(updateCommand)
	log.Println(result)
	return err
}
