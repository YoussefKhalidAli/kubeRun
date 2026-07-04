package targets

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"kuberun.com/controller/agent"
	"kuberun.com/controller/scale"
	"kuberun.com/controller/server"
	"kuberun.com/controller/store"
)

func CreateTarget(key string, svc corev1.ServiceSpec, metadata metav1.ObjectMeta, resourceName string, resource string) {

	store.Targets[key] = &store.TargetDto{
		LastAccessed: time.Now(),
		ResourceName: resourceName,
		Namespace:    metadata.Namespace,
		Resource:     resource,
		ServiceName:  metadata.Name,
		Server:       server.New(),
		Status:       "Awake",
		ServicePorts: MapServicePorts(svc.Ports),
		SelectorMap:  svc.Selector,
	}

	target := store.Targets[key]
	target.Server.ScaleUp = func() {
		target.Mux.Lock()
		shouldWake := target.Status != "Awake" && target.Status != "Waking"
		if shouldWake {
			go scale.ScaleResource(key, 1)
			target.Status = "Waking"
		}
		target.Mux.Unlock()
	}
}

func DeleteTarget(clientset *kubernetes.Clientset, key string) {
	target := store.Targets[key]
	if target.Status == "Asleep" {
		target.Server.Stop()
		target.Server.Signal.Unlock()
	}

	err := agent.UpdateAgentCM(clientset, key, "delete")
	if err != nil {
		fmt.Printf("Error occurred while updating agent config map: _%v_", err)
	}
	delete(store.Targets, key)
}

func MapServicePorts(portsMap []corev1.ServicePort) *[]int {
	targetPortsMap := make([]int, len(portsMap))
	for index, port := range portsMap {
		targetPortsMap[index] = port.TargetPort.IntValue()
	}
	return &targetPortsMap
}
