package config

type Config struct {
	Port string `yaml:"port"`
	//Db           string `yaml:"database"`
	Token        string `yaml:"token"`
	Collect      int    `yaml:"collect"`
	CheckExpired int    `yaml:"checkExpired"`
	Database     struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`
}
