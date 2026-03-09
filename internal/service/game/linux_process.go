package game

import (
	"dst-admin-go/internal/pkg/utils/dstUtils"
	"dst-admin-go/internal/pkg/utils/shellUtils"
	"dst-admin-go/internal/service/dstConfig"
	"dst-admin-go/internal/service/levelConfig"
	"log"
	"strings"
	"sync"
	"time"
)

type LinuxProcess struct {
	dstConfig        dstConfig.Config
	levelConfigUtils *levelConfig.LevelConfigUtils
	mu               sync.Mutex // 保护启动/停止操作，防止并发执行
}

func NewLinuxProcess(dstConfig dstConfig.Config, levelConfigUtils *levelConfig.LevelConfigUtils) *LinuxProcess {
	return &LinuxProcess{
		dstConfig:        dstConfig,
		levelConfigUtils: levelConfigUtils,
	}
}

func (p *LinuxProcess) SessionName(clusterName, levelName string) string {
	return "DST_8level_" + levelName + "_" + clusterName
}

func (p *LinuxProcess) Start(clusterName, levelName string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	err := p.stop(clusterName, levelName)
	if err != nil {
		return err
	}
	return p.launchLevel(clusterName, levelName)
}

func (p *LinuxProcess) launchLevel(clusterName, levelName string) error {
	cluster, err := p.dstConfig.GetDstConfig(clusterName)
	if err != nil {
		return err
	}
	bin := cluster.Bin
	dstInstallDir := cluster.Force_install_dir
	if cluster.Beta == 1 {
		dstInstallDir = dstInstallDir + "-beta"
	}
	ugcDirectory := cluster.Ugc_directory
	persistent_storage_root := cluster.Persistent_storage_root
	conf_dir := cluster.Conf_dir
	var startCmd = ""

	dstInstallDir = dstUtils.EscapePath(dstInstallDir)
	log.Println(dstInstallDir)
	sessionName := p.SessionName(clusterName, levelName)
	if bin == 64 {
		startCmd = "cd " + dstInstallDir + "/bin64 ; screen -d -m -S \"" + sessionName + "\"  ./dontstarve_dedicated_server_nullrenderer_x64 -console -cluster " + clusterName + " -shard " + levelName
	} else if bin == 100 {
		startCmd = "cd " + dstInstallDir + "/bin64 ; screen -d -m -S \"" + sessionName + "\"  ./dontstarve_dedicated_server_nullrenderer_x64_luajit -console -cluster " + clusterName + " -shard " + levelName
	} else if bin == 86 {
		startCmd = "cd " + dstInstallDir + "/bin64 ; screen -d -m -S \"" + sessionName + "\" box86 ./dontstarve_dedicated_server_nullrenderer_x64 -console -cluster " + clusterName + " -shard " + levelName
	} else if bin == 2664 {
		startCmd = "cd " + dstInstallDir + "/bin64 ; screen -d -m -S \"" + sessionName + "\" box64 ./dontstarve_dedicated_server_nullrenderer_x64 -console -cluster " + clusterName + " -shard " + levelName
	} else {
		startCmd = "cd " + dstInstallDir + "/bin ; screen -d -m -S \"" + sessionName + "\"  ./dontstarve_dedicated_server_nullrenderer -console -cluster " + clusterName + " -shard " + levelName
	}

	if ugcDirectory != "" {
		startCmd += " -ugc_directory " + ugcDirectory
	}
	if persistent_storage_root != "" {
		startCmd += " -persistent_storage_root " + persistent_storage_root
	}
	if conf_dir != "" {
		startCmd += " -conf_dir " + conf_dir
	}
	startCmd += "  ;"
	log.Println("正在启动世界", "cluster: ", clusterName, "level: ", levelName, "command: ", startCmd)
	_, err = shellUtils.Shell(startCmd)
	return err
}

func (p *LinuxProcess) shutdownLevel(clusterName, levelName string) error {
	ok, err := p.Status(clusterName, levelName)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	shell := "screen -S \"" + p.SessionName(clusterName, levelName) + "\" -p 0 -X stuff \"c_shutdown(true)\\n\""
	log.Println("正在shutdown世界", "cluster: ", clusterName, "level: ", levelName, "command: ", shell)
	_, err = shellUtils.Shell(shell)
	return err
}

func (p *LinuxProcess) killLevel(clusterName, level string) error {

	if ok, err := p.Status(clusterName, level); err != nil || !ok {
		return nil
	}
	cmd := " ps -ef | grep -v grep | grep -v tail |grep '" + clusterName + "'|grep " + level + " |sed -n '1P'|awk '{print $2}' |xargs kill -9"
	log.Println("正在kill世界", "cluster: ", clusterName, "level: ", level, "command: ", cmd)
	_, err := shellUtils.Shell(cmd)
	return err
}

func (p *LinuxProcess) Stop(clusterName, levelName string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.stop(clusterName, levelName)
}

// stop 内部实现，不加锁，供 Start 等方法内部调用
func (p *LinuxProcess) stop(clusterName, levelName string) error {
	p.shutdownLevel(clusterName, levelName)
	time.Sleep(3 * time.Second)

	if ok, err := p.Status(clusterName, levelName); err == nil && ok {
		var i uint8 = 1
		for {
			if ok, err := p.Status(clusterName, levelName); err == nil && ok {
				break
			}
			p.shutdownLevel(clusterName, levelName)
			log.Println("正在第", i, "次stop世界", "cluster: ", clusterName, "level: ", levelName)
			time.Sleep(1 * time.Second)
			i++
			if i > 3 {
				break
			}
		}
	}
	log.Println("使用kill命令强制结束世界", "cluster: ", clusterName, "level: ", levelName)
	return p.killLevel(clusterName, levelName)
}

func (p *LinuxProcess) StartAll(clusterName string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	err := p.stopAll(clusterName)
	if err != nil {
		return err
	}
	config, err := p.levelConfigUtils.GetLevelConfig(clusterName)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(len(config.LevelList))
	for i := range config.LevelList {
		go func(i int) {
			defer func() {
				wg.Done()
				if r := recover(); r != nil {
					log.Println(r)
				}
			}()
			levelName := config.LevelList[i].File
			err := p.launchLevel(clusterName, levelName)
			if err != nil {
				log.Println(err)
				return
			}
		}(i)
	}
	ClearScreen()
	wg.Wait()
	return nil
}

func (p *LinuxProcess) StopAll(clusterName string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.stopAll(clusterName)
}

// stopAll 内部实现，不加锁，供 StartAll 等方法内部调用
func (p *LinuxProcess) stopAll(clusterName string) error {
	config, err := p.levelConfigUtils.GetLevelConfig(clusterName)
	if err != nil {
		log.Panicln(err)
	}
	var wg sync.WaitGroup
	wg.Add(len(config.LevelList))
	for i := range config.LevelList {
		go func(i int) {
			defer func() {
				wg.Done()
				if r := recover(); r != nil {
					log.Println(r)
				}
			}()
			levelName := config.LevelList[i].File
			err := p.stop(clusterName, levelName)
			if err != nil {
				return
			}
		}(i)
	}
	wg.Wait()
	return nil
}

func (p *LinuxProcess) Status(clusterName, levelName string) (bool, error) {
	cmd := " ps -ef | grep -v grep | grep -v tail |grep '" + clusterName + "'|grep " + levelName + " |sed -n '1P'|awk '{print $2}' "
	result, err := shellUtils.Shell(cmd)
	if err != nil {
		return false, nil
	}
	res := strings.Split(result, "\n")[0]
	return res != "", nil
}

func (p *LinuxProcess) Command(clusterName, levelName, command string) error {
	cmd := "screen -S \"" + p.SessionName(clusterName, levelName) + "\" -p 0 -X stuff \"" + command + "\\n\""
	_, err := shellUtils.Shell(cmd)
	return err
}

func (p *LinuxProcess) PsAuxSpecified(clusterName, levelName string) DstPsAux {
	dstPsAux := DstPsAux{}
	cmd := "ps -aux | grep -v grep | grep -v tail | grep " + clusterName + "  | grep " + levelName + " | sed -n '2P' |awk '{print $3, $4, $5, $6}'"

	info, err := shellUtils.Shell(cmd)
	if err != nil {
		log.Println(cmd + " error: " + err.Error())
		return dstPsAux
	}
	if info == "" {
		return dstPsAux
	}

	arr := strings.Split(info, " ")
	dstPsAux.CpuUage = strings.Replace(arr[0], "\n", "", -1)
	dstPsAux.MemUage = strings.Replace(arr[1], "\n", "", -1)
	dstPsAux.VSZ = strings.Replace(arr[2], "\n", "", -1)
	dstPsAux.RSS = strings.Replace(arr[3], "\n", "", -1)

	return dstPsAux
}

const (
	// ClearScreenCmd 检查目前所有的screen作业，并删除已经无法使用的screen作业
	ClearScreenCmd = "screen -wipe "
)

func ClearScreen() bool {
	result, err := shellUtils.Shell(ClearScreenCmd)
	if err != nil {
		return false
	}
	res := strings.Split(result, "\n")[0]
	return res != ""
}
