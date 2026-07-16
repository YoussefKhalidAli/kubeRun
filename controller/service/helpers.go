package service

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"kuberun.com/controller/store"
	"kuberun.com/controller/utils"
)

func CreateshadowService(service *corev1.Service, clientset *kubernetes.Clientset, key string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shadowService := service
	target := store.Targets[key]

	target.Mux.Lock()
	servers := target.Servers
	name := target.ServiceName
	target.Mux.Unlock()

	for index, _ := range shadowService.Spec.Ports {
		shadowService.Spec.Ports[index].TargetPort = intstr.FromInt(servers[index].SwitchPort)
	}

	shadowService.ObjectMeta.Name = utils.GetShadowName(name)
	shadowService.ObjectMeta.Namespace = store.KubeRunNamespace
	shadowService.Spec.Type = corev1.ServiceTypeClusterIP
	shadowService.Spec.ClusterIP = ""
	shadowService.Spec.ClusterIPs = nil

	delete(shadowService.ObjectMeta.Labels, "kuberun/run")
	shadowService.ObjectMeta.Labels["kuberun/operator"] = "shadow"
	shadowService.Spec.Selector = map[string]string{
		"kuberun/operator": "controller",
	}

	shadowService.ResourceVersion = ""
	shadowService.UID = ""
	shadowService.CreationTimestamp = metav1.Time{}
	shadowService.Generation = 0
	shadowService.Status = corev1.ServiceStatus{}

	_, err := clientset.CoreV1().Services(store.KubeRunNamespace).Create(ctx, shadowService, metav1.CreateOptions{})

	if err != nil {
		utils.HandelError(err, "KRC1445", "Failed to create shadow service")
	}
}

func DeleteShadowService(clientset *kubernetes.Clientset, shadow string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := clientset.CoreV1().Services(store.KubeRunNamespace).Delete(ctx, shadow, metav1.DeleteOptions{})

	if err != nil {
		utils.HandelError(err, "KRC1445", "Failed to delete shadow service")
	}
}
