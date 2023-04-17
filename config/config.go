package config

type Config struct {
	Port string `yaml:"port"`
	Path string `yaml:"path"`
	Db   string `yaml:"db"`
}
