package service

import (
	"bufio"
	"fmt"
	"github.com/shirou/gopsutil/process"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
	"time"
)

type ClusterContainer struct {
	container map[string]*LevelInstance
	lock      sync.RWMutex
}

func NewClusterContainer() *ClusterContainer {
	return &ClusterContainer{
		container: map[string]*LevelInstance{},
		lock:      sync.RWMutex{},
	}
}

func (receiver *ClusterContainer) StartLevel(cluster, levelName string, bin int, steamcmd, dstServerInstall, ugcDirectory, persistent_storage_root, conf_dir string) {
	//receiver.lock.Lock()
	//defer receiver.lock.Unlock()

	key := cluster + "_" + levelName
	log.Println("正在启动 ", key)
	value, ok := receiver.container[key]
	if !ok {
		value = NewLevelInstance(cluster, levelName, bin, steamcmd, dstServerInstall, ugcDirectory, persistent_storage_root, conf_dir)
		receiver.container[key] = value
	} else {
		if value.dstSeverInstall != dstServerInstall {
			receiver.Remove(cluster, levelName)
			value = NewLevelInstance(cluster, levelName, bin, steamcmd, dstServerInstall, ugcDirectory, persistent_storage_root, conf_dir)
			receiver.container[key] = value
		}
	}
	value.Start()
}

func (receiver *ClusterContainer) StopLevel(cluster, levelName string) {
	key := cluster + "_" + levelName
	value, ok := receiver.container[key]
	log.Println("正在停止世界", cluster, levelName)
	if ok {
		value.Stop()
	}
}

func (receiver *ClusterContainer) Send(cluster, levelName, message string) {
	key := cluster + "_" + levelName
	value, ok := receiver.container[key]
	if ok {
		value.Send(message)
	}
}

func (receiver *ClusterContainer) Status(cluster, levelName string) bool {
	key := cluster + "_" + levelName
	value, ok := receiver.container[key]
	if ok {
		return value.Status()
	}
	return false
}

func (receiver *ClusterContainer) MemUsage(cluster, levelName string) float64 {
	key := cluster + "_" + levelName
	value, ok := receiver.container[key]
	if ok {
		return value.GetProcessMemInfo()
	}
	return 0
}

func (receiver *ClusterContainer) CpuUsage(cluster, levelName string) float64 {
	key := cluster + "_" + levelName
	value, ok := receiver.container[key]
	if ok {
		return value.GetProcessCpuInfo()
	}
	return 0
}

func (receiver *ClusterContainer) Remove(cluster, levelName string) {
	receiver.lock.RLock()
	defer receiver.lock.RUnlock()
	key := cluster + "_" + levelName
	delete(receiver.container, key)
}

type LevelInstance struct {
	running atomic.Bool
	lock    sync.Mutex
	stdin   io.WriteCloser
	stdout  io.ReadCloser
	pid     int

	cluster                 string
	levelName               string
	bin                     int
	steamcmd                string
	dstSeverInstall         string
	ugc_directory           string
	persistent_storage_root string
	conf_dir                string
}

func NewLevelInstance(cluster, levelName string, bin int, steamcmd, dstServerInstall, ugc_directory, persistent_storage_root, conf_dir string) *LevelInstance {
	running := atomic.Bool{}
	running.Store(false)
	game := &LevelInstance{
		lock:                    sync.Mutex{},
		running:                 running,
		cluster:                 cluster,
		levelName:               levelName,
		bin:                     bin,
		steamcmd:                steamcmd,
		dstSeverInstall:         dstServerInstall,
		ugc_directory:           ugc_directory,
		persistent_storage_root: persistent_storage_root,
		conf_dir:                conf_dir,
	}
	return game
}

func (receiver *LevelInstance) Status() bool {
	return receiver.running.Load()
}

func (receiver *LevelInstance) Start() {
	if receiver.running.Load() == true {
		return
	}
	receiver.lock.Lock()
	// 创建输出文件
	tLogsTxt := receiver.cluster + "_" + receiver.levelName + "_log"
	logFile, err := os.Create(tLogsTxt)
	if err != nil {
		log.Println("Error creating log file:", err)
		return
	}
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			receiver.lock.Unlock()
		}
	}(logFile)

	var args []string
	if receiver.bin == 32 {
		args = append(args, "./dontstarve_dedicated_server_nullrenderer.exe")
	} else {
		args = append(args, "./dontstarve_dedicated_server_nullrenderer.exe")
	}
	args = append(args, "-console")
	args = append(args, "-cluster")
	args = append(args, receiver.cluster)
	args = append(args, "-shard")
	args = append(args, receiver.levelName)
	if receiver.persistent_storage_root != "" {
		args = append(args, "-persistent_storage_root")
		args = append(args, receiver.persistent_storage_root)
	}
	if receiver.conf_dir != "" {
		args = append(args, "-conf_dir")
		args = append(args, receiver.conf_dir)
	}
	if receiver.ugc_directory != "" {
		args = append(args, "-ugc_directory")
		args = append(args, receiver.ugc_directory)
	}
	// 创建一个 cmd 对象
	cmd := exec.Command(args[0], args[1:]...)
	if receiver.bin == 32 {
		cmd.Dir = receiver.dstSeverInstall + "\\bin"
	} else {
		cmd.Dir = receiver.dstSeverInstall + "\\bin64"
	}

	receiver.running.Store(true)
	receiver.lock.Unlock()
	// 获取子进程的 stdin、stdout 和 stderr
	receiver.stdin, err = cmd.StdinPipe()
	if err != nil {
		log.Printf("Error getting stdin pipe: %v\n", err)
		return
	}
	receiver.stdout, err = cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error getting stdout pipe: %v\n", err)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("Error getting stderr pipe: %v\n", err)
		return
	}

	// 启动子进程
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting command: %v\n", err)
		return
	}

	receiver.pid = cmd.Process.Pid
	// 开启 goroutine 实时读取子进程的输出并写入文件和控制台
	go func() {
		io.Copy(io.MultiWriter(os.Stdout, logFile), receiver.stdout)
	}()
	go func() {
		io.Copy(io.MultiWriter(os.Stderr, logFile), stderr)
	}()

	// 等待子进程退出
	if err := cmd.Wait(); err != nil {
		log.Printf("Error waiting for command: %v\n", err)
		if receiver.stdin != nil {
			receiver.stdin.Close()
		}
		if receiver.stdout != nil {
			receiver.stdout.Close()
		}
		receiver.running.Store(false)
	}
	log.Println("process exit !!!")
	receiver.running.Store(false)
}

func (receiver *LevelInstance) Stop() {
	receiver.lock.Lock()
	defer receiver.lock.Unlock()

	if receiver.running.Load() == true {
		err := receiver.Send("c_shutdown(true)")
		if err != nil {
			log.Println("stop game error", err)
		} else {
			receiver.running.Store(false)
		}
	}
}

func (receiver *LevelInstance) Send(cmd string) error {
	// 向子进程写入命令
	writer := bufio.NewWriter(receiver.stdin)
	input := cmd + "\n"
	_, err := writer.WriteString(input)
	if err != nil {
		log.Printf("Error writing to stdin: %v\n", err)
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (receiver *LevelInstance) GetProcessMemInfo() float64 {
	p, err := process.NewProcess(int32(receiver.pid))
	if err != nil {
		fmt.Println("Error getting process info:", err)
		return 0
	}
	// 获取内存使用情况
	memInfo, err := p.MemoryInfo()
	if err != nil {
		fmt.Println("Error getting memory info:", err)
		return 0
	}
	memUsage := float64(memInfo.RSS) / (1024 * 1024) // 以 MB 为单位
	return memUsage
}

func (receiver *LevelInstance) GetProcessCpuInfo() float64 {
	p, err := process.NewProcess(int32(receiver.pid))
	if err != nil {
		fmt.Println("Error getting process info:", err)
		return 0
	}
	// 获取 CPU 使用情况
	cpuPercent, err := p.Percent(time.Second)
	if err != nil {
		fmt.Println("Error getting CPU percent:", err)
		return 0
	}
	return cpuPercent
}
