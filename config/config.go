package config

type Config struct {
	Port           string `yaml:"port"`
	Path           string `yaml:"path"`
	Db             string `yaml:"database"`
	Steamcmd       string `yaml:"steamcmd"`
	SteamAPIKey    string `yaml:"steamAPIKey"`
	OPENAI_API_KEY string `yaml:"OPENAI_API_KEY"`
	Prompt         string `yaml:"prompt"`
	Flag           string `yaml:"flag"`
	WanIP           string `yaml:"wanip"`
	Token string `yaml:"token"`

	AutoUpdateModinfo struct {
		Enable              bool `yaml:"enable"`
		CheckInterval       int  `yaml:"checkInterval"`
		UpdateCheckInterval int  `yaml:"updateCheckInterval"`
	} `yaml:"autoUpdateModinfo"`

	DstCliPort string `yaml:"dstCliPort"`
}
