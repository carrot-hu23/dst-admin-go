package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	BindAddress       string `yaml:"bindAddress"`
	Port              string `yaml:"port"`
	Path              string `yaml:"path"`
	Db                string `yaml:"database"`
	DataDir           string `yaml:"dataDir"`
	SteamAPIKey       string `yaml:"steamAPIKey"`
	Flag              string `yaml:"flag"`
	WanIP             string `yaml:"wanip"`
	WhiteAdminIP      string `yaml:"whiteadminip"`
	Token             string `yaml:"token"`
	DstVersionUrl     string `yaml:"dstVersionUrl"`
	AutoUpdateModinfo struct {
		Enable              bool `yaml:"enable"`
		CheckInterval       int  `yaml:"checkInterval"`
		UpdateCheckInterval int  `yaml:"updateCheckInterval"`
	} `yaml:"autoUpdateModinfo"`
}

const (
	DefaultConfigPath = "./config.yml"
	DefaultPort       = "8082"
	DefaultDataDir    = "./"
)

var Cfg *Config

func Load() *Config {
	yamlFile, err := ioutil.ReadFile(DefaultConfigPath)
	if err != nil {
		fmt.Println(err.Error())
	}
	var c *Config
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		fmt.Println(err.Error())
	}
	if c.Port == "" {
		c.Port = DefaultPort
	}
	if c.DataDir == "" {
		c.DataDir = DefaultDataDir
	}
	if c.AutoUpdateModinfo.UpdateCheckInterval == 0 {
		c.AutoUpdateModinfo.UpdateCheckInterval = 10
	}
	if c.AutoUpdateModinfo.CheckInterval == 0 {
		c.AutoUpdateModinfo.CheckInterval = 5
	}
	if c.DstVersionUrl == "" {
		c.DstVersionUrl = "http://ver.tugos.cn/getLocalVersion"
	}
	Cfg = c
	return c
}

// GetDbPath 获取数据库文件的完整路径
func (c *Config) GetDbPath() string {
	return filepath.Join(c.DataDir, c.Db)
}
