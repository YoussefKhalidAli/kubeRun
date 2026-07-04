package informer

import (
	"strings"

	discoveryv1 "k8s.io/api/discovery/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"kuberun.com/controller/slice"
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
				if selector != nil && !strings.Contains(*selector, "kuberun-controller") {
					owner := eSlice.ObjectMeta.OwnerReferences[0].Name
					endpoints := eSlice.Endpoints
					slice.AddSlice(owner, endpoints)
				}
			}
		},
	})
}
