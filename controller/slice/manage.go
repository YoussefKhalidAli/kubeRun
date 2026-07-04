package slice

import (
	discoveryv1 "k8s.io/api/discovery/v1"
	"kuberun.com/controller/agent"
	"kuberun.com/controller/client"
	"kuberun.com/controller/service"
)

func AddSlice(owner string, endpoints []discoveryv1.Endpoint) {
	clientset := client.GetClientset()
	var addresses []string
	for _, endpoint := range endpoints {
		addresses = append(addresses, endpoint.Addresses[0])
	}
	agent.UpdateAgentCM(clientset,
		service.GetHeadlessServiceKey(owner), "add", addresses...)
}
