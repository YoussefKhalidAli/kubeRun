package service

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"kuberun.com/controller/agent"
	"kuberun.com/controller/resource"
	"kuberun.com/controller/scale"
	"kuberun.com/controller/store"
	"kuberun.com/controller/targets"
)

func AddService(svc corev1.ServiceSpec, metadata metav1.ObjectMeta, clientset *kubernetes.Clientset) {
	resourceName, resourceKind := resource.FindResource(clientset, svc.Selector, metadata.Namespace, svc.ClusterIP)
	if resourceName == "kuberun-controller" || resourceKind == "DaemonSet" {
		println("Found unmanagable resource. Skipping")
		return
	}

	key := svc.ClusterIP
	if key == "None" {
		key = GetHeadlessServiceKey(metadata.Name)
	} else {
		agent.UpdateAgents(key)
		agent.UpdateAgentCM(clientset, key, "add")
	}
	targets.CreateTarget(key, svc, metadata, resourceName, resourceKind)
	labelService(clientset, metadata.Name, metadata.Namespace, key)

	store.PrintTargets()
}

func UpdateService(clussterIp string, service *corev1.Service, old *corev1.Service, clientset *kubernetes.Clientset) {
	key := service.Labels["kuberun/key"]
	target := store.Targets[key]

	if target == nil {
		targets.DeleteTarget(clientset, key)
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
	target.ServicePorts = targets.MapServicePorts(service.Spec.Ports)
	target.Mux.Unlock()
	if targetStatus == "Asleep" && service.Spec.Selector["KubeRun"] != "Controller" {
		scale.PatchService(key, 0)
	}
}

func GetHeadlessServiceKey(name string) string {
	return fmt.Sprintf("svc-%v", name)
}
