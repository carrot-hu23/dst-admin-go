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
}
