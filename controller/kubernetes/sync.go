package kubernetes

import (
	"time"

	"kuberun.com/controller/utils"
)

func SyncLoop() {
	ticker := time.NewTicker(utils.SyncTime)
	defer ticker.Stop()

	for range ticker.C {
		sync()
	}
}

func sync() {
	for _, targetVal := range utils.Targets {
		if time.Now().After(targetVal.LastAccessed.Add(utils.SyncTime)) {
			ScaleResource(targetVal, 0)
		}
	}
}
