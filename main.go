package main

import (
	"dst-admin-go/collect"
	"dst-admin-go/config"
	"dst-admin-go/config/global"
	"dst-admin-go/constant"
	"dst-admin-go/entity"
	"dst-admin-go/route"
	"dst-admin-go/utils/dstConfigUtils"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const log_path = "./dst-admin-go.log"

var f *os.File

var configData *config.Config

func InitConfig() {

	yamlFile, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		fmt.Println(err.Error())
	}
	var _config *config.Config
	err = yaml.Unmarshal(yamlFile, &_config)
	if err != nil {
		fmt.Println(err.Error())
	}
	configData = _config
	entity.Config = configData
}

func iniiDB() {
	db, err := gorm.Open(sqlite.Open(configData.Db), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}
	entity.DB = db
	entity.DB.AutoMigrate(&entity.Spawn{}, &entity.PlayerLog{}, &entity.Connect{}, &entity.Proxy{}, &entity.ModInfo{})

	proxyEntities := []entity.Proxy{}
	db.Find(&proxyEntities)

	if len(proxyEntities) > 0 {
		for _, proxyEntity := range proxyEntities {
			r, e := url.Parse("http://" + proxyEntity.Ip + ":" + proxyEntity.Port)
			if e != nil {
				panic(e)
			}
			p := httputil.NewSingleHostReverseProxy(r)
			entity.RoutingTable[proxyEntity.Name] = &entity.Route{Proxy: p, Url: r}
		}
	}

}

func logInit() {
	var err error
	f, err = os.OpenFile(log_path, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return
	}

	// 组合一下即可，os.Stdout代表标准输出流
	multiWriter := io.MultiWriter(os.Stdout, f)
	log.SetOutput(multiWriter)

	gin.ForceConsoleColor()
	gin.SetMode(gin.DebugMode)
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func init() {
	logInit()
	InitConfig()
	iniiDB()
}

func main() {

	defer func() {
		f.Close()
	}()

	fmt.Println(":pig, 你是好人")

	baseLogPath := filepath.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether", dstConfigUtils.GetDstConfig().Cluster)

	global.Collect = collect.NewCollect([]string{
		filepath.Join(baseLogPath, "Master", "server_log.txt"),
		filepath.Join(baseLogPath, "Caves", "server_log.txt"),
	}, []string{
		filepath.Join(baseLogPath, "Master", "server_chat_log.txt"),
		filepath.Join(baseLogPath, "Master", "server_chat_log.txt"),
	})
	global.Collect.StartCollect()

	app := route.NewRoute()
	app.Run(":" + configData.Port)

}
