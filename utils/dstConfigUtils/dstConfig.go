package dstConfigUtils

import (
	"dst-admin-go/utils/fileUtils"
	"log"
	"strings"
)

type DstConfig struct {
	Steamcmd            string `json:"steamcmd"`
	Force_install_dir   string `json:"force_install_dir"`
	DoNotStarveTogether string `json:"doNotStarveTogether"`
	Cluster             string `json:"cluster"`
	Backup              string `json:"backup"`
}

const dst_config_path = "./dst_config"

func NewDstConfig() *DstConfig {
	return &DstConfig{}
}

func GetDstConfig() DstConfig {

	dst_config := NewDstConfig()

	//判断是否存在，不存在创建一个
	if !fileUtils.Exists(dst_config_path) {
		if err := fileUtils.CreateFile(dst_config_path); err != nil {
			log.Panicln("create dst_config error", err)
		}

	}
	data, err := fileUtils.ReadLnFile(dst_config_path)
	if err != nil {
		log.Panicln("read dst_config error", err)
	}
	for _, value := range data {
		if value == "" {
			continue
		}
		if strings.Contains(value, "steamcmd") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				dst_config.Steamcmd = strings.Replace(s, "\\n", "", -1)
			}
		}
		if strings.Contains(value, "force_install_dir") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				dst_config.Force_install_dir = strings.Replace(s, "\\n", "", -1)
			}
		}
		if strings.Contains(value, "doNotStarveTogether") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				dst_config.DoNotStarveTogether = strings.Replace(s, "\\n", "", -1)
			}
		}
		if strings.Contains(value, "cluster") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				dst_config.Cluster = strings.Replace(s, "\\n", "", -1)
			}
		}
		if strings.Contains(value, "backup") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				dst_config.Backup = strings.Replace(s, "\\n", "", -1)
			}
		}
	}

	return *dst_config
}

func SaveDstConfig(dstConfig *DstConfig) {
	log.Println(dstConfig)

	error := fileUtils.WriterLnFile(dst_config_path, []string{
		"steamcmd=" + dstConfig.Steamcmd,
		"force_install_dir=" + dstConfig.Force_install_dir,
		"doNotStarveTogether=" + dstConfig.DoNotStarveTogether,
		"cluster=" + dstConfig.Cluster,
		"backup=" + dstConfig.Backup,
	})
	if error != nil {
		log.Panicln("write dst_config error:", error)
	}
}
