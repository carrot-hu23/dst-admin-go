package vo

import (
	"dst-admin-go/constant"
	"dst-admin-go/utils/dstConfigUtils"
	"dst-admin-go/utils/systemUtils"
	"path/filepath"
)

type DashboardVO struct {
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
	Version      string                `json:"version"`

	MasterPs *DstPsVo `json:"masterPs"`
	CavesPs  *DstPsVo `json:"cavesPs"`
}

func NewDashboardVO() *DashboardVO {
	return &DashboardVO{
		MasterLog: filepath.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether", dstConfigUtils.GetDstConfig().Cluster, "Master", "server_log.txt"),
		CavesLog:  filepath.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether", dstConfigUtils.GetDstConfig().Cluster, "Caves", "server_log.txt"),
	}
}
