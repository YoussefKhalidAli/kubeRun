package informer

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"kuberun.com/controller/store"
)

func podInformer(factory informers.SharedInformerFactory) {
	podInformer := factory.Core().V1().Pods().Informer()

	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj any) {
			setKubeRunPodIp(obj)
		},
		UpdateFunc: func(_, obj any) {
			setKubeRunPodIp(obj)
		},
	})
}

func setKubeRunPodIp(obj any) {
	pod, ok := obj.(*corev1.Pod)
	if !ok || pod == nil {
		return
	}
	store.KubeRunPodIp = pod.Status.PodIP
	println("KubeRunPodIp", store.KubeRunPodIp)
}
