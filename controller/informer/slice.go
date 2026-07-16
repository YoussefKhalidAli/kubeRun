package informer

import (
	"strings"
	"time"

	discoveryv1 "k8s.io/api/discovery/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"kuberun.com/controller/slice"
	"kuberun.com/controller/store"
	"kuberun.com/controller/utils"
)

func endpointSlicesInformer(factory informers.SharedInformerFactory) {

	endpointSlicesInformer := factory.Discovery().V1().EndpointSlices().Informer()

	endpointSlicesInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(_ any, obj any) {
			println("Got endpoint")
			eSlice, ok := obj.(*discoveryv1.EndpointSlice)
			if !ok || eSlice == nil {
				return
			}
			if len(eSlice.Endpoints) > 0 && len(eSlice.ObjectMeta.OwnerReferences) > 0 {
				selector := eSlice.Endpoints[0].Hostname
				owner := eSlice.ObjectMeta.OwnerReferences[0].Name

				shouldAddSlice := selector != nil &&
					!strings.Contains(*selector, "kuberun-controller") &&
					waitForTargetCreation(utils.GetHeadlessServiceKey(owner))

				if shouldAddSlice {
					endpoints := eSlice.Endpoints
					slice.AddSlice(owner, endpoints)
				}
			}
		},
	})
}

func waitForTargetCreation(key string) bool {

	ticker := time.NewTicker(time.Second / 2)
	defer ticker.Stop()
	timeout := time.After(5 * time.Second)

	for {
		select {
		case <-ticker.C:
			if store.Targets[key] != nil {
				return true
			}
		case <-timeout:
			return false
		}
	}
}
