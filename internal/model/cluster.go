package model

import "gorm.io/gorm"

type Cluster struct {
	gorm.Model
	ClusterName     string `gorm:"uniqueIndex" json:"clusterName"`
	Description     string `json:"description"`
	SteamCmd        string `json:"steamcmd"`
	ForceInstallDir string `json:"force_install_dir"`
	Backup          string `json:"backup"`
	ModDownloadPath string `json:"mod_download_path"`
	Uuid            string `json:"uuid"`
	Beta            int    `json:"beta"`
	Bin             int    `json:"bin"`

	Ugc_directory           string `json:"ugc_directory"`
	Persistent_storage_root string `json:"persistent_storage_root"`
	Conf_dir                string `json:"conf_dir"`
}
