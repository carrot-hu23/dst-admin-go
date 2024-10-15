package dockerClient

import (
	"dst-admin-go/config/database"
	"dst-admin-go/model"
	"fmt"
	"github.com/docker/docker/client"
	"log"
)

var Clients map[string]*client.Client
var queueMap map[string]model.QueueInfo

func GetDockerClient(code string) (*client.Client, bool) {
	v, exist := Clients[code]
	return v, exist
}

func InitZoneDockerClient() {

	db := database.DB
	var queues []model.QueueInfo
	db.Find(&queues)

	Clients = make(map[string]*client.Client)
	queueMap = make(map[string]model.QueueInfo)

	for i := range queues {
		code := queues[i].QueueCode
		host := fmt.Sprintf("tcp://%s:%d", queues[i].Ip, queues[i].Port)
		log.Println("正在初始 docker", i+1, host)
		cli, err := client.NewClientWithOpts(client.WithHost(host), client.WithAPIVersionNegotiation())
		if err != nil {
			log.Panicln(err)
		}
		Clients[code] = cli
		queueMap[code] = queues[i]
	}
}

func AddQueue(queue model.QueueInfo) error {
	code := queue.QueueCode
	host := fmt.Sprintf("tcp://%s:%d", queue.Ip, queue.Port)
	log.Println("正在添加 docker client", host)
	cli, err := client.NewClientWithOpts(client.WithHost(host), client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	Clients[code] = cli
	queueMap[code] = queue
	return err
}

func UpdateQueue(queue model.QueueInfo) error {

	// 删除之前的
	delete(Clients, queue.QueueCode)
	delete(queueMap, queue.QueueCode)

	code := queue.QueueCode
	host := fmt.Sprintf("tcp://%s:%d", queue.Ip, queue.Port)
	log.Println("正在更新 docker client", host)
	cli, err := client.NewClientWithOpts(client.WithHost(host), client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	Clients[code] = cli
	queueMap[code] = queue
	return nil
}

func DeleteQueue(code string) {
	// 删除之前的
	delete(Clients, code)
	delete(queueMap, code)
}

func Queue(code string) (model.QueueInfo, bool) {
	v, exist := queueMap[code]
	return v, exist
}
