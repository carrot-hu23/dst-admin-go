package config

type Config struct {
	BindAddress       string `yaml:"bindAddress"`
	Port              string `yaml:"port"`
	Path              string `yaml:"path"`
	Db                string `yaml:"database"`
	Steamcmd          string `yaml:"steamcmd"`
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

	DstCliPort string `yaml:"dstCliPort"`
}
