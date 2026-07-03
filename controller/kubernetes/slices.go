package kubernetes

import (
	discoveryv1 "k8s.io/api/discovery/v1"
)

func AddSlice(owner string, endpoints []discoveryv1.Endpoint) {
	clientset := GetClientset()
	var addresses []string
	for _, endpoint := range endpoints {
		addresses = append(addresses, endpoint.Addresses[0])
	}
	UpdateAgentCM(clientset, GetHeadlessServiceKey(owner), "add", addresses...)
}
