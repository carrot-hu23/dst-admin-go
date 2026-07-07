package dstPath

import (
	"dst-admin-go/internal/pkg/utils/fileUtils"
	"dst-admin-go/internal/service/dstConfig"
	"fmt"
	"path/filepath"
)

type LinuxDstPath struct {
	dstConfig dstConfig.Config
}

func NewLinuxDstPath(dstConfig dstConfig.Config) *LinuxDstPath {
	return &LinuxDstPath{dstConfig: dstConfig}
}

func (d LinuxDstPath) UpdateCommand(clusterName string) (string, error) {
	cluster, err := d.dstConfig.GetDstConfig(clusterName)
	if err != nil {
		return "", err
	}
	steamCmdPath := cluster.Steamcmd
	baseCmd := GetBaseUpdateCmd(cluster)
	var cmd string
	steamCmdScript := filepath.Join(steamCmdPath, "steamcmd.sh")
	if cluster.Bin == 86 {
		cmd = fmt.Sprintf("cd %s ; box86 ./linux32/steamcmd %s", steamCmdPath, baseCmd)
	} else {
		if fileUtils.Exists(steamCmdScript) {
			cmd = fmt.Sprintf("cd %s ; ./steamcmd.sh %s", steamCmdPath, baseCmd)
		} else {
			cmd = fmt.Sprintf("cd %s ; ./steamcmd %s", steamCmdPath, baseCmd)
		}
	}
	return cmd, nil
}
