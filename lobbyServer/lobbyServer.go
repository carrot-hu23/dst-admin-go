package lobbyServer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/glebarez/sqlite"
	lua "github.com/yuin/gopher-lua"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

var (
	Platforms = []string{"Steam", "PSN", "Rail", "XBone", "Switch"}
	Token     = "pds-g^KU_qE7e8rv1^VVrVXd/01kBDicd7UO5LeL+uYZH1+geZlrutzItvOaw="

	LobbyRegionUrl = "https://lobby-v2-cdn.klei.com/regioncapabilities-v2.json"
	LobbyListUrl   = "https://lobby-v2-cdn.klei.com/%s-Steam.json.gz"
	LobbyReadUrl   = "https://lobby-v2-%s.klei.com/lobby/read"
)

type LobbyServer struct {
	DB *gorm.DB
}

func NewLobbyServer2(db *gorm.DB) *LobbyServer {
	lobbyServer := &LobbyServer{}
	lobbyServer.DB = db
	return lobbyServer
}

func NewLobbyServer() *LobbyServer {

	lobbyServer := &LobbyServer{}

	db, err := gorm.Open(sqlite.Open("dst-db"), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
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

	Region string `json:"region"`
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
	Secondaries     map[string]SecondaryInfo `json:"secondaries"`
	SecondariesJson string                   `json:"secondariesJson"`
	Data            string                   `json:"data"`
	Worldgen        string                   `json:"worldgen"`
	Players         string                   `json:"players"`
	ModsInfo        []interface{}            `json:"mods_info"`
	Desc            string                   `json:"desc"`
	Tick            int                      `json:"tick"`
	Clientmodsoff   bool                     `json:"clientmodsoff"`
	Nat             int                      `json:"nat"`

	PlayerList []Player `json:"playerList"`
	DayData    DayData  `json:"dayData"`
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

func (l *LobbyServer) QueryPlayer(playerName string) {
	if regions, err := l.queryRegions(); err == nil {
		var wg sync.WaitGroup
		wg.Add(len(regions.LobbyRegions))
		for _, region := range regions.LobbyRegions {
			go func(region string) {
				defer wg.Done()
				if lobbyList, err := l.requestLobbyV2Api(region); err == nil {
					fmt.Println("Size of lobbyList:", len(lobbyList.GET))
					var w sync.WaitGroup
					w.Add(len(lobbyList.GET))
					for _, lobbyHome := range lobbyList.GET {
						go func(lobbyHome LobbyHome, region string) {
							homeDetail := l.QueryLobbyHomeInfo(region, lobbyHome.RowID)
							if playerName != "" && strings.Contains(homeDetail.Players, playerName) {
								log.Println(homeDetail.Name)
								l.getPlayer(homeDetail.Players)

								log.Println(homeDetail.Secondaries)
							}
							w.Done()
						}(lobbyHome, region)
					}
					w.Wait()
				}
			}(region.Region)
		}
		wg.Wait()
	}
}

func (l *LobbyServer) SaveLobbyList() {
	if regions, err := l.queryRegions(); err == nil {
		var wg sync.WaitGroup
		wg.Add(len(regions.LobbyRegions))
		for _, region := range regions.LobbyRegions {
			go func(region string) {
				defer wg.Done()

				if lobbyList, err := l.requestLobbyV2Api(region); err == nil {
					fmt.Println("Size of lobbyList:", len(lobbyList.GET))
					for i := range lobbyList.GET {
						if lobbyList.GET[i].Secondaries != nil {
							secondariesJson, _ := json.Marshal(lobbyList.GET[i].Secondaries)
							lobbyList.GET[i].SecondariesJson = string(secondariesJson)
						}
					}
					l.updateLobbyHome(lobbyList.GET, region)
				}

			}(region.Region)
		}
		wg.Wait()
	}

}

func (l *LobbyServer) updateLobbyHome(lobbyHomeList []LobbyHome, region string) {

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
	}

	// 查询数据库中已经存在的房间信息
	lobbyHomeListLen := len(lobbyHomeList)
	existingRooms := make([]LobbyHome, 0, lobbyHomeListLen)
	if err := l.DB.Where("region = ?", region).Find(&existingRooms).Error; err != nil {
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
			existingRooms[i].Addr = h.Addr
			existingRooms[i].Host = h.Host
			existingRooms[i].Clanonly = h.Clanonly
			existingRooms[i].Platform = h.Platform
			existingRooms[i].Mods = h.Mods
			existingRooms[i].Name = h.Name
			existingRooms[i].Pvp = h.Pvp
			existingRooms[i].Session = h.Session
			existingRooms[i].Fo = h.Fo
			existingRooms[i].Password = h.Password
			existingRooms[i].GUID = h.GUID
			existingRooms[i].Maxconnections = h.Maxconnections
			existingRooms[i].Dedicated = h.Dedicated
			existingRooms[i].Clienthosted = h.Clienthosted
			existingRooms[i].Connected = h.Connected
			existingRooms[i].Mode = h.Mode
			existingRooms[i].Port = h.Port
			existingRooms[i].V = h.V
			existingRooms[i].Tags = h.Tags
			existingRooms[i].Session = h.Session
			existingRooms[i].Lanonly = h.Lanonly
			existingRooms[i].Intent = h.Intent
			existingRooms[i].Allownewplayers = h.Allownewplayers
			existingRooms[i].Serverpaused = h.Serverpaused
			existingRooms[i].Steamid = h.Steamid

			updateRooms = append(updateRooms, existingRooms[i])
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

	var batchSize = 1000

	if len(deleteRooms) > 0 {
		l.DB.Delete(&deleteRooms)
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
				if err := l.DB.Save(&batch).Error; err != nil {
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
			if err := l.DB.Create(&batch).Error; err != nil {
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

func (l *LobbyServer) requestLobbyV2Api(region string) (LobbyHomeListBody, error) {
	url := fmt.Sprintf(LobbyListUrl, region)
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
	Colour     string
	EventLevel int
	Name       string
	NetID      string
	Prefab     string
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
		fmt.Printf("Player %d: %+v\n", i, p)
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
