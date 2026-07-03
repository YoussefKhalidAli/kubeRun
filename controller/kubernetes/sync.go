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
	for index, targetVal := range store.Targets {
		println("Syncing: ", index)

		targetVal.Mux.Lock()
		isResourcePresent := targetVal.Resource != "" && targetVal.ResourceName != ""
		shouldSleep := time.Now().After(targetVal.LastAccessed.Add(store.SyncTime)) && targetVal.Status == "Awake"

		if !isResourcePresent {
			checkResource(targetVal, index)
		} else if shouldSleep {
			targetVal.Status = "Sleeping"
			targetVal.Mux.Unlock()

			ScaleResource(targetVal, 0)
		} else {
			targetVal.Mux.Unlock()
		}
	}
}

func checkResource(target *store.TargetDto, index string) {

	resourceName, resource := FindResource(GetClientset(), target.SelectorMap, target.Namespace, index)
	if resourceName == "kuberun-controller" || resource == "DaemonSet" {
		println("Found unmanagable resource. Skipping")
		DeleteTarget(clientset, index)
		return
	} else if resourceName == "" && resource == "" {
		return
	}

	target.Resource = resource
	target.ResourceName = resourceName
	target.Mux.Unlock()

}
