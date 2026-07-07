package dstConfig

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

type Config interface {
	GetDstConfig(clusterName string) (DstConfig, error)
	SaveDstConfig(clusterName string, config DstConfig) error
}
