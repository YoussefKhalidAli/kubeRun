package kubernetes

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"kuberun.com/controller/server"
	"kuberun.com/controller/store"
)

func CreateService(svc corev1.ServiceSpec, metadata metav1.ObjectMeta, resourceName string, resource string) {

	store.Targets[svc.ClusterIP] = &store.TargetDto{
		LastAccessed: time.Now(),
		ResourceName: resourceName,
		Namespace:    metadata.Namespace,
		Resource:     resource,
		ServiceName:  metadata.Name,
		Server:       server.New(),
		IsSleep:      false,
		ServicePorts: MapServicePorts(svc.Ports),
		SelectorMap:  svc.Selector,
	}

}

func FilterAnnotations(anns map[string]string) bool {
	run := false
	if anns[store.RunAnnotation] == "true" {
		run = true
	}
	return run
}

func MapServicePorts(portsMap []corev1.ServicePort) *[]int {
	targetPortsMap := make([]int, len(portsMap))
	for index, port := range portsMap {
		targetPortsMap[index] = port.TargetPort.IntValue()
	}
	return &targetPortsMap
}
