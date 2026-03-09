package dstPath

import (
	"dst-admin-go/internal/service/dstConfig"
	"fmt"
)

type WindowDstPath struct {
	dstConfig dstConfig.Config
}

func NewWindowDst(dstConfig dstConfig.Config) *WindowDstPath {
	return &WindowDstPath{dstConfig: dstConfig}
}

func (d WindowDstPath) UpdateCommand(clusterName string) (string, error) {
	cluster, err := d.dstConfig.GetDstConfig(clusterName)
	if err != nil {
		return "", err
	}
	steamCmdPath := cluster.Steamcmd
	baseCmd := GetBaseUpdateCmd(cluster)

	var cmd string
	cmd = fmt.Sprintf("cd /d %s && Start steamcmd.exe %s", steamCmdPath, baseCmd)
	return cmd, nil
}
