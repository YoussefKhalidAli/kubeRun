package service

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"kuberun.com/controller/agent"
	"kuberun.com/controller/resource"
	"kuberun.com/controller/scale"
	"kuberun.com/controller/store"
	"kuberun.com/controller/targets"
	"kuberun.com/controller/utils"
)

func AddService(service *corev1.Service, clientset *kubernetes.Clientset) {
	svc, metadata := service.Spec, service.ObjectMeta
	resourceName, resourceKind := resource.FindResource(clientset, svc.Selector, metadata.Namespace, svc.ClusterIP)

	if resourceName == "kuberun-controller" || resourceKind == "DaemonSet" {
		println("Found unmanagable resource. Skipping")
		return
	}

	key := svc.ClusterIP
	var serviceType string

	if key == "None" {
		key = utils.GetHeadlessServiceKey(metadata.Name)
		serviceType = "Headless"
	} else {
		serviceType = string(corev1.ServiceTypeClusterIP)
		agent.UpdateAgents(key)
		agent.UpdateAgentCM(clientset, key, "add")
	}

	targets.CreateTarget(key, svc, metadata, resourceName, resourceKind)
	CreateshadowService(service, clientset, key)

	labelKeys := [2]string{"key", "type"}
	labelValues := [2]string{key, serviceType}
	err := labelService(clientset, metadata.Name, metadata.Namespace, key, labelKeys[:], labelValues[:])
	if err != nil {
		utils.HandelError(err, "KRC1452M", fmt.Sprintf("Couldn't update service %v after retrying", metadata.Name))
		return
	}

	store.PrintTargets()
}

func UpdateService(clussterIp string, service *corev1.Service, old *corev1.Service, clientset *kubernetes.Clientset) {
	key := ""
	if service.Labels != nil {
		key = service.Labels["kuberun/key"]
	}
	target, ok := store.Targets[key]

	if !ok || target == nil {
		targets.DeleteTarget(clientset, key)
		AddService(service, clientset)
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
	if targetStatus == "Asleep" && service.Spec.Selector["kuberun/operator"] != "controller" {
		scale.PatchService(key, 0)
	}
}
