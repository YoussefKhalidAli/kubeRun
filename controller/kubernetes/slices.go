package kubernetes

import (
	"fmt"

	discoveryv1 "k8s.io/api/discovery/v1"
)

func AddSlice(owner string, endpoints []discoveryv1.Endpoint) {
	clientset := GetClientset()
	fmt.Println("endpoints", endpoints)
	var addresses []string
	for _, endpoint := range endpoints {
		fmt.Println("endpoint", endpoint)
		addresses = append(addresses, endpoint.Addresses[0])
	}
	fmt.Println("addresses", addresses)
	UpdateAgentCM(clientset, "None", "add", addresses...)
}
