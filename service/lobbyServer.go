package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/glebarez/sqlite"
	lua "github.com/yuin/gopher-lua"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	PlatformsMapper = map[string]int{
		"Steam":  1,
		"Rail":   4,
		"PSN":    2,
		"XBone":  16,
		"Switch": 32,
	}
	Platforms = []string{"Steam", "PSN", "Rail", "XBone", "Switch"}
	Token     = "pds-g^KU_qE7e8rv1^VVrVXd/01kBDicd7UO5LeL+uYZH1+geZlrutzItvOaw="

	LobbyRegionUrl = "https://lobby-v2-cdn.klei.com/regioncapabilities-v2.json"
	LobbyListUrl   = "https://lobby-v2-cdn.klei.com/%s-%s.json.gz"
	LobbyReadUrl   = "https://lobby-v2-%s.klei.com/lobby/read"
)

type LobbyServer struct {
	DB  *gorm.DB
	DB2 *gorm.DB
}

func NewLobbyServerWithDB(db *gorm.DB, enableMemory bool) *LobbyServer {
	lobbyServer := &LobbyServer{}
	lobbyServer.DB = db

	if enableMemory {
		log.Println("正在使用内存模式 保存 lobby home")
		db2, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Error),
		})
		if err != nil {
			log.Println(err)
		}
		db2.AutoMigrate(&LobbyHome{})
		lobbyServer.DB2 = db2
	} else {
		db2, err := gorm.Open(sqlite.Open("lobby-home-db"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Error),
		})
		if err != nil {
			log.Println(err)
		}
		db2.AutoMigrate(&LobbyHome{})
		lobbyServer.DB2 = db2
	}

	return lobbyServer
}

func NewLobbyServer() *LobbyServer {

	lobbyServer := &LobbyServer{}

	db, err := gorm.Open(sqlite.Open("dst-db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		panic("failed to connect database")
	}
	lobbyServer.DB = db
	err = lobbyServer.DB.AutoMigrate(
		&LobbyHome{},
	)
	if err != nil {
		panic(err)
	}

	return lobbyServer
}

type LobbyRegionsBody struct {
	LobbyRegions []struct {
		Region string `json:"Region"`
	} `json:"LobbyRegions"`
}

type LobbyHome struct {
	gorm.Model

	Addr            string `json:"__addr"`
	RowID           string `json:"__rowId"`
	Host            string `json:"host"`
	Clanonly        bool   `json:"clanonly"`
	Platform        int    `json:"platform"`
	Mods            bool   `json:"mods"`
	Name            string `json:"name"`
	Pvp             bool   `json:"pvp"`
	Session         string `json:"session"`
	Fo              bool   `json:"fo"`
	Password        bool   `json:"password"`
	GUID            string `json:"guid"`
	Maxconnections  int    `json:"maxconnections"`
	Dedicated       bool   `json:"dedicated"`
	Clienthosted    bool   `json:"clienthosted"`
	Connected       int    `json:"connected"`
	Mode            string `json:"mode"`
	Port            int    `json:"port"`
	V               int    `json:"v"`
	Tags            string `json:"tags"`
	Season          string `json:"season"`
	Lanonly         bool   `json:"lanonly"`
	Intent          string `json:"intent"`
	Allownewplayers bool   `json:"allownewplayers"`
	Serverpaused    bool   `json:"serverpaused"`
	Steamid         string `json:"steamid"`
	Steamroom       string `json:"steamroom"`
	// 多层世界
	Secondaries     map[string]SecondaryInfo `gorm:"-" json:"secondaries"`
	SecondariesJson string                   `json:"secondariesJson"`

	Region    string `json:"region"`
	Platform2 string `json:"platform2"`
}

type LobbyHomeListBody struct {
	GET []LobbyHome "GET"
}

type SecondaryInfo struct {
	Addr    string `json:"__addr"`
	ID      string `json:"id"`
	SteamID string `json:"steamid"`
	Port    int    `json:"port"`
}

type DayData struct {
	Day                 int `json:"day"`
	Dayselapsedinseason int `json:"dayselapsedinseason"`
	Daysleftinseason    int `json:"daysleftinseason"`
}

type LobbyHomeDetail struct {
	gorm.Model

	Addr            string                   `json:"__addr"`
	RowID           string                   `json:"__rowId"`
	Host            string                   `json:"host"`
	Steamclanid     string                   `json:"steamclanid"`
	Clanonly        bool                     `json:"clanonly"`
	Platform        int                      `json:"platform"`
	Mods            bool                     `json:"mods"`
	Name            string                   `json:"name"`
	Pvp             bool                     `json:"pvp"`
	Session         string                   `json:"session"`
	Fo              bool                     `json:"fo"`
	Password        bool                     `json:"password"`
	GUID            string                   `json:"guid"`
	Maxconnections  int                      `json:"maxconnections"`
	Dedicated       bool                     `json:"dedicated"`
	Clienthosted    bool                     `json:"clienthosted"`
	Connected       int                      `json:"connected"`
	Mode            string                   `json:"mode"`
	Port            int                      `json:"port"`
	V               int                      `json:"v"`
	Tags            string                   `json:"tags"`
	Season          string                   `json:"season"`
	Lanonly         bool                     `json:"lanonly"`
	Intent          string                   `json:"intent"`
	Allownewplayers bool                     `json:"allownewplayers"`
	Serverpaused    bool                     `json:"serverpaused"`
	Steamid         string                   `json:"steamid"`
	Steamroom       string                   `json:"steamroom"`
	Secondaries     map[string]SecondaryInfo `gorm:"-" json:"secondaries"`
	SecondariesJson string                   `json:"secondariesJson"`
	Data            string                   `json:"data"`
	Worldgen        string                   `json:"worldgen"`
	Players         string                   `json:"players"`
	ModsInfo        []interface{}            `gorm:"-" json:"mods_info"`
	Desc            string                   `json:"desc"`
	Tick            int                      `json:"tick"`
	Clientmodsoff   bool                     `json:"clientmodsoff"`
	Nat             int                      `json:"nat"`

	PlayerList []Player `gorm:"-" json:"playerList"`
	DayData    DayData  `gorm:"-" json:"dayData"`
}
type LobbyHomeDetailBody struct {
	GET []LobbyHomeDetail `json:"GET"`
}

func (l *LobbyServer) queryRegions() (LobbyRegionsBody, error) {
	resp, err := http.Get(LobbyRegionUrl)
	if err != nil {
		fmt.Println("GET LobbyRegion request failed:", err)
		return LobbyRegionsBody{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return LobbyRegionsBody{}, err
	}

	var data LobbyRegionsBody
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Failed to unmarshal response body:", err)
		return LobbyRegionsBody{}, err
	}

	return data, err
}

func (l *LobbyServer) QueryPlayer(region, playerName string) []LobbyHomeDetail {
	var result []LobbyHomeDetail
	if lobbyList, err := l.requestLobbyV2Api(region, "Steam"); err == nil {
		fmt.Println("Size of lobbyList:", len(lobbyList.GET))
		var w sync.WaitGroup
		w.Add(len(lobbyList.GET))
		for _, lobbyHome := range lobbyList.GET {
			go func(lobbyHome LobbyHome, region string) {
				if lobbyHomeDetailBody, err := l.requestLobbyHomeDetailApi(region, lobbyHome.RowID); err == nil {
					if len(lobbyHomeDetailBody.GET) > 0 {
						lobbyHomeDetail := lobbyHomeDetailBody.GET[0]

						if playerName != "" && strings.Contains(lobbyHomeDetail.Players, playerName) {
							log.Println(lobbyHomeDetail.Name)

							lobbyHomeDetail.PlayerList = l.getPlayer(lobbyHomeDetail.Players)
							lobbyHomeDetail.DayData = l.getDayData(lobbyHomeDetail.Data)
							if lobbyHomeDetail.Secondaries != nil {
								secondariesJson, _ := json.Marshal(lobbyHomeDetail.Secondaries)
								lobbyHomeDetail.SecondariesJson = string(secondariesJson)
							}

							result = append(result, lobbyHomeDetail)
						}
					}
				}
				w.Done()
			}(lobbyHome, region)
		}
		w.Wait()
	}

	return result
}

func (l *LobbyServer) StartCollect(interval int64) {

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				l.saveLobbyList()
			}
		}
	}()
}

func (l *LobbyServer) saveLobbyList() {
	if regions, err := l.queryRegions(); err == nil {
		log.Println("开始同步大厅数据")
		defer log.Println("结束同步大厅数据")
		var wg sync.WaitGroup
		wg.Add(len(regions.LobbyRegions) * len(Platforms))
		for _, region := range regions.LobbyRegions {
			for _, platform := range Platforms {
				go func(region, platform string) {
					defer wg.Done()
					if lobbyList, err := l.requestLobbyV2Api(region, platform); err == nil {
						fmt.Println("Size of lobbyList:", len(lobbyList.GET))
						for i := range lobbyList.GET {
							if lobbyList.GET[i].Secondaries != nil {
								secondariesJson, _ := json.Marshal(lobbyList.GET[i].Secondaries)
								lobbyList.GET[i].SecondariesJson = string(secondariesJson)
							}
						}
						l.updateLobbyHome(lobbyList.GET, region, platform)
					}
				}(region.Region, platform)
			}
		}
		wg.Wait()
	}

}

func (l *LobbyServer) copyProperty(oldRoom, newRoom *LobbyHome) bool {

	var update bool

	if oldRoom.Host != newRoom.Host || oldRoom.Addr != newRoom.Addr || oldRoom.Clanonly != newRoom.Clanonly || oldRoom.Platform != newRoom.Platform || oldRoom.Mods != newRoom.Mods || oldRoom.Name != newRoom.Name || oldRoom.Pvp != newRoom.Pvp || oldRoom.Fo != newRoom.Fo || oldRoom.Password != newRoom.Password || oldRoom.GUID != newRoom.GUID || oldRoom.Maxconnections != newRoom.Maxconnections || oldRoom.Dedicated != newRoom.Dedicated || oldRoom.Clienthosted != newRoom.Clienthosted || oldRoom.Connected != newRoom.Connected || oldRoom.Mode != newRoom.Mode || oldRoom.Port != newRoom.Port || oldRoom.V != newRoom.V || oldRoom.Tags != newRoom.Tags || oldRoom.Session != newRoom.Session || oldRoom.Lanonly != newRoom.Lanonly || oldRoom.Intent != newRoom.Intent || oldRoom.Allownewplayers != newRoom.Allownewplayers || oldRoom.Serverpaused != newRoom.Serverpaused || oldRoom.Steamid != newRoom.Steamid {
		update = true
	}

	oldRoom.Addr = newRoom.Addr
	oldRoom.Host = newRoom.Host
	oldRoom.Clanonly = newRoom.Clanonly
	oldRoom.Platform = newRoom.Platform
	oldRoom.Mods = newRoom.Mods
	oldRoom.Name = newRoom.Name
	oldRoom.Pvp = newRoom.Pvp
	oldRoom.Session = newRoom.Session
	oldRoom.Fo = newRoom.Fo
	oldRoom.Password = newRoom.Password
	oldRoom.GUID = newRoom.GUID
	oldRoom.Maxconnections = newRoom.Maxconnections
	oldRoom.Dedicated = newRoom.Dedicated
	oldRoom.Clienthosted = newRoom.Clienthosted
	oldRoom.Connected = newRoom.Connected
	oldRoom.Mode = newRoom.Mode
	oldRoom.Port = newRoom.Port
	oldRoom.V = newRoom.V
	oldRoom.Tags = newRoom.Tags
	oldRoom.Session = newRoom.Session
	oldRoom.Lanonly = newRoom.Lanonly
	oldRoom.Intent = newRoom.Intent
	oldRoom.Allownewplayers = newRoom.Allownewplayers
	oldRoom.Serverpaused = newRoom.Serverpaused
	oldRoom.Steamid = newRoom.Steamid

	return update
}

func (l *LobbyServer) updateLobbyHome(lobbyHomeList []LobbyHome, region, platform string) {

	log.Println("开始同步中央房间")
	defer log.Println("结束同步中央房间")

	// 创建一个空的 Map
	lobbyHomeMap := make(map[string]LobbyHome)

	// 遍历切片，将字段值添加到 Map 中
	for i, _ := range lobbyHomeList {
		lobbyHomeMap[lobbyHomeList[i].RowID] = lobbyHomeList[i]
	}

	// 更新大厅数据
	for i, _ := range lobbyHomeList {
		lobbyHomeList[i].Region = region
		lobbyHomeList[i].Platform2 = platform
	}

	// 查询数据库中已经存在的房间信息
	lobbyHomeListLen := len(lobbyHomeList)
	existingRooms := make([]LobbyHome, 0, lobbyHomeListLen)
	if err := l.DB2.Where("region = ? and platform2 = ?", region, platform).Find(&existingRooms).Error; err != nil {
		fmt.Println("Error querying existing rooms:", err)
		return
	}
	var createRooms []LobbyHome
	var updateRooms []LobbyHome
	var deleteRooms []LobbyHome

	// 遍历已经存在的房间信息，将不存在于 JSON 数据中的房间信息标记为已删除
	for i, _ := range existingRooms {
		h, ok := lobbyHomeMap[existingRooms[i].RowID]
		if !ok {
			deleteRooms = append(deleteRooms, existingRooms[i])
		} else {
			// 如果存在 就只更新
			if l.copyProperty(&existingRooms[i], &h) {
				updateRooms = append(updateRooms, existingRooms[i])
			}
		}
	}

	for i := range lobbyHomeList {
		notFind := true
		for j := range existingRooms {
			if existingRooms[j].RowID == lobbyHomeList[i].RowID {
				notFind = false
				break
			}
		}
		if notFind {
			createRooms = append(createRooms, lobbyHomeList[i])
		}
	}

	log.Println("create len: ", len(createRooms), lobbyHomeListLen, region)
	log.Println("save len: ", len(updateRooms), lobbyHomeListLen, region)
	log.Println("delete len: ", len(deleteRooms), lobbyHomeListLen, region)

	var batchSize = 800

	if len(deleteRooms) > 0 {
		l.DB2.Delete(&deleteRooms)
	}

	if len(updateRooms) > 0 {
		var batchSize2 = 800
		var wg1 sync.WaitGroup
		for i := 0; i < len(updateRooms); i += batchSize2 {
			wg1.Add(1)
			go func(i int) {
				defer wg1.Done()

				end := i + batchSize2
				if end > len(updateRooms) {
					end = len(updateRooms)
				}

				// 每个批次处理一部分数据
				batch := updateRooms[i:end]
				if err := l.DB2.Save(&batch).Error; err != nil {
					// 处理错误
					log.Println("update table ", err)
				}
			}(i)
		}
		wg1.Wait()

	}

	var wg sync.WaitGroup
	for i := 0; i < len(createRooms); i += batchSize {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			end := i + batchSize
			if end > len(createRooms) {
				end = len(createRooms)
			}

			// 每个批次处理一部分数据
			batch := createRooms[i:end]
			if err := l.DB2.Create(&batch).Error; err != nil {
				// 处理错误
				log.Println("insert table ", err)
			}
		}(i)
	}
	wg.Wait()

}

func (l *LobbyServer) containsRoom(rooms []LobbyHome, room LobbyHome) bool {
	for _, r := range rooms {
		if r.RowID == room.RowID {
			return true
		}
	}
	return false
}

func (l *LobbyServer) requestLobbyV2Api(region, platform string) (LobbyHomeListBody, error) {
	url := fmt.Sprintf(LobbyListUrl, region, platform)
	log.Println("url: ", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("GET LobbyRegion request failed:", err)
		return LobbyHomeListBody{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return LobbyHomeListBody{}, err
	}

	var data LobbyHomeListBody
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Failed to unmarshal response body:", err)
		return LobbyHomeListBody{}, err
	}

	return data, nil
}

func (l *LobbyServer) requestLobbyHomeDetailApi(region, rowId string) (LobbyHomeDetailBody, error) {
	url := fmt.Sprintf(LobbyReadUrl, region)
	// 准备请求数据
	data1 := map[string]interface{}{
		"__gameId": "DontStarveTogether",
		"__token":  Token,
		"query": map[string]string{
			"__rowId": rowId,
		},
	}
	// 将请求数据转换为 JSON 字节
	jsonData, err := json.Marshal(data1)
	if err != nil {
		fmt.Println("Failed to marshal JSON:", err)
		return LobbyHomeDetailBody{}, err
	}

	// 创建请求的 Body
	body1 := bytes.NewBuffer(jsonData)

	// 发送 POST 请求
	resp, err := http.Post(url, "application/json", body1)
	if err != nil {
		// fmt.Println("Failed to send POST request:", err)
		return LobbyHomeDetailBody{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return LobbyHomeDetailBody{}, err
	}

	var data LobbyHomeDetailBody
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Failed to unmarshal response body:", err)
		return LobbyHomeDetailBody{}, err
	}
	return data, nil
}

func (l *LobbyServer) QueryLobbyHomeInfo(region, rowID string) LobbyHomeDetail {
	if lobbyHomeDetailBody, err := l.requestLobbyHomeDetailApi(region, rowID); err == nil {
		if len(lobbyHomeDetailBody.GET) > 0 {
			lobbyHomeDetail := lobbyHomeDetailBody.GET[0]
			lobbyHomeDetail.PlayerList = l.getPlayer(lobbyHomeDetail.Players)
			lobbyHomeDetail.DayData = l.getDayData(lobbyHomeDetail.Data)

			if lobbyHomeDetail.Secondaries != nil {
				secondariesJson, _ := json.Marshal(lobbyHomeDetail.Secondaries)
				lobbyHomeDetail.SecondariesJson = string(secondariesJson)
			}

			return lobbyHomeDetail
		}
	}
	return LobbyHomeDetail{}
}

type Player struct {
	Colour     string `json:"colour"`
	EventLevel int    `json:"eventLevel"`
	Name       string `json:"name"`
	NetID      string `json:"netID"`
	Prefab     string `json:"prefab"`
}

func (l *LobbyServer) getPlayer(playerLua string) []Player {
	if playerLua == "return {  }" {
		return []Player{}
	}
	L := lua.NewState()
	defer L.Close()

	err := L.DoString(playerLua)
	if err != nil {
		fmt.Println("Lua execution error:", err)
		return []Player{}
	}

	// 提取 Lua 表中的值
	luaTable := L.Get(-1) // 获取栈顶的值
	L.Pop(1)              // 从栈中移除该值

	players, ok := luaTable.(*lua.LTable)
	if !ok {
		fmt.Println("Invalid Lua table")
		return []Player{}
	}

	var playerList []Player
	// 解析并打印玩家信息
	for i := 1; ; i++ {
		player := players.RawGetInt(i)
		if player == lua.LNil {
			break
		}

		p := l.parsePlayer(player)
		//fmt.Printf("Player %d: %+v\n", i, p)
		playerList = append(playerList, p)
	}

	return playerList
}

func (l *LobbyServer) parsePlayer(player lua.LValue) Player {
	playerTable, ok := player.(*lua.LTable)
	if !ok {
		return Player{}
	}

	p := Player{}

	colour := playerTable.RawGetString("colour")
	if str, ok := colour.(lua.LString); ok {
		p.Colour = string(str)
	}

	eventLevel := playerTable.RawGetString("eventlevel")
	if num, ok := eventLevel.(lua.LNumber); ok {
		p.EventLevel = int(num)
	}

	name := playerTable.RawGetString("name")
	if str, ok := name.(lua.LString); ok {
		p.Name = string(str)
	}

	netID := playerTable.RawGetString("netid")
	if str, ok := netID.(lua.LString); ok {
		p.NetID = string(str)
	}

	prefab := playerTable.RawGetString("prefab")
	if str, ok := prefab.(lua.LString); ok {
		p.Prefab = string(str)
	}

	return p
}

func (l *LobbyServer) getDayData(dayLua string) DayData {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	// 创建一个 Lua 虚拟机
	L := lua.NewState()
	defer L.Close()

	// 解析 Lua 脚本
	if err := L.DoString(dayLua); err != nil {
		panic(err)
	}
	// 提取 Lua 表格数据并转换为 Go 结构体
	luaTable := L.Get(-1)
	dayData := DayData{
		Day:                 int(luaTable.(*lua.LTable).RawGetString("day").(lua.LNumber)),
		Dayselapsedinseason: int(luaTable.(*lua.LTable).RawGetString("dayselapsedinseason").(lua.LNumber)),
		Daysleftinseason:    int(luaTable.(*lua.LTable).RawGetString("daysleftinseason").(lua.LNumber)),
	}
	// 将 Lua 表转换为 Go 结构体

	return dayData
}

type LobbyStatistics struct {
	gorm.Model

	AllServerTotalCount int
	AllPlayerTotalCount int

	SteamServerTotalCount  int
	PSNServerTotalCount    int
	RailServerTotalCount   int
	XBoneServerTotalCount  int
	SwitchServerTotalCount int

	SteamPlayerTotalCount  int
	PSNPlayerTotalCount    int
	RailPlayerTotalCount   int
	XBonePlayerTotalCount  int
	SwitchPlayerTotalCount int
}

func (l *LobbyServer) StartStatistics(interval int64) {

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				l.statisticLobbyServer()
			}
		}
	}()
}

func (l *LobbyServer) statisticLobbyHomeDetail(data *sync.Map) {
	var lobbyHomeList []LobbyHome
	data.Range(func(key, value interface{}) bool {
		lobbyHomes := value.([]LobbyHome)
		lobbyHomes = append(lobbyHomeList, lobbyHomes...)
		return true
	})

	var result []LobbyHomeDetail
	var wg sync.WaitGroup
	wg.Add(len(lobbyHomeList))
	for i := range lobbyHomeList {
		go func(i int) {
			lobbyHomeDetailBody, err := l.requestLobbyHomeDetailApi(lobbyHomeList[i].Region, lobbyHomeList[i].RowID)
			if err == nil {
				if len(lobbyHomeDetailBody.GET) > 0 {
					lobbyHomeDetail := lobbyHomeDetailBody.GET[0]
					if lobbyHomeDetail.Secondaries != nil {
						secondariesJson, _ := json.Marshal(lobbyHomeDetail.Secondaries)
						lobbyHomeDetail.SecondariesJson = string(secondariesJson)
					}
					result = append(result, lobbyHomeDetail)
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	//批量保存
	var wg2 sync.WaitGroup
	var batchSize = 800
	for i := 0; i < len(result); i += batchSize {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			end := i + batchSize
			if end > len(result) {
				end = len(result)
			}
			// 每个批次处理一部分数据
			batch := result[i:end]
			if err := l.DB.Create(&batch).Error; err != nil {
				// 处理错误
				log.Println("insert table ", err)
			}
		}(i)
	}
	wg2.Wait()

}

func (l *LobbyServer) statisticLobbyServer() {
	if regions, err := l.queryRegions(); err == nil {

		var mapResult sync.Map
		var wg sync.WaitGroup
		wg.Add(len(regions.LobbyRegions) * len(Platforms))

		for _, region := range regions.LobbyRegions {
			for _, platform := range Platforms {
				go func(region, platform string) {
					defer wg.Done()
					if lobbyList, err := l.requestLobbyV2Api(region, platform); err == nil {
						fmt.Println("Size of lobbyList:", len(lobbyList.GET))
						for i := range lobbyList.GET {
							if lobbyList.GET[i].Secondaries != nil {
								secondariesJson, _ := json.Marshal(lobbyList.GET[i].Secondaries)
								lobbyList.GET[i].SecondariesJson = string(secondariesJson)
							}
						}
						var list []LobbyHome
						value, ok := mapResult.Load(platform)
						if ok {
							list = append(list, value.([]LobbyHome)...)
						}
						list = append(list, lobbyList.GET...)
						mapResult.Store(platform, list)
					}
				}(region.Region, platform)
			}
		}
		wg.Wait()

		statistics := &LobbyStatistics{}
		// 遍历 sync.Map 中的数据
		var allServerTotal int
		var allPlayerTotal int

		var lobbyHomeList []LobbyHome
		mapResult.Range(func(key, value interface{}) bool {
			lobbyHomes := value.([]LobbyHome)

			lobbyHomeList = append(lobbyHomeList, lobbyHomes...)

			playerCount := l.playerCount(lobbyHomes)
			serverCount := len(lobbyHomes)
			allServerTotal += serverCount
			allPlayerTotal += playerCount
			if key == "Steam" {
				statistics.SteamServerTotalCount = serverCount
				statistics.SteamPlayerTotalCount = playerCount
			}
			if key == "Rail" {
				statistics.RailServerTotalCount = serverCount
				statistics.RailPlayerTotalCount = playerCount
			}
			if key == "PSN" {
				statistics.PSNServerTotalCount = serverCount
				statistics.PSNPlayerTotalCount = playerCount
			}
			if key == "XBone" {
				statistics.XBoneServerTotalCount = serverCount
				statistics.XBonePlayerTotalCount = playerCount
			}
			if key == "Switch" {
				statistics.SwitchServerTotalCount = serverCount
				statistics.SwitchPlayerTotalCount = playerCount
			}
			return true
		})

		statistics.AllServerTotalCount = allServerTotal
		statistics.AllPlayerTotalCount = allPlayerTotal

		l.DB.Save(statistics)
		l.statisticLobbyServerBiref(lobbyHomeList)
	}
}

func (l *LobbyServer) playerCount(lobbyHomes []LobbyHome) int {
	var total int
	for i, _ := range lobbyHomes {
		connected := lobbyHomes[i].Connected
		total += connected
	}
	return total
}

type LobbyHomeBrief struct {
	gorm.Model

	Addr           string `json:"__addr"`
	RowID          string `json:"__rowId"`
	Host           string `json:"host"`
	Name           string `json:"name"`
	Session        string `json:"session"`
	Maxconnections int    `json:"maxconnections"`
	Connected      int    `json:"connected"`
	Mode           string `json:"mode"`
	Port           int    `json:"port"`
	Region         string `json:"region"`
}

func (l *LobbyServer) statisticLobbyServerBiref(lobbyHomeList []LobbyHome) {

	log.Println("开始采集房间信息: size ", len(lobbyHomeList))
	defer log.Println("结束采集房间信息: size ", len(lobbyHomeList))

	var lobbyHomeBriefList []LobbyHomeBrief
	for i := range lobbyHomeList {
		lb := LobbyHomeBrief{
			Addr:           lobbyHomeList[i].Addr,
			RowID:          lobbyHomeList[i].RowID,
			Host:           lobbyHomeList[i].Host,
			Name:           lobbyHomeList[i].Name,
			Session:        lobbyHomeList[i].Session,
			Maxconnections: lobbyHomeList[i].Maxconnections,
			Connected:      lobbyHomeList[i].Connected,
			Mode:           lobbyHomeList[i].Mode,
			Port:           lobbyHomeList[i].Port,
			Region:         lobbyHomeList[i].Region,
		}
		lobbyHomeBriefList = append(lobbyHomeBriefList, lb)
	}
	batchSize := 800
	var wg sync.WaitGroup
	for i := 0; i < len(lobbyHomeBriefList); i += batchSize {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			end := i + batchSize
			if end > len(lobbyHomeBriefList) {
				end = len(lobbyHomeBriefList)
			}
			// 每个批次处理一部分数据
			batch := lobbyHomeBriefList[i:end]
			if err := l.DB.Create(&batch).Error; err != nil {
				// 处理错误
				log.Println("插入 lobbyHomeBriefList 失败 ", err)
			}
		}(i)
	}
	wg.Wait()
}
