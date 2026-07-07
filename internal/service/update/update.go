package update

import (
	"dst-admin-go/internal/pkg/utils/fileUtils"
	"dst-admin-go/internal/service/dstConfig"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

type Update interface {
	Update(clusterName string) error
}

func EscapePath(path string) string {
	if runtime.GOOS == "windows" {
		return path
	}
	// 在这里添加需要转义的特殊字符
	escapedChars := []string{" ", "'", "(", ")"}
	for _, char := range escapedChars {
		path = strings.ReplaceAll(path, char, "\\"+char)
	}
	return path
}

func GetBaseUpdateCmd(cluster dstConfig.DstConfig) string {
	steamCmdPath := cluster.Steamcmd
	dstInstallDir := cluster.Force_install_dir
	if cluster.Beta == 1 {
		dstInstallDir = dstInstallDir + "-beta"
	}
	// 确保路径是跨平台兼容的
	dstInstallDir = filepath.Clean(EscapePath(dstInstallDir))
	steamCmdPath = filepath.Clean(steamCmdPath)

	// 构建基本命令
	baseCmd := "+login anonymous +force_install_dir %s +app_update 343050"
	if cluster.Beta == 1 {
		baseCmd += " -beta updatebeta"
	}
	baseCmd += " validate +quit"
	baseCmd = fmt.Sprintf(baseCmd, dstInstallDir)
	return baseCmd
}

func WindowUpdateCommand(cluster dstConfig.DstConfig) (string, error) {
	steamCmdPath := cluster.Steamcmd
	baseCmd := GetBaseUpdateCmd(cluster)

	var cmd string
	cmd = fmt.Sprintf("cd /d %s && Start steamcmd.exe %s", steamCmdPath, baseCmd)
	return cmd, nil
}

func LinuxUpdateCommand(cluster dstConfig.DstConfig) (string, error) {
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
