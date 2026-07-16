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
			svc, ok := obj.(*corev1.Service)
			if !ok || svc == nil {
				return
			}
			key := ""
			if svc.ObjectMeta.Labels != nil {
				key = svc.ObjectMeta.Labels["kuberun/key"]
			}
			if key != "" && store.Targets[key] == nil {
				service.AddService(svc, clientset)
			}
		},
		DeleteFunc: func(obj any) {
			svc, ok := obj.(*corev1.Service)
			if !ok || svc == nil {
				return
			}
			key := ""
			if svc.Labels != nil {
				key = svc.Labels["kuberun/key"]
			}

			target, ok := store.Targets[key]
			if !ok || target == nil {
				return
			}
			target.Mux.Lock()
			store.PrintTargets()
			if !strings.Contains(string(target.Status), "ing") {
				target.Mux.Unlock()
				targets.DeleteTarget(clientset, key)
				service.DeleteShadowService(clientset, utils.GetShadowName(target.ServiceName))
			} else {
				target.Mux.Unlock()
			}
		},
		UpdateFunc: func(old any, obj any) {
			svc, ok := obj.(*corev1.Service)
			if !ok || svc == nil {
				return
			}
			oldSvc, ok := old.(*corev1.Service)
			if !ok || oldSvc == nil {
				return
			}
			service.UpdateService(svc.Spec.ClusterIP, svc, oldSvc, clientset)
		},
	})
}
