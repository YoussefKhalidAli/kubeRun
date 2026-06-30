package main

import (
	"kuberun.com/controller/kubernetes"
	"kuberun.com/controller/store"
)

func main() {
	store.Targets = make(map[string]*store.TargetDto)

	go Alert()

	kubernetes.Kubernetes()
}
