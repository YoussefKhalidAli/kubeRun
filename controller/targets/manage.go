package targets

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"kuberun.com/controller/agent"
	"kuberun.com/controller/scale"
	"kuberun.com/controller/server"
	"kuberun.com/controller/store"
	"kuberun.com/controller/utils"
)

func CreateTarget(key string, svc corev1.ServiceSpec, metadata metav1.ObjectMeta, resourceName string, resource string) {

	store.Targets[key] = &store.TargetDto{
		LastAccessed: time.Now(),
		ResourceName: resourceName,
		Namespace:    metadata.Namespace,
		Resource:     resource,
		ServiceName:  metadata.Name,
		Status:       "Awake",
		ServicePorts: MapServicePorts(svc.Ports),
		SelectorMap:  svc.Selector,
	}

	target := store.Targets[key]
	target.Servers = make([]*server.Switch, len(*target.ServicePorts))

	scaleUp := func() {
		target.Mux.Lock()
		shouldWake := target.Status != "Awake" && target.Status != "Waking"
		if shouldWake {

			for _, port := range target.Servers {
				switchPort := strconv.Itoa(port.SwitchPort)
				go killSwitch(switchPort)
			}
			target.LastAccessed = time.Now()
			target.Status = "Waking"
			target.Mux.Unlock()
			go scale.ScaleResource(key, 1)
		} else {
			target.Mux.Unlock()
		}
	}

	for index, port := range *target.ServicePorts {
		target.Servers[index] = server.New(strconv.Itoa(port))
		target.Servers[index].ScaleUp = scaleUp
	}

}

func DeleteTarget(clientset *kubernetes.Clientset, key string) {
	target, ok := store.Targets[key]
	if !ok || target == nil {
		return
	}
	if target.Status == "Asleep" {
		for _, server := range target.Servers {
			if server != nil {
				server.Kill()
			}
		}
	}

	agent.UpdateAgentCM(clientset, key, "delete")
	delete(store.Targets, key)
	store.PrintTargets()
}

func MapServicePorts(portsMap []corev1.ServicePort) *[]int {
	targetPortsMap := make([]int, len(portsMap))
	for index, port := range portsMap {
		targetPortsMap[index] = port.TargetPort.IntValue()
	}
	return &targetPortsMap
}

func killSwitch(port string) {
	resp, err := http.Post("http://"+store.KubeRunPodIp+":"+port, "text/plain", strings.NewReader("Go to Sleep"))
	if err != nil {
		utils.HandelError(err, "KRC1451L", "Failed to kill switch on port "+port)
		return
	}
	resp.Body.Close()
}
