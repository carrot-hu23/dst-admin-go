package autoCheck

import (
	"dst-admin-go/config/database"
	"dst-admin-go/constant/consts"
	"dst-admin-go/constant/dst"
	"dst-admin-go/model"
	"dst-admin-go/utils/clusterUtils"
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
	"sync"
	"time"
)

type AutoCheckManager struct {
	AutoChecks []model.AutoCheck
	statusMap  map[string]chan int
	mutex      sync.Mutex
}

func (m *AutoCheckManager) Start(autoChecks []model.AutoCheck) {

	m.AutoChecks = autoChecks
	m.statusMap = make(map[string]chan int)

	for i := range autoChecks {
		taskId := autoChecks[i].Uuid
		m.statusMap[taskId] = make(chan int)
	}

	for i := range autoChecks {
		go func(index int) {
			defer func() {
				if r := recover(); r != nil {
					log.Println(r)
				}
			}()
			taskId := autoChecks[index].ClusterName + autoChecks[index].LevelName
			m.run(autoChecks[index], m.statusMap[taskId])
		}(i)
	}

}

func (m *AutoCheckManager) run(task model.AutoCheck, stop chan int) {
	for {
		select {
		case <-stop:
			return
		default:
			m.check(task)
		}
	}
}
func (m *AutoCheckManager) GetAutoCheck(clusterName, levelName string) *model.AutoCheck {
	db := database.DB
	autoCheck := model.AutoCheck{}
	db.Where("cluster_name = ? and level_name", clusterName, levelName).Find(&autoCheck)
	if autoCheck.Interval == 0 {
		autoCheck.Interval = 10
	}
	return &autoCheck
}

// TODO 这里要修改
func (m *AutoCheckManager) check(task model.AutoCheck) {
	autoCheck := m.GetAutoCheck(task.ClusterName, task.LevelName)
	if autoCheck.Enable != 1 {
		time.Sleep(10 * time.Second)
	} else {
		checkInterval := time.Duration(autoCheck.Interval) * time.Minute
		strategy := StrategyMap[task.CheckType]
		if !strategy.Check(task.ClusterName, task.Uuid) {
			log.Println(task.ClusterName, task.Uuid, " is not running, waiting for ", checkInterval)
			time.Sleep(checkInterval)
			if !strategy.Check(task.ClusterName, task.Uuid) {
				log.Println(task.ClusterName, task.Uuid, "has not started, starting it...")
				err := strategy.Run(task.ClusterName, task.Uuid)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		time.Sleep(checkInterval)
	}
}

func (m *AutoCheckManager) AddAutoCheckTasks(task model.AutoCheck) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	taskId := task.Uuid
	if oldChan, ok := m.statusMap[taskId]; ok {
		close(oldChan) // 关闭旧的通道
	}
	m.statusMap[taskId] = make(chan int)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		m.run(task, m.statusMap[taskId])
	}()

}

func (m *AutoCheckManager) DeleteAutoCheck(taskId string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if ch, ok := m.statusMap[taskId]; ok {
		close(ch)                   // 关闭通道
		delete(m.statusMap, taskId) // 从 statusMap 中删除键值对
	}

}

var StrategyMap = map[string]CheckStrategy{}

func init() {
	StrategyMap[consts.UPDATE_GAME] = &GameUpdateCheck{}
	StrategyMap[consts.LEVEL_MOD] = &LevelModCheck{}
	StrategyMap[consts.LEVEL_DOWN] = &LevelModCheck{}
}

type CheckStrategy interface {
	Check(string, string) bool
	Run(string, string) error
}

type LevelModCheck struct{}

func (s *LevelModCheck) Check(clusterName, levelName string) bool {
	// 找到当前存档的modId, 然后根据判断当前存档的
	cluster := clusterUtils.GetCluster(clusterName)
	modoverridesPath := dst.GetLevelModoverridesPath(clusterName, levelName)
	content, err := fileUtils.ReadFile(modoverridesPath)
	if err != nil {
		return true
	}
	workshopIds := dstUtils.WorkshopIds(content)
	if len(workshopIds) == 0 {
		return true
	}

	acfPath := filepath.Join(cluster.ForceInstallDir, "ugc_mods", cluster.ClusterName, levelName, "appworkshop_322330.acf")
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
	return diffFetchModInfo2(activeModMap)
}

// Run 更新会重启所有世界
func (s *LevelModCheck) Run(clusterName, levelName string) error {

	SendAnnouncement2(clusterName, levelName)

	gameService.StopLevel(clusterName, levelName)
	time.Sleep(1 * time.Minute)
	gameService.StartGame(clusterName)
	return nil
}

type LevelDownCheck struct{}

func (s *LevelDownCheck) Check(clusterName, levelName string) bool {
	return gameService.GetLevelStatus(clusterName, levelName)
}

func (s *LevelDownCheck) Run(clusterName, levelName string) error {
	cluster := clusterUtils.GetCluster(clusterName)
	bin := cluster.Bin
	beta := cluster.Beta
	gameService.LaunchLevel(clusterName, levelName, bin, beta)
	return nil
}

type GameUpdateCheck struct{}

func (s *GameUpdateCheck) Check(clusterName, levelName string) bool {
	localDstVersion := gameService.GetLocalDstVersion(clusterName)
	lastDstVersion := gameService.GetLastDstVersion()
	return lastDstVersion > localDstVersion
}

func (s *GameUpdateCheck) Run(clusterName, levelName string) error {

	SendAnnouncement2(clusterName, levelName)

	return gameService.UpdateGame(clusterName)
}

func SendAnnouncement2(clusterName string, levelName string) {
	db := database.DB
	autoCheck := model.AutoCheck{}
	db.Where("uuid = ?", levelName).Find(&autoCheck)
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

func diffFetchModInfo2(activeModMap map[string]dstUtils.WorkshopItem) bool {

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
