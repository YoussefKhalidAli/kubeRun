package kubernetes

import (
	"time"

	"kuberun.com/controller/store"
)

func SyncLoop() {
	ticker := time.NewTicker(store.SyncTime)
	defer ticker.Stop()

	for range ticker.C {
		sync()
	}
}

func sync() {
	for _, targetVal := range store.Targets {
		if time.Now().After(targetVal.LastAccessed.Add(store.SyncTime)) && targetVal.Status == "Awake" {
			targetVal.Status = "Sleeping"
			ScaleResource(targetVal, 0)
		}
	}
}
