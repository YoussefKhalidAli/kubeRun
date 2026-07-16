package informer

import (
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"kuberun.com/controller/client"
	"kuberun.com/controller/service"
	"kuberun.com/controller/store"
	"kuberun.com/controller/targets"
	"kuberun.com/controller/utils"
)

func serviceInformer(factory informers.SharedInformerFactory) {
	clientset := client.GetClientset()

	serviceInformer := factory.Core().V1().Services().Informer()

	serviceInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj any) {
			svc := obj.(*corev1.Service)
			key := svc.ObjectMeta.Labels["kuberun/key"]
			if store.Targets[key] == nil {
				service.AddService(svc, clientset)
			}
		},
		DeleteFunc: func(obj any) {
			svc, ok := obj.(*corev1.Service)
			if !ok {
				panic("couldn't convert object to service")
			}
			key := svc.Labels["kuberun/key"]

			target := store.Targets[key]
			target.Mux.Lock()
			store.PrintTargets()
			if !strings.Contains(string(target.Status), "ing") {
				targets.DeleteTarget(clientset, key)
				service.DeleteShadowService(clientset, utils.GetShadowName(target.ServiceName))
			}
			target.Mux.Unlock()
		},
		UpdateFunc: func(old any, obj any) {
			svc := obj.(*corev1.Service)
			service.UpdateService(svc.Spec.ClusterIP, svc, old.(*corev1.Service), clientset)
		},
	})
}
