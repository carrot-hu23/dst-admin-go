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

	// ContainerMode 由 DST_ADMIN_CONTAINER_MODE 环境变量决定，不写入配置文件。
	// 为 true 时，前端会锁定 steamcmd/force_install_dir/cluster 等由 Docker 镜像
	// 固定路径决定的字段，改用 volume 挂载来适配不同的存档/安装路径。
	ContainerMode bool `json:"container_mode"`
}

type Config interface {
	GetDstConfig(clusterName string) (DstConfig, error)
	SaveDstConfig(clusterName string, config DstConfig) error
}
