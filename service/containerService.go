package service

import (
	"context"
	"dst-admin-go/config/dockerClient"
	"dst-admin-go/model"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"log"
	"math/rand"
	"strings"
	"time"
)

type ContainerService struct {
}

type Container struct {
	ID          uint   `gorm:"primary_key" json:"id"`
	ContainerId string `json:"containerId"`
	Core        int    `json:"core"`
	Memory      int    `json:"memory"`
	Disk        int    `json:"disk"`
	Image       string `json:"image"`

	Username string `json:"username"`
	Password string `json:"password"`

	Port       int `json:"port"`
	LevelNum   int `json:"levelNum"`
	MaxPlayers int `json:"maxPlayers"`
	MasterPort int `json:"masterPort"`
}

func generateContainerName(prefix string) string {
	rand.Seed(time.Now().UnixNano())
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return fmt.Sprintf("%s-%s", prefix, string(b))
}

func (t *ContainerService) CreateContainer(c model.Cluster) (string, error) {
	// 创建 Docker 客户端
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", fmt.Errorf("无法创建 Docker 客户端: %v", err)
	}

	// 设置容器的环境变量
	env := []string{
		fmt.Sprintf("%s%d", "levelNum=", c.LevelNum),
		fmt.Sprintf("%s%d", "maxPlayers=", c.MaxPlayers),
		fmt.Sprintf("%s%d", "masterPort=", c.MasterPort),
		fmt.Sprintf("%s%d", "maxBackup=", c.MaxBackup),
	}

	// 设置容器卷挂载
	mounts := []string{
		"/root/dst-dedicated-server:/app/dst-dedicated-server",
		"/root/steamcmd:/app/steamcmd",
	}

	// 配置容器资源限制
	config := &container.Config{
		Image: c.Image, // 从数据库中读取镜像名
		Env:   env,
		//Volumes: map[string]struct{}{ // 设置卷
		//	"/app/dst-dedicated-server": {},
		//	"/app/steamcmd":             {},
		//},
		ExposedPorts: nat.PortSet{
			"8082/tcp": struct{}{}, // 暴露容器的 8082 端口
		},
	}

	portBindings := nat.PortMap{}

	portBindings["8082/tcp"] = []nat.PortBinding{
		{
			// 将容器的 8082 端口映射到主机的 8084 端口
			HostPort: fmt.Sprintf("%d", c.Port),
		},
	}

	// 暴露 udp 端口
	for i := 0; i < c.LevelNum; i++ {
		key := fmt.Sprintf("%d/%s", c.MasterPort+1, "udp")
		portBindings[nat.Port(key)] = []nat.PortBinding{
			{
				HostPort: fmt.Sprintf("%d", c.MasterPort+1),
			},
		}
	}

	//// 定义端口映射，将主机的8084端口映射到容器的8082端口
	//portBindings := nat.PortMap{
	//	"8082/tcp": []nat.PortBinding{
	//		{
	//			// 将容器的 8082 端口映射到主机的 8084 端口
	//			HostPort: fmt.Sprintf("%d", c.Port),
	//		},
	//	},
	//}

	// 设置容器资源限制
	hostConfig := &container.HostConfig{
		Resources: container.Resources{
			NanoCPUs: int64(c.Core * 1000000000),           // 核心数，按比例分配 CPU
			Memory:   int64(c.Memory) * 1024 * 1024 * 1024, // 内存，转为字节数
		},
		PortBindings: portBindings,
		Binds:        mounts,
		//StorageOpt: map[string]string{
		//	"size": fmt.Sprintf("%d%s", c.Disk, "g"),
		//},
	}

	// 配置容器网络
	networkConfig := &network.NetworkingConfig{}

	// 暴露的端口配置（根据需要修改）
	portSet := nat.PortSet{}

	config.ExposedPorts = portSet

	// 创建容器
	resp, err := cli.ContainerCreate(context.Background(), config, hostConfig, networkConfig, nil, generateContainerName("dst"))
	if err != nil {
		return "", fmt.Errorf("创建容器失败: %v", err)
	}

	// 启动容器
	if err := cli.ContainerStart(context.Background(), resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("启动容器失败: %v", err)
	}

	return resp.ID, nil
}

func (t *ContainerService) DeleteContainer(containerID string) error {
	// 创建 Docker 客户端
	cli := dockerClient.Client
	log.Println("正在停止容器", containerID)
	// 删除容器
	err := cli.ContainerStop(context.Background(), containerID, container.StopOptions{})
	if err != nil {
		return err
	}
	err = cli.ContainerRemove(context.Background(), containerID, container.RemoveOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (t *ContainerService) RestartContainer(containerID string) error {
	// 创建 Docker 客户端
	cli := dockerClient.Client
	log.Println("正在停止容器", containerID)
	// 重启容器
	err := cli.ContainerRestart(context.Background(), containerID, container.StopOptions{})
	return err
}

func (t *ContainerService) ContainerDstInstallStatus(containerID string) bool {

	cli := dockerClient.Client

	// 要检查的文件路径
	filePath := "/app/dst-dedicated-server/bin/dontstarve_dedicated_server_nullrenderer"

	// 执行命令来检查文件是否存在
	execConfig := types.ExecConfig{
		Cmd:          []string{"sh", "-c", "test -f " + filePath},
		AttachStdout: true,
		AttachStderr: true,
	}

	resp, err := cli.ContainerExecCreate(context.Background(), containerID, execConfig)
	if err != nil {
		panic(err)
	}

	execResp, err := cli.ContainerExecAttach(context.Background(), resp.ID, types.ExecStartCheck{})
	if err != nil {
		panic(err)
	}
	defer execResp.Close()

	var output strings.Builder
	_, err = io.Copy(&output, execResp.Reader)
	if err != nil {
		panic(err)
	}

	if strings.TrimSpace(output.String()) == "" {
		fmt.Printf("文件 %s 不存在\n", filePath)
		return false
	} else {
		fmt.Printf("文件 %s 存在，服务已安装好！\n", filePath)
		return true
	}
}

func (t *ContainerService) ContainerStatusInfo(containerID string) (types.ContainerJSON, error) {

	cli := dockerClient.Client
	containerInfo, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		fmt.Printf("Error inspecting container: %v\n", err)
	}
	return containerInfo, err
}
