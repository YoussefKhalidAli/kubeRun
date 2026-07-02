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

		targetVal.Mux.Lock()
		resource := targetVal.Resource
		resourceName := targetVal.ResourceName
		isResourcePresent := resource != "" && resourceName != ""
		shouldSleep := time.Now().After(targetVal.LastAccessed.Add(store.SyncTime)) && targetVal.Status == "Awake"
		targetVal.Mux.Unlock()

		if !isResourcePresent {
			isResourcePresent = checkResource(targetVal, index)
		} else if shouldSleep {
			targetVal.Mux.Lock()
			targetVal.Status = "Sleeping"
			targetVal.Mux.Unlock()

			ScaleResource(targetVal, 0)
		}
	}
}

func checkResource(target *store.TargetDto, index string) bool {

	target.Mux.Lock()
	namespace := target.Namespace
	selectorMap := target.SelectorMap
	target.Mux.Unlock()

	resourceName, resource := FindResource(GetClientset(), selectorMap, namespace, index)
	if resourceName == "kuberun-controller" || resource == "DaemonSet" {
		println("Found unmanagable resource. Skipping")
		RemoveService(clientset, index)
		return false
	} else if resourceName == "" && resource == "" {
		return false
	}

	target.Mux.Lock()
	target.Resource = resource
	target.ResourceName = resourceName
	target.Mux.Unlock()

	return true
}
