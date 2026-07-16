package concile

import (
	"time"

	"kuberun.com/controller/client"
	"kuberun.com/controller/resource"
	"kuberun.com/controller/scale"
	"kuberun.com/controller/store"
	"kuberun.com/controller/targets"
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
			targetVal.Mux.Unlock()
		} else if shouldSleep {
			targetVal.Status = "Sleeping"
			targetVal.Mux.Unlock()

			scale.ScaleResource(index, 0)
		} else {
			targetVal.Mux.Unlock()
		}
		store.PrintTargets()
	}
}

func checkResource(target *store.TargetDto, index string) {
	clientset := client.GetClientset()

	resourceName, resource := resource.FindResource(clientset, target.SelectorMap, target.Namespace, index)
	if resourceName == "kuberun-controller" || resource == "DaemonSet" {
		println("Found unmanagable resource. Skipping")
		targets.DeleteTarget(clientset, index)
		return
	} else if resourceName == "" && resource == "" {
		return
	}

	target.Resource = resource
	target.ResourceName = resourceName
}
