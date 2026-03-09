package game

type DstPsAux struct {
	CpuUage string `json:"cpuUage"`
	MemUage string `json:"memUage"`
	VSZ     string `json:"VSZ"`
	RSS     string `json:"RSS"`
}

type Process interface {
	SessionName(clusterName, levelName string) string

	Start(clusterName, levelName string) error
	Stop(clusterName, levelName string) error
	StartAll(clusterName string) error
	StopAll(clusterName string) error

	Status(clusterName, levelName string) (bool, error)

	Command(clusterName, levelName, command string) error

	PsAuxSpecified(clusterName, levelName string) DstPsAux
}
