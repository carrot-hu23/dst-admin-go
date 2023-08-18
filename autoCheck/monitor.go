package autoCheck

import (
	"dst-admin-go/config/database"
	"dst-admin-go/config/global"
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
	t             string
	beta          int
	bin           int
	clusterName   string
	checkFunc     func(clusterName string, int2 int, int3 int) bool
	checkInterval time.Duration
	startFunc     func(clusterName string, int2 int, int3 int) error
	stopCh        chan int
}

func NewMonitor(clusterName string, t string, bin int, beta int, checkFunc func(string2 string, int2 int, int3 int) bool, checkInterval time.Duration, startFunc func(string2 string, int2 int, int3 int) error) *Monitor {
	return &Monitor{
		t:             t,
		beta:          beta,
		bin:           bin,
		clusterName:   clusterName,
		checkFunc:     checkFunc,
		checkInterval: checkInterval,
		startFunc:     startFunc,
		stopCh:        make(chan int, 1),
	}
}

func (m *Monitor) Start() {
	log.Println("Starting monitor cluster: ", m.clusterName, m.t, "bin:", m.bin, "beta:", m.beta)
	for {
		select {
		case <-m.stopCh:
			log.Println("Stopping monitor cluster: ", m.clusterName, m.t)
			return
		default:
			if !m.checkFunc(m.clusterName, m.bin, m.beta) {
				log.Println(m.clusterName, m.t, " is not running, waiting for ", m.checkInterval)
				time.Sleep(m.checkInterval)
				if !m.checkFunc(m.clusterName, m.bin, m.beta) {
					log.Println(m.clusterName, m.t, "has not started, starting it...")
					err := m.startFunc(m.clusterName, m.bin, m.beta)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
			time.Sleep(m.checkInterval)
		}
	}
}

func (m *Monitor) Stop() {
	m.stopCh <- 1
}

func IsMasterRunning(clusterName string, bin, beta int) bool {
	db := database.DB
	autoCheck := model.AutoCheck{}
	db.Where("name = ?", consts.MasterRunning).Find(&autoCheck)
	if autoCheck.Enable == 1 {
		return gameService.GetLevelStatus(clusterName, "Master")
	}
	return true
}

func IsCavesRunning(clusterName string, bin, beta int) bool {
	db := database.DB
	autoCheck := model.AutoCheck{}
	db.Where("name = ?", consts.CavesRunning).Find(&autoCheck)
	if autoCheck.Enable == 1 {
		return gameService.GetLevelStatus(clusterName, "Caves")
	}
	return true
}

func StartMasterProcess(clusterName string, bin, beta int) error {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	gameService.StartGame(clusterName, bin, beta, consts.StartMaster)
	return nil
}

func StartCavesProcess(clusterName string, bin, beta int) error {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	gameService.StartGame(clusterName, bin, beta, consts.StartCaves)
	return nil
}

func IsGameUpdateVersionProcess(clusterName string, bin, beta int) bool {
	db := database.DB
	autoCheck := model.AutoCheck{}
	db.Where("name = ?", consts.UpdateGameVersion).Find(&autoCheck)
	if autoCheck.Enable == 1 {
		// diff dst version
		localVersion := gameService.GetLocalDstVersion(clusterName)
		version := gameService.GetLastDstVersion()
		log.Println("localVersion: ", localVersion, "lastVersion: ", version)
		return localVersion == version
	}
	return true
}

func UpdateGameVersionProcess(clusterName string, bin, beta int) error {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	SendAnnouncement(clusterName, consts.UpdateGameVersion)
	err := gameService.UpdateGame(clusterName)
	if err != nil {
		return err
	}
	gameService.StartGame(clusterName, consts.StartGame, bin, beta)
	return nil
}

func IsMasterModUpdateProcess(clusterName string, bin, beta int) bool {
	db := database.DB
	autoCheck := model.AutoCheck{}
	db.Where("name = ?", consts.UpdateMasterMod).Find(&autoCheck)
	if autoCheck.Enable == 1 {
		return isLevelModUpdateProcess(clusterName, bin, beta, consts.Master)
	}
	return true
}

func UpdateMasterModUpdateProcess(clusterName string, bin, beta int) error {
	SendAnnouncement(clusterName, consts.UpdateMasterMod)
	return updateLevelModUpdateProcess(clusterName, bin, beta, consts.StartMaster)
}

func IsCavesModUpdateProcess(clusterName string, bin, beta int) bool {
	db := database.DB
	autoCheck := model.AutoCheck{}
	db.Where("name = ?", consts.UpdateCavesMod).Find(&autoCheck)
	if autoCheck.Enable == 1 {
		return isLevelModUpdateProcess(clusterName, bin, beta, consts.Caves)
	}
	return true
}

func UpdateCavesModUpdateProcess(clusterName string, bin, beta int) error {
	SendAnnouncement(clusterName, consts.UpdateCavesMod)
	for i := 0; i < 3; i++ {
		gameConsoleService.SentBroadcast(clusterName, global.Config.AutoCheck.ModUpdatePrompt)
		time.Sleep(3 * time.Second)
	}
	return updateLevelModUpdateProcess(clusterName, bin, beta, consts.StartCaves)
}

// TODO 有问题
func isLevelModUpdateProcess(clusterName string, bin, beta int, levelName string) bool {

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
	log.Println("acf workshops: ", acfWorkshops)

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

// TODO 有问题
func updateLevelModUpdateProcess(clusterName string, bin, beta int, startOpt int) error {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	log.Println("更新模组")
	gameService.StartGame(clusterName, bin, beta, startOpt)
	return nil
}

func SendAnnouncement(clusterName string, name string) {
	db := database.DB
	autoCheck := model.AutoCheck{}
	db.Where("name = ?", name).Find(&autoCheck)
	size := autoCheck.Times
	for i := 0; i < size; i++ {
		announcement := autoCheck.Announcement
		if announcement != "" {
			lines := strings.Split(announcement, "\n")
			for j := range lines {
				gameConsoleService.SentBroadcast(clusterName, lines[j])
				time.Sleep(300 * time.Millisecond)
			}
		}
		time.Sleep(time.Duration(autoCheck.Sleep) * time.Second)
	}
}
