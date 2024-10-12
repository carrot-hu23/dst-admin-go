package config

type Zone struct {
	Name     string `yaml:"name" json:"name"`
	ZoneCode string `yaml:"zoneCode" json:"zoneCode"`
	Ip       string `yaml:"ip" json:"ip"`
	Port     int    `yaml:"port" json:"port"`
}

type Config struct {
	Port         string `yaml:"port"`
	Db           string `yaml:"database"`
	Token        string `yaml:"token"`
	Collect      int    `yaml:"collect"`
	CheckExpired int    `yaml:"checkExpired"`
	Zones        []Zone `yaml:"zones"`
}
