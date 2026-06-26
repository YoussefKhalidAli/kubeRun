package main

import (
	"time"

	"kuberun.com/controller/kubernetes"
	"kuberun.com/controller/utils"
)

func main() {
	utils.Targets = make(map[string]*utils.TargetDto)
	utils.Targets["10.110.91.113"] = &utils.TargetDto{
		LastAccessed: time.Now().Add(-24 * time.Hour),
		ResourceName: "auth-service-pod",
		Namespace:    "production",
		Resource:     "pods",
	}

	utils.Targets["10.106.160.35"] = &utils.TargetDto{
		LastAccessed: time.Now().Add(-24 * time.Hour),
		ResourceName: "payment-gateway-svc",
		Namespace:    "staging",
		Resource:     "services",
	}
	go Alert()
	kubernetes.Kubernetes()
}
