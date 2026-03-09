package game

import (
	"dst-admin-go/internal/service/dstConfig"
	"dst-admin-go/internal/service/levelConfig"
	"runtime"
)

func NewGame(dstConfig dstConfig.Config, levelConfigUtils *levelConfig.LevelConfigUtils) Process {
	if runtime.GOOS == "windows" {
		return NewWindowProcess(&dstConfig, levelConfigUtils)
	}
	return NewLinuxProcess(dstConfig, levelConfigUtils)
}
