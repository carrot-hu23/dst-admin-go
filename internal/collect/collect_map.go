package collect

import (
	"dst-admin-go/internal/service/archive"
	"sync"
)

var CollectorMap *CollectMap

type CollectMap struct {
	cache   sync.Map
	archive *archive.PathResolver
}

func NewCollectMap() *CollectMap {
	return &CollectMap{
		cache: sync.Map{},
	}
}

func (cm *CollectMap) AddNewCollect(clusterName string, baseLogPath string) {
	_, ok := cm.cache.Load(clusterName)
	if !ok {
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
