package lobbyServer

import (
	"bytes"
	"encoding/json"
	"fmt"
	lua "github.com/yuin/gopher-lua"
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

type LobbyServer struct{}

type LobbyRegionsBody struct {
	LobbyRegions []struct {
		Region string `json:"Region"`
	} `json:"LobbyRegions"`
}

type LobbyHome struct {
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
	Secondaries     interface{} `json:"secondaries"`
	SecondariesJson string      `json:"secondariesJson"`
}

type LobbyHomeListBody struct {
	GET []LobbyHome "GET"
}

type Secondaries struct {
	Addr    string `json:"__addr"`
	ID      string `json:"id"`
	Steamid string `json:"steamid"`
	Port    int    `json:"port"`
}

type LobbyHomeDetail struct {
	Addr            string        `json:"__addr"`
	RowID           string        `json:"__rowId"`
	Host            string        `json:"host"`
	Steamclanid     string        `json:"steamclanid"`
	Clanonly        bool          `json:"clanonly"`
	Platform        int           `json:"platform"`
	Mods            bool          `json:"mods"`
	Name            string        `json:"name"`
	Pvp             bool          `json:"pvp"`
	Session         string        `json:"session"`
	Fo              bool          `json:"fo"`
	Password        bool          `json:"password"`
	GUID            string        `json:"guid"`
	Maxconnections  int           `json:"maxconnections"`
	Dedicated       bool          `json:"dedicated"`
	Clienthosted    bool          `json:"clienthosted"`
	Connected       int           `json:"connected"`
	Mode            string        `json:"mode"`
	Port            int           `json:"port"`
	V               int           `json:"v"`
	Tags            string        `json:"tags"`
	Season          string        `json:"season"`
	Lanonly         bool          `json:"lanonly"`
	Intent          string        `json:"intent"`
	Allownewplayers bool          `json:"allownewplayers"`
	Serverpaused    bool          `json:"serverpaused"`
	Steamid         string        `json:"steamid"`
	Steamroom       string        `json:"steamroom"`
	Secondaries     interface{}   `json:"secondaries"`
	Data            string        `json:"data"`
	Worldgen        string        `json:"worldgen"`
	Players         string        `json:"players"`
	ModsInfo        []interface{} `json:"mods_info"`
	Desc            string        `json:"desc"`
	Tick            int           `json:"tick"`
	Clientmodsoff   bool          `json:"clientmodsoff"`
	Nat             int           `json:"nat"`
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

func (l *LobbyServer) queryLobbyList(playerName string) {
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
							homeDetail := l.GetLobbyHomeInfo(region, lobbyHome.RowID)
							if playerName != "" && strings.Contains(homeDetail.Players, playerName) {
								log.Println(homeDetail.Name)
								l.getPlayer(homeDetail.Players)
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

func (l *LobbyServer) GetLobbyHomeInfo(region, rowID string) LobbyHomeDetail {
	if lobbyHomeDetailBody, err := l.requestLobbyHomeDetailApi(region, rowID); err == nil {
		if len(lobbyHomeDetailBody.GET) > 0 {
			lobbyHomeDetail := lobbyHomeDetailBody.GET[0]
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

func (l *LobbyServer) getPlayer(playerLua string) {
	if playerLua == "return {  }" {

	}
	L := lua.NewState()
	defer L.Close()

	err := L.DoString(playerLua)
	if err != nil {
		fmt.Println("Lua execution error:", err)
		return
	}

	// 提取 Lua 表中的值
	luaTable := L.Get(-1) // 获取栈顶的值
	L.Pop(1)              // 从栈中移除该值

	players, ok := luaTable.(*lua.LTable)
	if !ok {
		fmt.Println("Invalid Lua table")
		return
	}

	// 解析并打印玩家信息
	for i := 1; ; i++ {
		player := players.RawGetInt(i)
		if player == lua.LNil {
			break
		}

		p := l.parsePlayer(player)
		fmt.Printf("Player %d: %+v\n", i, p)
	}
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
