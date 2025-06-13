package bootstrap

import (
	"dst-admin-go/autoCheck"
	"dst-admin-go/config"
	"dst-admin-go/config/database"
	"dst-admin-go/config/global"
	"dst-admin-go/mod"
	"dst-admin-go/model"
	"dst-admin-go/schedule"
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
	"time"
)

const logPath = "./dst-admin-go.log"

var f *os.File

func Init() {

	initConfig()
	initLog()
	initDB()
	initCollect()
	initSchedule()

	initUpdateModinfos()
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
		&model.User{},
		&model.UserCluster{},
		&model.KV{},
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
	if _config.AutoUpdateModinfo.CheckInterval == 0 {
		_config.AutoUpdateModinfo.CheckInterval = 5
	}
	if _config.AutoUpdateModinfo.UpdateCheckInterval == 0 {
		_config.AutoUpdateModinfo.UpdateCheckInterval = 10
	}
	if _config.Token == "" {
		_config.Token = "pds-g^KU_qE7e8rv1^VVrVXd/01kBDicd7UO5LeL+uYZH1+geZlrutzItvOaw="
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

	var clusters []model.Cluster
	database.DB.Find(&clusters)
	for _, cluster := range clusters {
		if cluster.ClusterType == "远程" {
			continue
		}
		global.CollectMap.AddNewCollect(cluster.ClusterName)
	}

	autoCheckManager := autoCheck.AutoCheckManager{}
	autoCheckManager.Start()
	autoCheck.Manager = &autoCheckManager
}

func initSchedule() {
	schedule.ScheduleSingleton = schedule.NewSchedule()
	// service.InitAnnounce()

}

func initUpdateModinfos() {
	if global.Config.AutoUpdateModinfo.Enable {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Println(r)
				}
			}()
			t := global.Config.AutoUpdateModinfo.CheckInterval
			ticker := time.NewTicker(time.Duration(t) * time.Minute)
			for {
				select {
				case <-ticker.C:
					log.Println("正在定时检查模组配置是否更新 间隔: ", 5, "分钟")
					// 每隔10分钟执行的任务
					mod.CheckModInfoUpdate()
				}
			}
		}()
	}
}
