package main

import (
	"kuberun.com/controller/kubernetes"
	"kuberun.com/controller/utils"
)

func main() {
	utils.Targets = make(map[string]*utils.TargetDto)

	go Alert()

	kubernetes.Kubernetes()
}
