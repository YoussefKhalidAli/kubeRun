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

func UpdateService(clussterIp string, service *corev1.Service) {
	target := store.Targets[clussterIp]
	target.Mux.Lock()
	println("new", &service)
	if target.UpdateMarker == service.ResourceVersion {
		target.Mux.Unlock()
		return
	}

	targetStatus := target.Status
	target.SelectorMap = service.Spec.Selector
	target.ServicePorts = MapServicePorts(service.Spec.Ports)
	target.Mux.Unlock()
	if targetStatus == "Asleep" && service.Spec.Selector["KubeRun"] != "Controller" {
		PatchService(target, 0)
	}
}

func DeleteService(clusterIP string, clientset kubernetes.Interface) {

	UpdateAgentCM(clientset, clusterIP, "delete")
	store.Targets[clusterIP].Server.Stop()
	delete(store.Targets, clusterIP)

	store.PrintTargets()
}
