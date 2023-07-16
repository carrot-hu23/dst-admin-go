package vo

import (
	"dst-admin-go/constant"
	"dst-admin-go/utils/systemUtils"
	"path/filepath"
)

type ClusterDashboardVO struct {
	IsInstallDst bool                  `json:"isInstallDst"`
	MasterStatus bool                  `json:"masterStatus"`
	CavesStatus  bool                  `json:"cavesStatus"`
	HostInfo     *systemUtils.HostInfo `json:"host"`
	CpuInfo      *systemUtils.CpuInfo  `json:"cpu"`
	MemInfo      *systemUtils.MemInfo  `json:"mem"`
	DiskInfo     *systemUtils.DiskInfo `json:"disk"`
	MemStates    uint64                `json:"memStates"`
	MasterLog    string                `json:"masterLog"`
	CavesLog     string                `json:"cavesLog"`
	Version      int64                 `json:"version"`

	MasterPs  *DstPsVo `json:"masterPs"`
	CavesPs   *DstPsVo `json:"cavesPs"`
	IpConnect string   `json:"ipConnect"`
}

func NewDashboardVO(clusterName string) *ClusterDashboardVO {
	return &ClusterDashboardVO{
		MasterLog: filepath.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether", clusterName, "Master", "server_log.txt"),
		CavesLog:  filepath.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether", clusterName, "Caves", "server_log.txt"),
	}
}
