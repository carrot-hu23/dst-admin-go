package dockerClient

import (
	"dst-admin-go/config"
	"fmt"
	"github.com/docker/docker/client"
	"log"
)

var Clients map[string]*client.Client

var zones []config.Zone

func GetZoneDockerClient(zoneCode string) (*client.Client, bool) {
	v, exist := Clients[zoneCode]
	return v, exist
}

func InitZoneDockerClient(z []config.Zone) {
	Clients = make(map[string]*client.Client)
	zones = z
	if len(zones) == 0 {
		log.Panicln("请配置远程docker client 配置")
	} else {
		for i := range zones {
			zoneCode := zones[i].ZoneCode
			host := fmt.Sprintf("tcp://%s:%d", zones[i].Ip, zones[i].Port)
			log.Println("正在初始 docker", i+1, host)
			cli, err := client.NewClientWithOpts(client.WithHost(host), client.WithAPIVersionNegotiation())
			if err != nil {
				log.Panicln(err)
			}
			Clients[zoneCode] = cli
		}
	}

}
