package main

import (
	"dst-admin-go/config"
	"dst-admin-go/config/database"
	"dst-admin-go/config/global"
	"dst-admin-go/model"
	"dst-admin-go/router"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"io/ioutil"
	"log"
	"net/http/httputil"
	"net/url"
	"os"
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
	global.Config = configData
}

func initDB() {
	db, err := gorm.Open(sqlite.Open(configData.Db), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}
	database.DB = db
	err = database.DB.AutoMigrate(&model.Spawn{}, &model.PlayerLog{}, &model.Connect{}, &model.Proxy{}, &model.ModInfo{}, &model.Cluster{})
	if err != nil {
		return
	}

	proxyEntities := []model.Proxy{}
	db.Find(&proxyEntities)

	if len(proxyEntities) > 0 {
		for _, proxyEntity := range proxyEntities {
			r, e := url.Parse("http://" + proxyEntity.Ip + ":" + proxyEntity.Port)
			if e != nil {
				panic(e)
			}
			p := httputil.NewSingleHostReverseProxy(r)
			global.RoutingTable[proxyEntity.Name] = &global.Route{Proxy: p, Url: r}
		}
	}

	//---------------
	var clusters []model.Cluster
	db.Find(&clusters)
	for _, cluster := range clusters {
		global.CollectMap.AddNewCollect(cluster.ClusterName)
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
	initDB()
}

func main() {

	defer func() {
		f.Close()
	}()

	fmt.Println(":pig, 你是好人")

	app := router.NewRoute()
	pprof.Register(app)
	app.Run(":" + configData.Port)

}
