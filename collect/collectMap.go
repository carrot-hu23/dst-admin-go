package collect

import (
	"dst-admin-go/constant"
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
		baseLogPath := filepath.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether", clusterName)
		collect := NewCollect(baseLogPath)
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
