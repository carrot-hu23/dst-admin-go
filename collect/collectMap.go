package collect

import (
	"dst-admin-go/utils/systemUtils"
	"path/filepath"
	"sync"
)

type CollectMap struct {
	cache sync.Map
}

func NewCollectMap() *CollectMap {
	return &CollectMap{
		cache: sync.Map{},
	}
}

func (cm *CollectMap) AddNewCollect(clusterName string) {
	_, ok := cm.cache.Load(clusterName)
	if !ok {
		home, err := systemUtils.Home()
		if err != nil {
			panic("Home path error: " + err.Error())
		}
		baseLogPath := filepath.Join(home, ".klei/DoNotStarveTogether", clusterName)
		collect := NewCollect(baseLogPath, clusterName)
		collect.StartCollect()
		cm.cache.Store(clusterName, collect)
	}
}

func (cm *CollectMap) RemoveCollect(clusterName string) {
	value, loaded := cm.cache.LoadAndDelete(clusterName)
	if loaded {
		value.(*Collect).Stop()
	}
}
