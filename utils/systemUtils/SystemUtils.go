package systemUtils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

// Home returns the home directory for the executing user.
//
// This uses an OS-specific method for discovering the home directory.
// An error is returned if a home directory cannot be detected.
func Home() (string, error) {
	user, err := user.Current()
	if nil == err {
		return user.HomeDir, nil
	}

	// cross compile support

	if runtime.GOOS == "windows" {
		return homeWindows()
	}

	// Unix-like system, so just assume Unix
	return homeUnix()
}

func HomePath() string {
	home, err := Home()
	if err != nil {
		panic("Home path error: " + err.Error())
	}
	return home
}

func homeUnix() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}

type HostInfo struct {
	Os         string `json:"os"`
	HostName   string `json:"hostname"`
	Platform   string `json:"platform"`
	KernelArch string `json:"kernelArch"`
}

type MemInfo struct {
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"usedPercent"`
}

type CpuInfo struct {
	Cores      int     `json:"cores"`
	CpuPercent float64 `json:"cpuPercent"`
}

type DiskInfo struct {
	Devices []deviceInfo `json:"devices"`
}
type deviceInfo struct {
	Device      string  `json:"device"`
	Mountpoint  string  `json:"mountpoint"`
	Fstype      string  `json:"fstype"`
	Opts        string  `json:"opts"`
	Total       uint64  `json:"total"`
	Usage       float64 `json:"usage"`
	InodesUsage float64 `json:"inodesUsage"`
}

func GetDiskInfo() *DiskInfo {
	//info, _ := disk.IOCounters() //所有硬盘的io信息
	diskPart, err := disk.Partitions(false)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(diskPart)
	var devices []deviceInfo

	for _, dp := range diskPart {
		device := deviceInfo{
			Device:     dp.Device,
			Mountpoint: dp.Mountpoint,
			Fstype:     dp.Fstype,
			Opts:       dp.Opts,
		}

		diskUsed, _ := disk.Usage(dp.Mountpoint)
		// fmt.Printf("分区总大小: %d MB \n", diskUsed.Total/1024/1024)
		// fmt.Printf("分区使用率: %.3f %% \n", diskUsed.UsedPercent)
		// fmt.Printf("分区inode使用率: %.3f %% \n", diskUsed.InodesUsedPercent)
		device.Total = diskUsed.Total / 1024 / 1024
		device.Usage = diskUsed.UsedPercent
		device.InodesUsage = diskUsed.InodesUsedPercent
		devices = append(devices, device)
	}
	return &DiskInfo{
		Devices: devices,
	}
}

func GetCpuInfo() *CpuInfo {

	cpuPercent, _ := cpu.Percent(time.Duration(time.Millisecond*100), false)
	cpuNumber, _ := cpu.Counts(false)
	return &CpuInfo{
		Cores:      cpuNumber,
		CpuPercent: cpuPercent[0],
	}
}

func GetHostInfo() *HostInfo {
	info, _ := host.Info()

	return &HostInfo{
		Os:         info.OS,
		HostName:   info.Hostname,
		Platform:   info.Platform,
		KernelArch: info.KernelArch,
	}
}

func GetMemInfo() *MemInfo {
	v, _ := mem.VirtualMemory()

	// almost every return value is a struct
	//fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)
	// convert to JSON. String() is also implemented
	//fmt.Println(v)

	return &MemInfo{
		Total:       v.Total,
		Free:        v.Free,
		UsedPercent: v.UsedPercent,
	}
}
