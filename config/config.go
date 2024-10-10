package config

type Config struct {
	Port    string `yaml:"port"`
	Db      string `yaml:"database"`
	Token   string `yaml:"token"`
	Collect int    `yaml:"collect"`
}
