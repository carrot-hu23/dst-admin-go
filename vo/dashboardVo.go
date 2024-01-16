package vo

import (
	"dst-admin-go/utils/systemUtils"
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
