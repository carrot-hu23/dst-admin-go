package config

type Zone struct {
	Name     string `yaml:"name"`
	ZoneCode string `yaml:"zoneCode"`
	Ip       string `yaml:"ip"`
	Port     int    `yaml:"port"`
}

type Config struct {
	Port         string `yaml:"port"`
	Db           string `yaml:"database"`
	Token        string `yaml:"token"`
	Collect      int    `yaml:"collect"`
	CheckExpired int    `yaml:"checkExpired"`
	Zones        []Zone `yaml:"zones"`
}
