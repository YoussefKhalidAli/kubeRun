package service

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"kuberun.com/controller/store"
)

func CreateshadowService(service *corev1.Service, clientset *kubernetes.Clientset, key string) {
	ctx := context.Background()

	shadowService := service
	target := store.Targets[key]

	target.Mux.Lock()
	servers := target.Servers
	target.Mux.Unlock()

	for index, _ := range shadowService.Spec.Ports {
		shadowService.Spec.Ports[index].TargetPort = intstr.FromInt(servers[index].SwitchPort)
	}

	shadowService.Spec.Selector = map[string]string{
		"kuberun/operator": "controller",
	}
	shadowService.ObjectMeta.Namespace = store.KubeRunNamespace
	delete(shadowService.ObjectMeta.Labels, "kuberun/run")
	shadowService.ObjectMeta.Labels["kuberun/operator"] = "shadow"

	CreatedShadowService, _ := clientset.CoreV1().Services(store.KubeRunNamespace).Create(ctx, shadowService, metav1.CreateOptions{})

	target.Mux.Lock()
	target.Shadow = CreatedShadowService.ResourceVersion
	target.Mux.Unlock()

}
