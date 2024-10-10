package bootstrap

import (
	"dst-admin-go/config"
	"dst-admin-go/config/database"
	"dst-admin-go/config/dockerClient"
	"dst-admin-go/config/global"
	"dst-admin-go/model"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"io/ioutil"
	"log"
	"os"
)

const logPath = "./dst-admin-go.log"

var f *os.File

func Init() {

	initConfig()
	initLog()
	initDB()
	initDockerClient()
}

func initDockerClient() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Panicln(err)
	}
	dockerClient.Client = cli
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
		&model.Cluster{},
		&model.User{},
		&model.UserCluster{},
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
	if _config.Collect == 0 {
		_config.Collect = 30
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
