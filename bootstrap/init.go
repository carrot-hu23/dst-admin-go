package bootstrap

import (
	"dst-admin-go/autoCheck"
	"dst-admin-go/bot/discord"
	"dst-admin-go/collect"
	"dst-admin-go/config"
	"dst-admin-go/config/database"
	"dst-admin-go/config/global"
	"dst-admin-go/model"
	"dst-admin-go/schedule"
	"dst-admin-go/utils/dstConfigUtils"
	"dst-admin-go/utils/systemUtils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const logPath = "./dst-admin-go.log"

var f *os.File

func Init() {

	initConfig()
	initLog()
	initDB()
	initCollect()
	initSchedule()
}

func initDB() {
	db, err := gorm.Open(sqlite.Open(global.Config.Db), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		panic("failed to connect database")
	}
	database.DB = db
	err = database.DB.AutoMigrate(
		&model.Spawn{},
		&model.PlayerLog{},
		&model.Connect{},
		&model.Regenerate{},
		&model.Proxy{},
		&model.ModInfo{},
		&model.Cluster{},
		&model.JobTask{},
		&model.AutoCheck{},
		&model.Announce{},
		&model.WebLink{},
	)
	if err != nil {
		return
	}
}

func initConfig() {
	yamlFile, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		fmt.Println(err.Error())
	}
	var _config *config.Config
	err = yaml.Unmarshal(yamlFile, &_config)
	if err != nil {
		fmt.Println(err.Error())
	}
	if _config.AutoCheck.MasterInterval == 0 {
		_config.AutoCheck.MasterInterval = 5
	}
	if _config.AutoCheck.CavesInterval == 0 {
		_config.AutoCheck.CavesInterval = 5
	}
	if _config.AutoCheck.MasterModInterval == 0 {
		_config.AutoCheck.MasterModInterval = 10
	}
	if _config.AutoCheck.CavesModInterval == 0 {
		_config.AutoCheck.CavesModInterval = 10
	}
	if _config.AutoCheck.GameUpdateInterval == 0 {
		_config.AutoCheck.GameUpdateInterval = 20
	}
	log.Println("config: ", _config)
	global.Config = _config
}

func initLog() {
	var err error
	f, err = os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
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

func initCollect() {

	//var clusters []model.Cluster
	//database.DB.Find(&clusters)
	//for _, cluster := range clusters {
	//	global.CollectMap.AddNewCollect(cluster.ClusterName)
	//}

	home, _ := systemUtils.Home()
	dstConfig := dstConfigUtils.GetDstConfig()
	clusterName := dstConfig.Cluster
	newCollect := collect.NewCollect(filepath.Join(home, ".klei/DoNotStarveTogether", clusterName), clusterName)
	newCollect.StartCollect()
	global.Collect = newCollect

	// autoCheck.AutoCheckObject = autoCheck.NewAutoCheckConfig(clusterName, dstConfig.Bin, dstConfig.Beta)

	autoCheckManager := autoCheck.AutoCheckManager{}
	autoCheckManager.Start()
	autoCheck.Manager = &autoCheckManager
}

func initSchedule() {
	schedule.ScheduleSingleton = schedule.NewSchedule()
	// service.InitAnnounce()

	discord.DiscordClient = discord.NewDiscordClient(discord.Token)
	log.Println("开始启动discord机器人")
	go discord.DiscordClient.Start()

}
