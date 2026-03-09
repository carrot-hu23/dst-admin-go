package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	BindAddress       string `yaml:"bindAddress"`
	Port              string `yaml:"port"`
	Path              string `yaml:"path"`
	Db                string `yaml:"database"`
	SteamAPIKey       string `yaml:"steamAPIKey"`
	Flag              string `yaml:"flag"`
	WanIP             string `yaml:"wanip"`
	WhiteAdminIP      string `yaml:"whiteadminip"`
	Token             string `yaml:"token"`
	AutoUpdateModinfo struct {
		Enable              bool `yaml:"enable"`
		CheckInterval       int  `yaml:"checkInterval"`
		UpdateCheckInterval int  `yaml:"updateCheckInterval"`
	} `yaml:"autoUpdateModinfo"`
}

const (
	ConfigPath  = "./config.yml"
	DefaultPort = "8083"
)

var Cfg *Config

func Load() *Config {
	yamlFile, err := ioutil.ReadFile(ConfigPath)
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
	if c.AutoUpdateModinfo.UpdateCheckInterval == 0 {
		c.AutoUpdateModinfo.UpdateCheckInterval = 10
	}
	if c.AutoUpdateModinfo.CheckInterval == 0 {
		c.AutoUpdateModinfo.CheckInterval = 5
	}
	Cfg = c
	return c
}
