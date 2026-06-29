package kubernetes

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"kuberun.com/controller/store"
)

func AddService(svc corev1.ServiceSpec, metadata metav1.ObjectMeta, clientset kubernetes.Interface) {
	resourceName, resource := FindResource(clientset, svc.Selector, metadata.Namespace)
	if resourceName == "kuberun-controller" || resource == "DaemonSet" {
		println("Found unmanagable resource. Skipping")
		return
	}

	CreateService(svc, metadata, resourceName, resource)

	UpdateAgentCM(clientset, svc.ClusterIP, "add")

	store.PrintTargets()
}

func UpdateService(clussterIp string, service corev1.ServiceSpec) {
	target := store.Targets[clussterIp]
	target.SelectorMap = service.Selector
	target.ServicePorts = MapServicePorts(service.Ports)
	if target.Status == "Asleep" {
		PatchService(target, 0)
	}
}

func DeleteService(clusterIP string, clientset kubernetes.Interface) {

	UpdateAgentCM(clientset, clusterIP, "delete")
	store.Targets[clusterIP].Server.Stop()
	delete(store.Targets, clusterIP)

	store.PrintTargets()
}
