package autoCheck

import (
	"dst-admin-go/config/database"
	"dst-admin-go/constant/consts"
	"dst-admin-go/constant/dst"
	"dst-admin-go/model"
	"dst-admin-go/service"
	"dst-admin-go/utils/dstConfigUtils"
	"dst-admin-go/utils/dstUtils"
	"dst-admin-go/utils/fileUtils"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var gameService = service.GameService{}
var gameConsoleService = service.GameConsoleService{}

type Monitor struct {
	name        string
	beta        int
	bin         int
	clusterName string
	levelName   string
	checkFunc   func(clusterName string, levelName string, int2 int, int3 int) bool
	startFunc   func(clusterName string, levelName string, int2 int, int3 int) error
	stopCh      chan int
}

func NewMonitor(clusterName, levelName string, name string, bin int, beta int, checkFunc func(string2 string, string3 string, int2 int, int3 int) bool, startFunc func(string2 string, string3 string, int2 int, int3 int) error) *Monitor {
	return &Monitor{
		name:        name,
		beta:        beta,
		bin:         bin,
		clusterName: clusterName,
		levelName:   levelName,
		checkFunc:   checkFunc,
		startFunc:   startFunc,
		stopCh:      make(chan int, 1),
	}
}

func (m *Monitor) Start() {
	log.Println("Starting monitor cluster: ", m.clusterName, m.name, "bin:", m.bin, "beta:", m.beta)
	for {
		select {
		case <-m.stopCh:
			log.Println("Stopping monitor cluster: ", m.clusterName, m.name)
			return
		default:
			autoCheck := GetAutoCheckByName(m.name)
			if autoCheck.Enable != 1 {
				time.Sleep(10 * time.Second)
			} else {
				checkInterval := time.Duration(autoCheck.Interval) * time.Minute
				if !m.checkFunc(m.clusterName, m.levelName, m.bin, m.beta) {
					log.Println(m.clusterName, m.name, " is not running, waiting for ", checkInterval)
					time.Sleep(checkInterval)
					if !m.checkFunc(m.clusterName, m.levelName, m.bin, m.beta) {
						log.Println(m.clusterName, m.name, "has not started, starting it...")
						err := m.startFunc(m.clusterName, m.levelName, m.bin, m.beta)
						if err != nil {
							log.Fatal(err)
						}
					}
				}
				time.Sleep(checkInterval)
			}

		}
	}
}

func (m *Monitor) Stop() {
	m.stopCh <- 1
}

func IsLevelRunning(clusterName string, levelName string, bin, beta int) bool {
	return gameService.GetLevelStatus(clusterName, levelName)
}

func StartLevelProcess(clusterName string, levelName string, bin, beta int) error {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	gameService.LaunchLevel(clusterName, levelName, bin, beta)
	return nil
}

func IsGameUpdateVersionProcess(clusterName string, levelName string, bin, beta int) bool {
	// diff dst version
	localVersion := gameService.GetLocalDstVersion(clusterName)
	version := gameService.GetLastDstVersion()
	log.Println("localVersion: ", localVersion, "lastVersion: ", version)
	return localVersion == version
}

func UpdateGameVersionProcess(clusterName string, levelName string, bin, beta int) error {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	SendAnnouncement(clusterName, levelName, consts.UpdateGameVersion)
	err := gameService.UpdateGame(clusterName)
	if err != nil {
		return err
	}
	gameService.StartGame(clusterName, consts.StartGame, bin, beta)
	return nil
}

func IsLevelModUpdateProcess(clusterName string, levelName string, bin, beta int) bool {

	// 找到当前存档的modId, 然后根据判断当前存档的
	dstConfig := dstConfigUtils.GetDstConfig()
	cluster := dstConfig.Cluster
	modoverridesPath := dst.GetLevelModoverridesPath(clusterName, levelName)
	content, err := fileUtils.ReadFile(modoverridesPath)
	if err != nil {
		return true
	}
	workshopIds := dstUtils.WorkshopIds(content)
	if len(workshopIds) == 0 {
		return true
	}

	acfPath := filepath.Join(dstConfig.Force_install_dir, "ugc_mods", cluster, levelName, "appworkshop_322330.acf")
	acfWorkshops := dstUtils.ParseACFFile(acfPath)

	log.Println("acf path: ", acfPath)
	// log.Println("acf workshops: ", acfWorkshops)

	activeModMap := make(map[string]dstUtils.WorkshopItem)
	for i := range workshopIds {
		key := workshopIds[i]
		value, ok := acfWorkshops[key]
		if ok {
			activeModMap[key] = value
		}
	}
	return diffFetchModInfo(activeModMap)
}

func UpdateModProcess(clusterName string, levelName string, bin, beta int) error {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	SendAnnouncement(clusterName, levelName, levelName+"Mod")
	gameService.StopLevel(clusterName, levelName)
	gameService.LaunchLevel(clusterName, levelName, bin, beta)
	return nil
}

const (
	steamAPIKey = "73DF9F781D195DFD3D19DED1CB72EEE6"
	appID       = 322330
	language    = 6
)

func diffFetchModInfo(activeModMap map[string]dstUtils.WorkshopItem) bool {

	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	var modIds []string
	for key := range activeModMap {
		modIds = append(modIds, key)
	}

	urlStr := "http://api.steampowered.com/IPublishedFileService/GetDetails/v1/"
	data := url.Values{}
	data.Set("key", steamAPIKey)
	data.Set("language", "6")
	for i := range modIds {
		data.Set("publishedfileids["+strconv.Itoa(i)+"]", modIds[i])
	}
	urlStr = urlStr + "?" + data.Encode()

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return true
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return true
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return true
	}

	dataList, ok := result["response"].(map[string]interface{})["publishedfiledetails"].([]interface{})
	if !ok {
		return true
	}
	for i := range dataList {

		data2 := dataList[i].(map[string]interface{})
		_, find := data2["time_updated"]
		if find {
			timeUpdated := data2["time_updated"].(float64)
			modId := data2["publishedfileid"].(string)
			value, ok := activeModMap[modId]
			if ok {
				if timeUpdated > float64(value.TimeUpdated) {
					return false
				}
			}
		}

	}

	return true
}

func SendAnnouncement(clusterName string, levelName string, name string) {
	db := database.DB
	autoCheck := model.AutoCheck{}
	db.Where("name = ?", name).Find(&autoCheck)
	size := autoCheck.Times
	for i := 0; i < size; i++ {
		announcement := autoCheck.Announcement
		if announcement != "" {
			lines := strings.Split(announcement, "\n")
			for j := range lines {
				gameConsoleService.SentBroadcast2(clusterName, levelName, lines[j])
				time.Sleep(300 * time.Millisecond)
			}
		}
		time.Sleep(time.Duration(autoCheck.Sleep) * time.Second)
	}
}

func GetAutoCheckByName(name string) *model.AutoCheck {
	db := database.DB
	autoCheck := model.AutoCheck{}
	db.Where("name = ?", name).Find(&autoCheck)
	if autoCheck.Interval == 0 {
		autoCheck.Interval = 10
	}
	return &autoCheck
}
