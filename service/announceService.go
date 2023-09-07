package service

import (
	"dst-admin-go/config/database"
	"dst-admin-go/model"
	"dst-admin-go/utils/dstConfigUtils"
	"log"
	"math/rand"
	"strings"
	"time"
)

type AnnounceService struct {
}

var gameConsoleService GameConsoleService

func shuffleStrings(arr []string) {
	rand.Seed(time.Now().UnixNano())
	for i := len(arr) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
}

func InitAnnounce() {
	go func() {
		for {
			config := dstConfigUtils.GetDstConfig()
			clusterName := config.Cluster

			db := database.DB
			var announce model.Announce
			db.First(&announce)
			if !announce.Enable {
				time.Sleep(1 * time.Minute)
				continue
			}
			content := announce.Content
			if content == "" {
				time.Sleep(5 * time.Second)
				continue
			}
			lines := strings.Split(content, "\n")
			if announce.Method != "order" {
				shuffleStrings(lines)
			}
			log.Println(lines)
			for _, line := range lines {
				line = strings.Replace(line, "\n", "", -1)
				line = strings.Replace(line, "\r", "", -1)

				gameConsoleService.SentBroadcast(clusterName, line)
				time.Sleep(200 * time.Millisecond)
			}
			d := time.Duration(announce.Interval)
			if announce.IntervalUnit == "S" {
				time.Sleep(d * time.Second)
			} else if announce.IntervalUnit == "M" {
				time.Sleep(d * time.Minute)
			} else if announce.IntervalUnit == "H" {
				time.Sleep(d * time.Hour)
			} else {
				time.Sleep(d * time.Second)
			}
		}
	}()
}

func (a *AnnounceService) GetAnnounceSetting() model.Announce {
	db := database.DB
	var announce model.Announce
	db.First(&announce)
	return announce
}

func (a *AnnounceService) SaveAnnounceSetting(announce *model.Announce) {
	db := database.DB
	db.Save(announce)
}
