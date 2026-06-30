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
		targetVal.Mux.Lock()
		shouldSleep := time.Now().After(targetVal.LastAccessed.Add(store.SyncTime)) && targetVal.Status == "Awake"
		targetVal.Mux.Unlock()
		if shouldSleep {
			targetVal.Mux.Lock()
			targetVal.Status = "Sleeping"
			targetVal.Mux.Unlock()
			ScaleResource(targetVal, 0)
		}
	}
}
