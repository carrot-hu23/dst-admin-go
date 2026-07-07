package game

import (
	"dst-admin-go/internal/pkg/utils/dstUtils"
	"dst-admin-go/internal/pkg/utils/shellUtils"
	"dst-admin-go/internal/service/dstConfig"
	"dst-admin-go/internal/service/levelConfig"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type LinuxProcess struct {
	dstConfig        dstConfig.Config
	levelConfigUtils *levelConfig.LevelConfigUtils
	locksMu          sync.Mutex
	levelLocks       map[string]*sync.Mutex // 按 cluster/level 保护启停；不同世界不能互相阻塞
}

func NewLinuxProcess(dstConfig dstConfig.Config, levelConfigUtils *levelConfig.LevelConfigUtils) *LinuxProcess {
	return &LinuxProcess{
		dstConfig:        dstConfig,
		levelConfigUtils: levelConfigUtils,
		levelLocks:       make(map[string]*sync.Mutex),
	}
}

func (p *LinuxProcess) levelLock(clusterName, levelName string) *sync.Mutex {
	key := clusterName + "/" + levelName
	p.locksMu.Lock()
	defer p.locksMu.Unlock()
	lock := p.levelLocks[key]
	if lock == nil {
		lock = &sync.Mutex{}
		p.levelLocks[key] = lock
	}
	return lock
}

func (p *LinuxProcess) SessionName(clusterName, levelName string) string {
	return "DST_8level_" + levelName + "_" + clusterName
}

func (p *LinuxProcess) Start(clusterName, levelName string) error {
	lock := p.levelLock(clusterName, levelName)
	lock.Lock()
	defer lock.Unlock()

	if ok, err := p.Status(clusterName, levelName); err != nil {
		return err
	} else if ok {
		// 单世界“启动”按钮应是幂等操作：已经运行就直接成功，不要先 stop 再 restart。
		// 否则用户连续点击/批量点击时会把已启动世界重新关停，并让 UI 等待一次完整 shutdown。
		return nil
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
	tokenPath := p.clusterTokenPath(clusterName, persistent_storage_root, conf_dir)
	if stat, err := os.Stat(tokenPath); err != nil {
		return fmt.Errorf("cluster token not found at %s: %w", tokenPath, err)
	} else if stat.Size() == 0 {
		return fmt.Errorf("cluster token is empty at %s", tokenPath)
	} else {
		log.Println("cluster token check ok", "cluster:", clusterName, "path:", tokenPath, "bytes:", stat.Size())
	}
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
	launchCmd := p.wrapLaunchCommand(sessionName, startCmd)
	log.Println("正在启动世界", "cluster: ", clusterName, "level: ", levelName, "command: ", launchCmd)
	_, err = shellUtils.Shell(launchCmd)
	return err
}

func (p *LinuxProcess) wrapLaunchCommand(sessionName, startCmd string) string {
	// Docker/container deployments usually do not have systemd as PID 1 and must
	// not depend on systemd-run/systemctl.  If a real systemd runtime is available
	// we use a transient scope to protect game processes from panel service
	// restarts; otherwise we fall back to plain screen, which is the Docker path.
	if os.Getenv("DST_DISABLE_SYSTEMD_SCOPE") != "1" && systemdScopeAvailable() {
		scopeUnit := "dst-screen-" + sessionName
		return "systemd-run --quiet --scope --collect --unit " + shellQuote(scopeUnit) + " bash -lc " + shellQuote(startCmd)
	}
	return startCmd
}

func systemdScopeAvailable() bool {
	if _, err := exec.LookPath("systemd-run"); err != nil {
		return false
	}
	if _, err := os.Stat("/run/systemd/system"); err != nil {
		return false
	}
	return true
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}

func (p *LinuxProcess) clusterTokenPath(clusterName, persistentStorageRoot, confDir string) string {
	if confDir == "" {
		confDir = "DoNotStarveTogether"
	}
	base := persistentStorageRoot
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil || home == "" {
			home = "/root"
		}
		base = filepath.Join(home, ".klei")
	}
	return filepath.Join(base, confDir, clusterName, "cluster_token.txt")
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
	pids, err := p.levelPids(clusterName, level)
	if err != nil {
		return err
	}
	if len(pids) == 0 {
		return nil
	}
	log.Println("正在kill世界", "cluster: ", clusterName, "level: ", level, "pids: ", strings.Join(pids, ","))
	args := append([]string{"-9"}, pids...)
	return exec.Command("kill", args...).Run()
}

func (p *LinuxProcess) Stop(clusterName, levelName string) error {
	go func() {
		lock := p.levelLock(clusterName, levelName)
		lock.Lock()
		defer lock.Unlock()
		if err := p.stop(clusterName, levelName); err != nil {
			log.Println("停止世界失败", "cluster:", clusterName, "level:", levelName, "err:", err)
		}
	}()
	return nil
}

// stop 内部实现，不加锁，供 Start 等方法内部调用
func (p *LinuxProcess) stop(clusterName, levelName string) error {
	if ok, err := p.Status(clusterName, levelName); err != nil {
		return err
	} else if !ok {
		return nil
	}

	p.shutdownLevel(clusterName, levelName)
	time.Sleep(3 * time.Second)

	if ok, err := p.Status(clusterName, levelName); err == nil && ok {
		var i uint8 = 1
		for {
			if ok, err := p.Status(clusterName, levelName); err == nil && !ok {
				return nil
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
	if ok, err := p.Status(clusterName, levelName); err != nil {
		return err
	} else if !ok {
		return nil
	}
	log.Println("使用kill命令强制结束世界", "cluster: ", clusterName, "level: ", levelName)
	return p.killLevel(clusterName, levelName)
}

func (p *LinuxProcess) StartAll(clusterName string) error {
	config, err := p.levelConfigUtils.GetLevelConfig(clusterName)
	if err != nil {
		return err
	}
	for i := range config.LevelList {
		go func(i int) {
			defer func() {
				if r := recover(); r != nil {
					log.Println(r)
				}
			}()
			levelName := config.LevelList[i].File
			lock := p.levelLock(clusterName, levelName)
			lock.Lock()
			defer lock.Unlock()

			// The top-left Start button means "start every configured world".
			// It must be idempotent: do not stop/restart worlds that are already running.
			// Stopped worlds are launched; running worlds are left untouched.
			ok, err := p.Status(clusterName, levelName)
			if err != nil {
				log.Println(err)
				return
			}
			if ok {
				log.Println("世界已运行，跳过启动", "cluster: ", clusterName, "level: ", levelName)
				return
			}
			if err := p.launchLevel(clusterName, levelName); err != nil {
				log.Println(err)
				return
			}
		}(i)
	}
	ClearScreen()
	return nil
}

func (p *LinuxProcess) StopAll(clusterName string) error {
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
			lock := p.levelLock(clusterName, levelName)
			lock.Lock()
			defer lock.Unlock()
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
	pids, err := p.levelPids(clusterName, levelName)
	if err != nil {
		return false, err
	}
	return len(pids) > 0, nil
}

func (p *LinuxProcess) levelPids(clusterName, levelName string) ([]string, error) {
	output, err := exec.Command("ps", "-eo", "pid,comm,args", "--no-headers").Output()
	if err != nil {
		return nil, err
	}
	var pids []string
	var candidates []string
	for _, line := range strings.Split(string(output), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		pid, comm, args := fields[0], fields[1], strings.Join(fields[2:], " ")
		if strings.Contains(args, "dontstarve_dedicated_server") || strings.Contains(args, "DST_8level_") {
			candidates = append(candidates, line)
		}
		if strings.Contains(comm, "dontstarve") && strings.Contains(args, clusterName) && strings.Contains(args, levelName) {
			pids = append(pids, pid)
		}
	}
	if len(pids) == 0 && len(candidates) > 0 {
		msg := "levelPids no match cluster=" + clusterName + " level=" + levelName + " candidates=" + strings.Join(candidates, " || ") + "\n"
		log.Print(msg)
		if f, err := os.OpenFile("/tmp/dst-levelpids.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err == nil {
			_, _ = f.WriteString(msg)
			_ = f.Close()
		}
	}
	return pids, nil
}

func (p *LinuxProcess) Command(clusterName, levelName, command string) error {
	// Do not route remote console commands through `sh -c`: panel commands often
	// contain Lua string quotes (for example c_give("cutgrass", 1)).  Shell
	// parsing strips or reinterprets those quotes before `screen` receives them,
	// causing DST to execute c_give(cutgrass, 1) or print(FOO) instead of the
	// intended Lua.  Pass argv directly so the command text is delivered exactly.
	return exec.Command(
		"screen",
		"-S", p.SessionName(clusterName, levelName),
		"-p", "0",
		"-X", "stuff",
		command+"\n",
	).Run()
}

func (p *LinuxProcess) PsAuxSpecified(clusterName, levelName string) DstPsAux {
	dstPsAux := DstPsAux{}
	output, err := exec.Command("ps", "-eo", "pcpu,pmem,vsz,rss,etime,comm,args", "--no-headers").Output()
	if err != nil {
		log.Println("ps status error: " + err.Error())
		return dstPsAux
	}

	for _, line := range strings.Split(string(output), "\n") {
		arr := strings.Fields(line)
		if len(arr) < 7 {
			continue
		}
		args := strings.Join(arr[6:], " ")
		if !strings.Contains(arr[5], "dontstarve") || !strings.Contains(args, clusterName) || !strings.Contains(args, levelName) {
			continue
		}
		dstPsAux.CpuUage = arr[0]
		dstPsAux.MemUage = arr[1]
		dstPsAux.VSZ = arr[2]
		dstPsAux.RSS = arr[3]
		dstPsAux.Elapsed = arr[4]
		return dstPsAux
	}

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
