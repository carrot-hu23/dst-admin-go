package utils

import "runtime"

func IsWindow() bool {
	os := runtime.GOOS
	return os == "windows"
}
