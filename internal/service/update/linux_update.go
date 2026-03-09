package update

import (
	"dst-admin-go/internal/pkg/utils/shellUtils"
	"dst-admin-go/internal/service/dstConfig"
	"log"
)

type LinuxUpdate struct {
	dstConfig dstConfig.Config
}

func NewLinuxUpdate(dstConfig dstConfig.Config) *LinuxUpdate {
	return &LinuxUpdate{
		dstConfig: dstConfig,
	}
}

func (u LinuxUpdate) Update(clusterName string) error {
	config, err := u.dstConfig.GetDstConfig(clusterName)
	if err != nil {
		return err
	}
	updateCommand, err := LinuxUpdateCommand(config)
	if err != nil {
		return err
	}
	log.Println("正在更新游戏", "cluster: ", clusterName, "command: ", updateCommand)
	_, err = shellUtils.Shell(updateCommand)
	return err
}
