package bootstrap

import (
	"EmptyClassroom/config"
	"EmptyClassroom/logs"
	"EmptyClassroom/snapshot"
	"sync"
)

var initOnce sync.Once

func Init(isMain bool) {
	initOnce.Do(func() {
		logs.Init(isMain)
		config.InitConfig()
	})
}

func NewSnapshotStore() snapshot.Store {
	return snapshot.NewDefaultStore()
}
