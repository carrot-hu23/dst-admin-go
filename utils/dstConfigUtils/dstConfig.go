package dstConfigUtils

import (
	"dst-admin-go/constant/consts"
	"dst-admin-go/utils/fileUtils"
	"log"
	"strings"
)

type DstConfig struct {
	Steamcmd                   string `json:"steamcmd"`
	Force_install_dir          string `json:"force_install_dir"`
	DoNotStarveServerDirectory string `json:"donot_starve_server_directory"`
	Persistent_storage_root    string `json:"persistent_storage_root"`
	Cluster                    string `json:"cluster"`
	Backup                     string `json:"backup"`
	Mod_download_path          string `json:"mod_download_path"`
}

const dst_config_path = "./dst_config"

func NewDstConfig() *DstConfig {
	return &DstConfig{}
}

func GetDstConfig() DstConfig {

	dstConfig := NewDstConfig()

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
				dstConfig.Steamcmd = strings.Replace(s, "\\n", "", -1)
			}
		}
		if strings.Contains(value, "force_install_dir") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				dstConfig.Force_install_dir = strings.Replace(s, "\\n", "", -1)
			}
		}
		if strings.Contains(value, "donot_starve_server_directory") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				dstConfig.DoNotStarveServerDirectory = strings.Replace(s, "\\n", "", -1)
			}
		}
		if strings.Contains(value, "persistent_storage_root") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				dstConfig.Persistent_storage_root = strings.Replace(s, "\\n", "", -1)
			}
		}
		if strings.Contains(value, "cluster") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				dstConfig.Cluster = strings.Replace(s, "\\n", "", -1)
			}
		}
		if strings.Contains(value, "backup") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				dstConfig.Backup = strings.Replace(s, "\\n", "", -1)
			}
		}
		if strings.Contains(value, "mod_download_path") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				dstConfig.Mod_download_path = strings.Replace(s, "\\n", "", -1)
			}
		}
	}
	// 设置默认值
	if dstConfig.Cluster == "" {
		dstConfig.Cluster = "Cluster1"
	}
	if dstConfig.Backup == "" {
		dstConfig.Backup = consts.KleiDstPath
	}
	if dstConfig.Mod_download_path == "" {
		dstConfig.Backup = consts.KleiDstPath
	}
	return *dstConfig
}

func SaveDstConfig(dstConfig *DstConfig) {
	log.Println(dstConfig)

	err := fileUtils.WriterLnFile(dst_config_path, []string{
		"steamcmd=" + dstConfig.Steamcmd,
		"force_install_dir=" + dstConfig.Force_install_dir,
		"donot_starve_server_directory=" + dstConfig.DoNotStarveServerDirectory,
		"persistent_storage_root=" + dstConfig.Persistent_storage_root,
		"cluster=" + dstConfig.Cluster,
		"backup=" + dstConfig.Backup,
		"mod_download_path=" + dstConfig.Mod_download_path,
	})
	if err != nil {
		log.Panicln("write dst_config error:", err)
	}
}
