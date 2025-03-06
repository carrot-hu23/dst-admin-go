package api

import (
	"bufio"
	"dst-admin-go/config/global"
	"dst-admin-go/constant/consts"
	"dst-admin-go/service/autoCheck"
	"dst-admin-go/utils/dstConfigUtils"
	"dst-admin-go/utils/dstUtils"
	"dst-admin-go/utils/shellUtils"
	"dst-admin-go/utils/systemUtils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type InstallSteamCmd struct{}

// 安装饥荒环境
func (i *InstallSteamCmd) InstallSteamCmd(ctx *gin.Context) {

	//if !atomic.CompareAndSwapInt32(&flag, 0, 1) {
	//	// 已经处理过请求，直接返回结果
	//	ctx.JSON(200, gin.H{"message": "already handled"})
	//	return
	//}
	//atomic.StoreInt32(&flag, 1)
	//defer atomic.StoreInt32(&flag, 0)

	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Access-Control-Allow-Origin", "*")

	// 使用一个channel来接收SSE事件
	eventCh := make(chan string)
	stopCh := make(chan string)

	defer func() {
		if err := recover(); err != nil {
			log.Println("安装依赖错误:", err)
			fmt.Fprintf(ctx.Writer, "data: 安装依赖错误 \n\n")
		}
		close(eventCh)
		close(stopCh)
	}()

	// 在单独的goroutine中发送SSE事件
	go func() {
		i.handle(eventCh, stopCh)
	}()

	// 循环读取channel中的事件并发送给客户端
	for {
		select {
		case event := <-eventCh:
			_, err := fmt.Fprintf(ctx.Writer, event)
			if err != nil {
				// 处理错误情况，例如日志记录或返回错误响应
				fmt.Println("Error writing SSE event:", err)
				return
			}
			if event == "data: end\n\n" {
				log.Println("结束安装")
				return
			}
			ctx.Writer.Flush()
		case <-stopCh:
			return
		case <-ctx.Writer.CloseNotify():
			// 如果客户端断开连接，则停止发送事件
			return
		}
	}

}

func installDependence(eventCh chan string, stopCh chan string) error {
	eventCh <- "data: 正在检测当前系统。。。\n\n"

	info := systemUtils.GetHostInfo()
	eventCh <- "data: " + "OS:" + info.Os + " hostname:" + info.HostName + " platform: " + info.Platform + " kernelArch: " + info.KernelArch + "\n\n"

	// TODO 容器安装依赖环境 不保证成功
	if os.Getenv("DOCKER_CONTAINER") != "" {
		eventCh <- "data: Running in a Docker container \n\n"
	} else if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		eventCh <- "data: Running in a Kubernetes cluster \n\n"
	}

	if strings.Contains(strings.ToLower(info.Platform), "centos") {

		eventCh <- "data: 正在 sudo dpkg --add-architecture i386 \n\n"
		err := command(eventCh, "sudo dpkg --add-architecture i386", "")
		if err != nil {
			eventCh <- "安装失败 \n\n"
		}

		eventCh <- "data: 正在安装 yum update \n\n"
		err = command(eventCh, "yum update -y", "")
		if err != nil {
			eventCh <- "安装失败 \n\n"
		}
		// yum install glibc.i686
		eventCh <- "data: yum install glibc.i686 \n\n"
		err = command(eventCh, "yum install -y glibc.i686", "")
		if err != nil {
			eventCh <- "安装失败 \n\n"
		}
		//yum install libgcc_s.so.1
		eventCh <- "data: yum install libgcc_s.so.1 \n\n"
		err = command(eventCh, "yum install -y libgcc_s.so.1", "")
		if err != nil {
			eventCh <- "安装失败 \n\n"
		}
		eventCh <- "data: 正在安装 glibc.i686 libstdc++.i686 ncurses-libs.i686 screen libcurl.i686 依赖 \n\n"
		err = command(eventCh, "yum install -y lib32gcc1 libcurl4-gnutls-dev:i386 glibc screen wget", "")
		if err != nil {
			eventCh <- "安装失败 \n\n"
		}
		eventCh <- "data: 正在安装 yum install -y SDL2.x86_64 SDL2_gfx-devel.x86_64 SDL2_image-devel.x86_64 SDL2_ttf-devel.x86_64 \n\n"
		err = command(eventCh, "yum install -y SDL2.x86_64 SDL2_gfx-devel.x86_64 SDL2_image-devel.x86_64 SDL2_ttf-devel.x86_64", "")
		if err != nil {
			eventCh <- "安装失败 \n\n"
		}

		eventCh <- "data: 正在建立libcurl-gnutls.so.4软连接 \n\n"
		err = command(eventCh, "ln -s /usr/lib/libcurl.so.4 /usr/lib/libcurl-gnutls.so.4", "")
		if err != nil {
			eventCh <- "建立libcurl-gnutls.so.4软连接失败 \n\n"
		}

	} else if strings.Contains(strings.ToLower(info.Platform), "ubuntu") {

		eventCh <- "data: 正在 sudo dpkg --add-architecture i386 \n\n"
		err := command(eventCh, "sudo dpkg --add-architecture i386", "")
		if err != nil {
			eventCh <- "安装失败 sudo dpkg --add-architecture i386 \n\n"
		}

		// apt-get install -y sudo
		eventCh <- "data: 正在 apt-get install -y sudo \n\n"
		err = command(eventCh, "apt-get install -y sudo", "")
		if err != nil {
			eventCh <- "安装失败 apt-get install -y sudo \n\n"
		}

		eventCh <- "data: 正在 apt-get update \n\n"
		err = command(eventCh, "sudo apt-get update", "")
		if err != nil {
			eventCh <- "安装失败 apt-get update \n\n"
		}

		err = command(eventCh, "sudo apt-get install -y lib32gcc1 libcurl4-gnutls-dev:i386 glibc screen wget sudo", "")
		if err != nil {
			eventCh <- "安装失败 lib32gcc1 libcurl4-gnutls-dev:i386 glibc screen wget sudo \n\n"
		}
		err = command(eventCh, "sudo apt-get install -y libsdl-image1.2-dev libsdl-mixer1.2-dev libsdl-ttf2.0-dev libsdl-gfx1.2-dev", "")
		if err != nil {
			eventCh <- "安装失败 libsdl-image1.2-dev libsdl-mixer1.2-dev libsdl-ttf2.0-dev libsdl-gfx1.2-dev \n\n"
		}

	} else {
		eventCh <- "data: 暂不支持 " + info.Platform + " 请手动安装依赖 \n\n"
		return errors.New("not support yet")
	}

	return nil
}

// 检测是否已经安装了 steamcmd
func installCmd(eventCh chan string, stopCh chan string) error {

	eventCh <- "data: 正在安装steamcmd。。。\n\n"

	// 直接调用脚本安装
	scriptPath := "./static/script/install_steamcmd.sh"
	shellUtils.Chmod(scriptPath)
	err := commandShell(eventCh, scriptPath, consts.HomePath, consts.HomePath)
	if err != nil {
		eventCh <- "data: 安装steamcmd失败！！！ \n\n"
		return err
	}

	return nil
}

func (i *InstallSteamCmd) handle(eventCh chan string, stopCh chan string) {

	err := installDependence(eventCh, stopCh)
	if err != nil {
		return
	}

	err = installCmd(eventCh, stopCh)
	if err != nil {
		return
	}
	saveDstConfig()
	eventCh <- "data: [successed]\n\n"
	eventCh <- "data: end\n\n"
}

func saveDstConfig() {
	// 写入到配置文件里面
	config := dstConfigUtils.GetDstConfig()
	config.Steamcmd = filepath.Join(consts.HomePath, "steamcmd")
	config.Force_install_dir = filepath.Join(consts.HomePath, "dst-dedicated-server")
	config.Backup = dstConfigUtils.KleiDstPath()
	config.Mod_download_path = dstConfigUtils.KleiDstPath()
	config.Cluster = "MyDediServer"

	dstConfigUtils.SaveDstConfig(&config)
	global.Collect.ReCollect(filepath.Join(dstUtils.GetKleiDstPath(), config.Cluster), config.Cluster)
	autoCheck.Manager.ReStart(config.Cluster)
	// autoCheck.AutoCheckObject.RestartAutoCheck(config.Cluster, config.Bin, config.Beta)

	initEvnService.InitBaseLevel(&config, "默认初始", "", true)
}

func commandShell(eventCh chan string, name string, arg ...string) error {
	cmd := exec.Command(name, arg...)

	// 创建管道来获取命令的输出
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating StdoutPipe:", err)
		return err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("Error creating StderrPipe:", err)
		return err
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return err
	}

	// 创建字符串通道来接收输出
	stdoutCh := make(chan string, 1)

	errputCh := make(chan string, 1)

	// 创建协程读取并处理标准输出
	go readAndSend(stdoutPipe, stdoutCh)
	// 创建协程读取并处理标准错误输出
	go readAndSend(stderrPipe, errputCh)

	// 从字符串通道接收输出并处理
	for output := range stdoutCh {
		fmt.Println(output)
		eventCh <- "data: " + output + "\n\n"
	}

	for errput := range errputCh {
		fmt.Println(errput)
		eventCh <- "data: " + errput + "\n\n"
	}

	// 等待命令执行完成
	if err := cmd.Wait(); err != nil {
		fmt.Println("Command finished with error:", err)
	}

	return nil
}

func command(eventCh chan string, name string, arg ...string) error {
	cmd := exec.Command("sh", "-c", name)

	// 创建管道来获取命令的输出
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating StdoutPipe:", err)
		return nil
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("Error creating StderrPipe:", err)
		return nil
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return nil
	}

	// 创建字符串通道来接收输出
	stdoutCh := make(chan string)

	errputCh := make(chan string)

	// 创建协程读取并处理标准输出
	go readAndSend(stdoutPipe, stdoutCh)
	// 创建协程读取并处理标准错误输出
	go readAndSend(stderrPipe, errputCh)

	// 从字符串通道接收输出并处理
	for output := range stdoutCh {
		fmt.Println(output)
		eventCh <- "data: " + output + "\n\n"
	}

	for errput := range errputCh {
		fmt.Println(errput)
		eventCh <- "data: " + errput + "\n\n"
	}

	// 等待命令执行完成
	if err := cmd.Wait(); err != nil {
		fmt.Println("Command finished with error:", err)
	}

	return nil
}

// 读取io.Reader并将每行内容发送到字符串通道
func readAndSend(reader io.Reader, ch chan<- string) {
	bufReader := bufio.NewReader(reader)
	for {
		line, err := bufReader.ReadString('\n')
		if line != "" {
			ch <- line
		}
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from reader:", err)
			}
			break
		}
	}
	// 关闭通道
	defer close(ch)
}
