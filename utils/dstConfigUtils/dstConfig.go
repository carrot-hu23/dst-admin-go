package dstConfigUtils

import (
	"dst-admin-go/constant/consts"
	"dst-admin-go/utils/fileUtils"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

type DstConfig struct {
	Steamcmd                   string `json:"steamcmd"`
	Force_install_dir          string `json:"force_install_dir"`
	DoNotStarveServerDirectory string `json:"donot_starve_server_directory"`
	Cluster                    string `json:"cluster"`
	Backup                     string `json:"backup"`
	Mod_download_path          string `json:"mod_download_path"`
	Bin                        int    `json:"bin"`
	Beta                       int    `json:"beta"`

	Ugc_directory string `json:"ugc_directory"`
	// 根目录位置
	Persistent_storage_root string `json:"persistent_storage_root"`
	// 存档相对位置
	Conf_dir string `json:"conf_dir"`
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
		// TODO 这里解析有问题，如果路径含有 steamcmd 就会存在问题
		if strings.Contains(value, "steamcmd=") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				dstConfig.Steamcmd = strings.Replace(s, "\\n", "", -1)
			}
		}
		if strings.Contains(value, "force_install_dir=") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				dstConfig.Force_install_dir = strings.Replace(s, "\\n", "", -1)
			}
		}
		if strings.Contains(value, "donot_starve_server_directory=") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				dstConfig.DoNotStarveServerDirectory = strings.Replace(s, "\\n", "", -1)
			}
		}
		if strings.Contains(value, "persistent_storage_root=") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				dstConfig.Persistent_storage_root = strings.Replace(s, "\\n", "", -1)
			}
		}
		if strings.Contains(value, "cluster=") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				dstConfig.Cluster = strings.Replace(s, "\\n", "", -1)
			}
		}
		if strings.Contains(value, "backup=") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				dstConfig.Backup = strings.Replace(s, "\\n", "", -1)
			}
		}
		if strings.Contains(value, "mod_download_path=") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				dstConfig.Mod_download_path = strings.Replace(s, "\\n", "", -1)
			}
		}
		if strings.Contains(value, "bin=") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				bin, _ := strconv.ParseInt(strings.Replace(s, "\\n", "", -1), 10, 64)
				dstConfig.Bin = int(bin)
			}
		}
		if strings.Contains(value, "beta=") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				beta, _ := strconv.ParseInt(strings.Replace(s, "\\n", "", -1), 10, 64)
				dstConfig.Beta = int(beta)
			}
		}
		if strings.Contains(value, "ugc_directory=") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				dstConfig.Ugc_directory = strings.Replace(s, "\\n", "", -1)
			}
		}
		if strings.Contains(value, "conf_dir=") {
			split := strings.Split(value, "=")
			if len(split) > 1 {
				s := strings.TrimSpace(split[1])
				dstConfig.Conf_dir = strings.Replace(s, "\\n", "", -1)
			}
		}
	}
	// 设置默认值
	if dstConfig.Cluster == "" {
		dstConfig.Cluster = "Cluster1"
	}
	if dstConfig.Backup == "" {
		defaultPath := filepath.Join(consts.DefaultKleiDstPath, "backup")
		fileUtils.CreateDirIfNotExists(defaultPath)
		dstConfig.Backup = defaultPath
	}
	if dstConfig.Mod_download_path == "" {
		defaultPath := filepath.Join(consts.DefaultKleiDstPath, "mod_config_download")
		fileUtils.CreateDirIfNotExists(defaultPath)
		dstConfig.Mod_download_path = defaultPath
	}
	if dstConfig.Bin == 0 {
		dstConfig.Bin = 32
	}
	return *dstConfig
}

func SaveDstConfig(dstConfig *DstConfig) {
	log.Println(dstConfig)

	oldDstConfig := GetDstConfig()
	if dstConfig.Steamcmd == "" {
		dstConfig.Steamcmd = oldDstConfig.Steamcmd
	}
	if dstConfig.Force_install_dir == "" {
		dstConfig.Force_install_dir = oldDstConfig.Force_install_dir
	}
	//if dstConfig.Ugc_directory == "" {
	//	dstConfig.Ugc_directory = oldDstConfig.Ugc_directory
	//}
	//if dstConfig.Conf_dir == "" {
	//	dstConfig.Conf_dir = oldDstConfig.Conf_dir
	//}
	//if dstConfig.Persistent_storage_root == "" {
	//	dstConfig.Persistent_storage_root = oldDstConfig.Persistent_storage_root
	//}
	if dstConfig.Cluster == "" {
		dstConfig.Cluster = oldDstConfig.Cluster
	}
	if dstConfig.Backup == "" {
		dstConfig.Backup = oldDstConfig.Backup
	}
	if dstConfig.Mod_download_path == "" {
		dstConfig.Mod_download_path = oldDstConfig.Mod_download_path
	}

	err := fileUtils.WriterLnFile(dst_config_path, []string{
		"steamcmd=" + dstConfig.Steamcmd,
		"force_install_dir=" + dstConfig.Force_install_dir,
		"donot_starve_server_directory=" + dstConfig.DoNotStarveServerDirectory,
		"ugc_directory=" + dstConfig.Ugc_directory,
		"conf_dir=" + dstConfig.Conf_dir,
		"persistent_storage_root=" + dstConfig.Persistent_storage_root,
		"cluster=" + dstConfig.Cluster,
		"backup=" + dstConfig.Backup,
		"mod_download_path=" + dstConfig.Mod_download_path,
		"bin=" + strconv.Itoa(dstConfig.Bin),
		"beta=" + strconv.Itoa(dstConfig.Beta),
	})
	if err != nil {
		log.Panicln("write dst_config error:", err)
	}
	// global.Collect.ReCollect(filepath.Join(dstUtils.GetKleiDstPath(), dstConfig.Cluster), dstConfig.Cluster)
}
