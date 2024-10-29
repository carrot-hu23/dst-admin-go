package model

import "gorm.io/gorm"

type Cluster struct {
	gorm.Model
	ClusterName     string `gorm:"uniqueIndex" json:"clusterName"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	SteamCmd        string `json:"steamcmd"`
	ForceInstallDir string `json:"force_install_dir"`
	Backup          string `json:"backup"`
	ModDownloadPath string `json:"mod_download_path"`
	Uuid            string `json:"uuid"`
	Beta            int    `json:"beta"`
	Bin             int    `json:"bin"`

	LevelType   string `json:"levelType"`
	ClusterType string `json:"clusterType"`
	Ip          string `json:"ip"`
	Port        int    `json:"port"`
	Username    string `json:"username"`
	Password    string `json:"password"`

	Ugc_directory           string `json:"ugc_directory"`
	Persistent_storage_root string `json:"persistent_storage_root"`
	Conf_dir                string `json:"conf_dir"`

	RemoteClusterName string `json:"remoteClusterName"`
	// 远程集群名称列表
	RemoteClusterNameList []string `gorm:"-"`
}
