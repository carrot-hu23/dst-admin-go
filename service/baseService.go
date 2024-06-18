package service

import (
	"runtime"
)

func isWindows() bool {
	return runtime.GOOS == "windows"
}

var WindowService WindowsGameService
var WindowGameConsoleService WindowsGameConsoleService

var clusterContainer = NewClusterContainer()
