package informer

import (
	"strings"
	"time"

	discoveryv1 "k8s.io/api/discovery/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"kuberun.com/controller/service"
	"kuberun.com/controller/slice"
	"kuberun.com/controller/store"
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
			if len(eSlice.Endpoints) > 0 {
				selector := eSlice.Endpoints[0].Hostname
				owner := eSlice.ObjectMeta.OwnerReferences[0].Name

				shouldAddSlice := selector != nil &&
					!strings.Contains(*selector, "kuberun-controller") &&
					store.Targets[service.GetHeadlessServiceKey(owner)] != nil

				if shouldAddSlice {
					time.Sleep(2 * time.Second)
					endpoints := eSlice.Endpoints
					slice.AddSlice(owner, endpoints)
				}
			}
		},
	})
}
