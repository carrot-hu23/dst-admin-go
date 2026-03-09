package game

import (
	"dst-admin-go/internal/service/dstConfig"
	"dst-admin-go/internal/service/levelConfig"
	"fmt"
	"log"
	"sync"
)

type WindowProcess struct {
	dstConfig        dstConfig.Config
	cli              *ClusterContainer
	levelConfigUtils *levelConfig.LevelConfigUtils
}

func NewWindowProcess(dstConfig *dstConfig.Config, levelConfigUtils *levelConfig.LevelConfigUtils) *WindowProcess {
	return &WindowProcess{
		dstConfig:        *dstConfig,
		cli:              NewClusterContainer(),
		levelConfigUtils: levelConfigUtils,
	}
}

func (p *WindowProcess) SessionName(clusterName, levelName string) string {
	return clusterName + "_" + levelName
}

func (p *WindowProcess) Start(clusterName, levelName string) error {
	config, err := p.dstConfig.GetDstConfig(clusterName)
	if err != nil {
		return err
	}
	go func() {
		p.cli.StartLevel(clusterName, levelName, config.Bin, config.Steamcmd, config.Force_install_dir, config.Ugc_directory, config.Persistent_storage_root, config.Conf_dir)
	}()

	return err
}

func (p *WindowProcess) Stop(clusterName, levelName string) error {
	p.cli.StopLevel(clusterName, levelName)
	return nil
}

func (p *WindowProcess) StartAll(clusterName string) error {

	err := p.StopAll(clusterName)
	if err != nil {
		return err
	}
	config, err := p.levelConfigUtils.GetLevelConfig(clusterName)
	if err != nil {
		log.Panicln(err)
	}
	var wg sync.WaitGroup
	wg.Add(len(config.LevelList))
	for i := range config.LevelList {
		go func(i int) {
			defer func() {
				wg.Done()
				if r := recover(); r != nil {
					log.Println(r)
				}
			}()
			levelName := config.LevelList[i].File
			err := p.Start(clusterName, levelName)
			if err != nil {
				log.Println(err)
				return
			}
		}(i)
	}
	wg.Wait()
	return nil
}

func (p *WindowProcess) StopAll(clusterName string) error {

	config, err := p.levelConfigUtils.GetLevelConfig(clusterName)
	if err != nil {
		log.Panicln(err)
	}
	var wg sync.WaitGroup
	wg.Add(len(config.LevelList))
	for i := range config.LevelList {
		go func(i int) {
			defer func() {
				wg.Done()
				if r := recover(); r != nil {
					log.Println(r)
				}
			}()
			levelName := config.LevelList[i].File
			err := p.Stop(clusterName, levelName)
			if err != nil {
				return
			}
		}(i)
	}
	wg.Wait()
	return nil
}

func (p *WindowProcess) Status(clusterName, levelName string) (bool, error) {
	return p.cli.Status(clusterName, levelName), nil
}

func (p *WindowProcess) Command(clusterName, levelName, command string) error {
	p.cli.Send(clusterName, levelName, command)
	return nil
}

func (p *WindowProcess) PsAuxSpecified(clusterName, levelName string) DstPsAux {
	cpuUsage := p.cli.CpuUsage(clusterName, levelName)
	memUsage := p.cli.MemUsage(clusterName, levelName)
	dstPsAux := DstPsAux{}
	dstPsAux.RSS = fmt.Sprintf("%f", memUsage*1024)
	dstPsAux.CpuUage = fmt.Sprintf("%f", cpuUsage)
	return dstPsAux
}
