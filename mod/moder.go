package mod

import (
	"dst-admin-go/entity"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	lua "github.com/yuin/gopher-lua"

	"math"
	"net/http"
	"net/url"
)

const (
	steamAPIKey = "73DF9F781D195DFD3D19DED1CB72EEE6"
	appID       = 322330
	language    = 6
)

type searchResult struct {
	Page      int       `json:"page"`
	Size      int       `json:"size"`
	Total     int       `json:"total"`
	TotalPage int       `json:"totalPage"`
	Data      []ModInfo `json:"data"`
}

// ModInfo 存储 mod 信息的结构体
type ModInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
	Desc   string `json:"desc"`
	Time   int    `json:"time"`
	Sub    int    `json:"sub"`
	Img    string `json:"img"`
	Vote   struct {
		Star int `json:"star"`
		Num  int `json:"num"`
	} `json:"vote"`
	Child []string `json:"child,omitempty"`
}

func get_mod_info_config(mod_id string) map[string]interface{} {
	// 检查 mod 文件是否已经存在
	mod_download_path := "/root/mine/dst/dst-admin-go-1.0.0/go-mod/mod"
	mod_path := mod_download_path + "/steamapps/workshop/content/322330/" + mod_id
	if _, err := os.Stat(mod_path); err == nil {
		fmt.Println("Mod already downloaded to:", mod_path)
	} else {
		// 调用 SteamCMD 命令下载 mod
		cmd := exec.Command("/root/steamcmd/steamcmd.sh", "+login anonymous", "+force_install_dir", mod_download_path, "+workshop_download_item 322330 "+mod_id, "+quit")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error executing command:", err)
			return make(map[string]interface{})
		}

		// 解析 SteamCMD 输出，提取 mod 文件路径
		re := regexp.MustCompile(`Downloaded item \d+ to "([^"]+)"`)
		match := re.FindStringSubmatch(string(output))
		if len(match) < 2 {
			fmt.Println("Error parsing output")
			return make(map[string]interface{})
		}
		path := match[1]
		fmt.Println("Mod downloaded to:", path)
	}

	// 查找 modinfo.lua 文件
	modinfo_path := filepath.Join(mod_path, "modinfo.lua")
	if _, err := os.Stat(modinfo_path); err != nil {
		fmt.Println("Error finding modinfo.lua:", err)
		return make(map[string]interface{})
	}
	return read_mod_info(modinfo_path)
}

func read_mod_info(modinfo_path string) map[string]interface{} {
	// 读取 modinfo.lua 文件内容
	script, err := ioutil.ReadFile(modinfo_path)
	if err != nil {
		fmt.Println("Error reading modinfo.lua:", err)
		return make(map[string]interface{})
	}

	// 打印 modinfo.lua 文件内容
	// fmt.Println("Modinfo.lua content:")
	// fmt.Println(string(content))

	return parseModInfoLua(string(script))
}

func parseModInfoLua(script string) map[string]interface{} {
	L := lua.NewState()
	defer L.Close()

	// 运行Lua脚本文件
	L.DoString(script)
	// if err := L.DoFile("hello.lua"); err != nil {
	//      panic(err)
	// }

	// 获取所有全局变量
	global := L.Get(lua.GlobalsIndex).(*lua.LTable)
	m := make(map[string]interface{})
	global.ForEach(func(k lua.LValue, v lua.LValue) {
		if !excludeList[k.String()] && v.Type() != lua.LTFunction {
			// data, _ := json.Marshal(toInterface(v))
			// fmt.Printf("%s = %v\n", k.String(), string(data))
			m[k.String()] = toInterface(v)
		}
	})

	// data, _ := json.Marshal(m)
	// fmt.Println(string(data))
	return m
}

// 定义需要排除的 Lua 自带对象的名称
var excludeList = map[string]bool{
	"_G": true, "assert": true, "collectgarbage": true, "dofile": true, "error": true,
	"getmetatable": true, "ipairs": true, "load": true, "loadfile": true, "module": true,
	"next": true, "pairs": true, "pcall": true, "print": true, "rawequal": true, "rawget": true,
	"rawset": true, "require": true, "select": true, "setmetatable": true, "tonumber": true,
	"tostring": true, "type": true, "unpack": true, "xpcall": true, "debug": true, "_VERSION": true,
	"os": true, "_GOPHER_LUA_VERSION": true, "string": true, "math": true, "io": true, "channel": true,
	"package": true, "coroutine": true, "table": true,
}

// 将 Lua 值转换为 interface{} 类型
func toInterface(lv lua.LValue) interface{} {
	switch lv.Type() {
	case lua.LTNil:
		return nil
	case lua.LTBool:
		return bool(lv.(lua.LBool))
	case lua.LTNumber:
		return float64(lv.(lua.LNumber))
	case lua.LTString:
		return lv.String()
	case lua.LTTable:
		t := lv.(*lua.LTable)
		if isTableArray(t) {
			arr := make([]interface{}, t.Len())
			t.ForEach(func(i lua.LValue, v lua.LValue) {
				index := int(float64(i.(lua.LNumber)) - 1)
				arr[index] = toInterface(v)
			})
			return arr
		} else {
			return toMap(t)
		}
	default:
		return lv.String()
	}
}

// 将 Lua table 转换为 map[string]interface{}
func toMap(t *lua.LTable) map[string]interface{} {
	m := make(map[string]interface{})
	t.ForEach(func(k lua.LValue, v lua.LValue) {
		key := ""
		switch k.Type() {
		case lua.LTString:
			key = k.String()
		case lua.LTNumber:
			key = fmt.Sprintf("%g", float64(k.(lua.LNumber)))
		default:
			key = fmt.Sprintf("%v", k)
		}
		m[key] = toInterface(v)
	})
	return m
}

// 判断 Lua table 是否为数组类型
func isTableArray(t *lua.LTable) bool {
	maxIndex := 0
	isSequential := true
	t.ForEach(func(k lua.LValue, v lua.LValue) {
		if i, ok := k.(lua.LNumber); ok {
			if i != lua.LNumber(int(i)) {
				isSequential = false
			} else if int(i) > maxIndex {
				maxIndex = int(i)
			}
		} else {
			isSequential = false
		}
	})
	return isSequential && maxIndex == t.Len()
}

// SearchModList 搜索 mod 的函数
func SearchModList(text string, page int, num int) (map[string]interface{}, error) {
	urlStr := "http://api.steampowered.com/IPublishedFileService/QueryFiles/v1/"
	data := url.Values{
		"page":             {fmt.Sprintf("%d", page)},
		"key":              {steamAPIKey},
		"appid":            {"322330"},
		"language":         {"6"},
		"return_tags":      {"true"},
		"numperpage":       {fmt.Sprintf("%d", num)},
		"search_text":      {text},
		"return_vote_data": {"true"},
		"return_children":  {"true"},
	}
	urlStr = urlStr + "?" + data.Encode()

	var modData map[string]interface{}
	for i := 0; i < 2; i++ {
		resp, err := http.Get(urlStr)
		if err != nil {
			fmt.Println("搜索 mod 失败")
			return nil, err
		}
		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(&modData)
		if err != nil {
			fmt.Println("解析 mod 数据失败")
			return nil, err
		}
		if modData["response"] != nil {
			break
		}
	}
	if modData["response"] == nil {
		return nil, fmt.Errorf("no response found in mod data")
	}

	modResponse := modData["response"].(map[string]interface{})
	total := int(modResponse["total"].(float64))
	modInfoRaw := modResponse["publishedfiledetails"].([]interface{})

	modList := make([]ModInfo, 0)
	if len(modInfoRaw) > 0 {
		for _, modInfoRaw := range modInfoRaw {
			modInfo := modInfoRaw.(map[string]interface{})
			img := modInfo["preview_url"].(string)
			voteData := modInfo["vote_data"].(map[string]interface{})
			auth := modInfo["creator"].(string)
			var authorURL string
			if auth != "" {
				authorURL = fmt.Sprintf("https://steamcommunity.com/profiles/%s/?xml=1", auth)
			}
			mod := ModInfo{
				ID:     fmt.Sprintf("%v", modInfo["publishedfileid"]),
				Name:   fmt.Sprintf("%v", modInfo["title"]),
				Author: authorURL,
				Desc:   fmt.Sprintf("%v", modInfo["file_description"]),
				Time:   int(modInfo["time_updated"].(float64)),
				Sub:    int(modInfo["subscriptions"].(float64)),
				Img:    img,
				Vote: struct {
					Star int `json:"star"`
					Num  int `json:"num"`
				}{
					Star: int(voteData["score"].(float64)*5) + 1,
					Num:  int(voteData["votes_up"].(float64) + voteData["votes_down"].(float64)),
				},
			}
			if modInfo["num_children"].(float64) != 0 {
				children := modInfo["children"].([]interface{})
				child := make([]string, len(children))
				for i, c := range children {
					child[i] = fmt.Sprintf("%v", c.(map[string]interface{})["publishedfileid"])
				}
				mod.Child = child
			}
			modList = append(modList, mod)
		}
	}

	return map[string]interface{}{
		"page":      page,
		"size":      num,
		"total":     total,
		"totalPage": int(math.Ceil(float64(total) / float64(num))),
		"data":      modList,
	}, nil
}

func GetModInfo(modID string) entity.ModInfo {
	urlStr := "http://api.steampowered.com/IPublishedFileService/GetDetails/v1/"
	data := url.Values{}
	data.Set("key", steamAPIKey)
	data.Set("language", "6")
	data.Set("publishedfileids[0]", modID)
	urlStr = urlStr + "?" + data.Encode()

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		fmt.Println(err)
		return entity.ModInfo{}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return entity.ModInfo{}
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return entity.ModInfo{}
	}

	dataList, ok := result["response"].(map[string]interface{})["publishedfiledetails"].([]interface{})
	if !ok || len(dataList) == 0 {
		fmt.Println("get mod error")
		return entity.ModInfo{}
	}

	data2 := dataList[0].(map[string]interface{})
	img := data2["preview_url"].(string)
	auth := data2["creator"].(string)
	var authorURL string
	if auth != "" {
		authorURL = fmt.Sprintf("https://steamcommunity.com/profiles/%s/?xml=1", auth)
	} else {
		authorURL = ""
	}

	modId := data2["publishedfileid"].(string)
	name := data2["title"].(string)
	last_time := data2["time_updated"].(float64)
	description := data2["file_description"].(string)
	auth = authorURL
	file_url := data2["file_url"]
	img = fmt.Sprintf("%s?imw=64&imh=64&ima=fit&impolicy=Letterbox&imcolor=%%23000000&letterbox=true", img)
	v := getVersion(data["tags"])
	creator_appid := data2["creator_appid"].(float64)
	consumer_appid := data2["consumer_appid"].(float64)

	// modInfoRaw := map[string]interface{}{
	// 	"id":             data2["publishedfileid"].(string),
	// 	"name":           data2["title"].(string),
	// 	"last_time":      data2["time_updated"].(int64),
	// 	"description":    data2["file_description"].(string),
	// 	"auth":           authorURL,
	// 	"file_url":       data2["file_url"],
	// 	"img":            fmt.Sprintf("%s?imw=64&imh=64&ima=fit&impolicy=Letterbox&imcolor=%%23000000&letterbox=true", img),
	// 	"v":              getVersion(data["tags"]),
	// 	"creator_appid":  data2["creator_appid"].(int64),
	// 	"consumer_appid": data2["consumer_appid"].(int64),
	// 	 "mod_config":     get_mod_info_config(modID),
	// }

	if modInfo, ok := getModInfoConfig(modID, last_time); ok {
		return modInfo
	}
	var fileUrl = ""
	if file_url != nil {
		fileUrl = file_url.(string)
	}
	modConfigJson, _ := json.Marshal(get_mod_info_config(modID))
	newModInfo := entity.ModInfo{
		Auth:          auth,
		ConsumerAppid: consumer_appid,
		CreatorAppid:  creator_appid,
		Description:   description,
		FileUrl:       fileUrl,
		Modid:         modId,
		Img:           img,
		LastTime:      last_time,
		Name:          name,
		V:             v,
		ModConfig:     string(modConfigJson),
	}

	db := entity.DB
	db.Create(&newModInfo)
	return newModInfo
}

func getModInfoConfig(modid string, lastTime float64) (entity.ModInfo, bool) {
	db := entity.DB
	modInfo := entity.ModInfo{}
	db.Where("modid = ? and last_time = ?", modid, lastTime).First(&modInfo)

	if modInfo.Modid == "" {
		return modInfo, false
	}
	return modInfo, true
}

func getVersion(tags interface{}) string {
	tagList, ok := tags.([]interface{})
	if !ok {
		return ""
	}
	for _, tag := range tagList {
		tagStr := tag.(map[string]interface{})["tag"].(string)
		if len(tagStr) > 8 && tagStr[:8] == "version:" {
			return tagStr[8:]
		}
	}
	return ""
}
