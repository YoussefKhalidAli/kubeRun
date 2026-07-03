package kubernetes

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"kuberun.com/controller/store"
)

func AddService(svc corev1.ServiceSpec, metadata metav1.ObjectMeta, clientset *kubernetes.Clientset) {
	resourceName, resource := FindResource(clientset, svc.Selector, metadata.Namespace, svc.ClusterIP)
	if resourceName == "kuberun-controller" || resource == "DaemonSet" {
		println("Found unmanagable resource. Skipping")
		return
	}

	key := svc.ClusterIP
	if key == "None" {
		key = GetHeadlessServiceKey(metadata.Name)
	} else {
		UpdateAgents(key)
		UpdateAgentCM(clientset, key, "add")
	}

	CreateTarget(key, svc, metadata, resourceName, resource)

	store.PrintTargets()
}

func UpdateService(clussterIp string, service *corev1.Service, old *corev1.Service, clientset *kubernetes.Clientset) {
	var key string

	if clussterIp == "None" {
		key = GetHeadlessServiceKey(service.ObjectMeta.Name)
	} else {
		key = clussterIp
	}
	target := store.Targets[key]

	if target == nil {
		DeleteTarget(clientset, key)
		AddService(service.Spec, service.ObjectMeta, clientset)
		return
	}

	target.Mux.Lock()
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

func GetHeadlessServiceKey(name string) string {
	return fmt.Sprintf("svc-%v", name)
}
