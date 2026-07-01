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
		isAbsent := false

		targetVal.Mux.Lock()
		resource := targetVal.Resource
		resourceName := targetVal.ResourceName
		namespace := targetVal.Namespace
		selectorMap := targetVal.SelectorMap
		shouldSleep := time.Now().After(targetVal.LastAccessed.Add(store.SyncTime)) && targetVal.Status == "Awake"
		targetVal.Mux.Unlock()

		if resource == "" && resourceName == "" {
			resourceName, resource := FindResource(GetClientset(), selectorMap, namespace, index)
			if resourceName == "kuberun-controller" || resource == "DaemonSet" {
				println("Found unmanagable resource. Skipping")
				RemoveService(clientset, index)
				return
			} else if resourceName == "" && resource == "" {
				continue
			}
			isAbsent = true
		}

		if shouldSleep {
			targetVal.Mux.Lock()
			targetVal.Status = "Sleeping"

			if isAbsent {
				targetVal.Resource = resource
				targetVal.ResourceName = resourceName
			}

			targetVal.Mux.Unlock()
			ScaleResource(targetVal, 0)
		}
	}
}
