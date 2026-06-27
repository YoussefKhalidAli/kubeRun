package main

import (
	"kuberun.com/controller/kubernetes"
	"kuberun.com/controller/server"
	"kuberun.com/controller/utils"
)

func main() {
	utils.Targets = make(map[string]*utils.TargetDto)

	go server.Alert()

	kubernetes.Kubernetes()
}
