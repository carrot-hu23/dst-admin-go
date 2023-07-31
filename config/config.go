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

	Token string `yaml:"token"`

	AutoCheck struct {
		MasterInterval     int    `yaml:"masterInterval"`
		CavesInterval      int    `yaml:"cavesInterval"`
		MasterModInterval  int    `yaml:"masterModInterval"`
		CavesModInterval   int    `yaml:"cavesModInterval"`
		GameUpdateInterval int    `yaml:"gameUpdateInterval"`
		ModUpdatePrompt    string `yaml:"modUpdatePrompt"`
		GameUpdatePrompt   string `yaml:"gameUpdatePrompt"`
	} `yaml:"autoCheck"`
}
