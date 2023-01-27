package main

import (
	"bytes"
	"dst-admin-go/vo"
	"fmt"
	"text/template"

	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func main() {

	GetHostInfo()
	GetCpuInfo()
	GetMemInfo()
	GetDiskInfo()

	tmpl, err := template.ParseFiles("cluster.ini")
	CheckErr(err)
	buf := new(bytes.Buffer)
	gameConfigVo := vo.GameConfigVO{
		ClusterIntention:   "这是一个测试",
		ClusterName:        "饥荒萌萌哒",
		ClusterDescription: "一起来玩啊",
		GameMode:           "endless",
		Pvp:                false,
		MaxPlayers:         6,
		Token:              "dadmam,dman,dmna m,dnamd",
		PauseWhenNobody:    false,
	}
	tmpl.Execute(buf, gameConfigVo)
	fmt.Printf("buf.String():\n%v\n", buf.String())
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GetHostInfo() {
	info, _ := host.Info()
	fmt.Println(info)
}

func GetCpuInfo() {

}

func GetMemInfo() {
	v, _ := mem.VirtualMemory()

	// almost every return value is a struct
	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)

	// convert to JSON. String() is also implemented
	//fmt.Println(v)
}

func GetDiskInfo() {
	//info, _ := disk.IOCounters() //所有硬盘的io信息
	diskPart, err := disk.Partitions(false)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(diskPart)
	for _, dp := range diskPart {
		fmt.Println(dp)
		diskUsed, _ := disk.Usage(dp.Mountpoint)
		fmt.Printf("分区总大小: %d MB \n", diskUsed.Total/1024/1024)
		fmt.Printf("分区使用率: %.3f %% \n", diskUsed.UsedPercent)
		fmt.Printf("分区inode使用率: %.3f %% \n", diskUsed.InodesUsedPercent)
	}
}
