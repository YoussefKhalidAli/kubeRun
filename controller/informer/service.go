package informer

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"kuberun.com/controller/client"
	"kuberun.com/controller/service"
	"kuberun.com/controller/store"
	"kuberun.com/controller/targets"
)

func serviceInformer(factory informers.SharedInformerFactory) {
	clientset := client.GetClientset()

	serviceInformer := factory.Core().V1().Services().Informer()

	serviceInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj any) {
			svc := obj.(*corev1.Service)
			service.AddService(svc.Spec, svc.ObjectMeta, clientset)
		},
		DeleteFunc: func(obj any) {
			svc, ok := obj.(*corev1.Service)
			if !ok {
				panic("couldn't convert object to service")
			}
			key := svc.Spec.ClusterIP
			if key == "None" {
				key = service.GetHeadlessServiceKey(svc.ObjectMeta.Name)
			}
			targets.DeleteTarget(clientset, key)
			store.PrintTargets()
		},
		UpdateFunc: func(old any, obj any) {
			svc := obj.(*corev1.Service)
			service.UpdateService(svc.Spec.ClusterIP, svc, old.(*corev1.Service), clientset)
		},
	})
}
