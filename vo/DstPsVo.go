package vo

// import (
// 	"dst-admin-go/constant"
// 	"dst-admin-go/utils/systemUtils"
// )

type DstPsVo struct {
	CpuUage string `json:"cpuUage"`
	MemUage string `json:"memUage"`
	VSZ     string `json:"VSZ"`
	RSS     string `json:"RSS"`
}

func NewDstPsVo() *DstPsVo {
	return &DstPsVo{}
}
