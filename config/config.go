package config

type Config struct {
	Port        string `yaml:"port"`
	Path        string `yaml:"path"`
	Db          string `yaml:"db"`
	Steamcmd    string `yaml:"steamcmd"`
	SteamAPIKey string `yaml:"steamAPIKey"`
}
