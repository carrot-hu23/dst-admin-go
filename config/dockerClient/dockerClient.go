package dockerClient

import (
	"dst-admin-go/config/database"
	"dst-admin-go/model"
	"fmt"
	"github.com/docker/docker/client"
	"log"
)

var Clients map[string]*client.Client
var ZoneMap map[string]model.ZoneInfo

func GetZoneDockerClient(zoneCode string) (*client.Client, bool) {
	v, exist := Clients[zoneCode]
	return v, exist
}

func InitZoneDockerClient() {

	db := database.DB
	var zones []model.ZoneInfo
	db.Find(&zones)

	Clients = make(map[string]*client.Client)
	ZoneMap = make(map[string]model.ZoneInfo)

	for i := range zones {
		zoneCode := zones[i].ZoneCode
		host := fmt.Sprintf("tcp://%s:%d", zones[i].Ip, zones[i].Port)
		log.Println("正在初始 docker", i+1, host)
		cli, err := client.NewClientWithOpts(client.WithHost(host), client.WithAPIVersionNegotiation())
		if err != nil {
			log.Panicln(err)
		}
		Clients[zoneCode] = cli
		ZoneMap[zoneCode] = zones[i]
	}
}

func AddZone(zone model.ZoneInfo) {
	zoneCode := zone.ZoneCode
	host := fmt.Sprintf("tcp://%s:%d", zone.Ip, zone.Port)
	log.Println("正在添加 docker client", host)
	cli, err := client.NewClientWithOpts(client.WithHost(host), client.WithAPIVersionNegotiation())
	if err != nil {
		log.Panicln(err)
	}
	Clients[zoneCode] = cli
	ZoneMap[zoneCode] = zone
}

func Zone(zone string) (model.ZoneInfo, bool) {
	v, exist := ZoneMap[zone]
	return v, exist
}
