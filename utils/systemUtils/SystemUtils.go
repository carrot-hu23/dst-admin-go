package systemUtils

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
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
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

type CpuInfo struct {
	Cores          int       `json:"cores"`
	CpuPercent     []float64 `json:"cpuPercent"`
	CpuUsedPercent float64   `json:"cpuUsedPercent"`
	CpuUsed        float64   `json:"cpuUsed"`
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
	cpuInfo := &CpuInfo{}
	cpuPercent, _ := cpu.Percent(0, false)
	cpuInfo.Cores, _ = cpu.Counts(true)
	if len(cpuPercent) == 1 {
		cpuInfo.CpuUsedPercent = cpuPercent[0]
		cpuInfo.CpuUsed = cpuInfo.CpuUsedPercent * 0.01 * float64(cpuInfo.Cores)
	}
	cpuInfo.CpuPercent, _ = cpu.Percent(0, true)
	return cpuInfo
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
		Available:   v.Available,
		Used:        v.Used,
		UsedPercent: v.UsedPercent,
	}
}
