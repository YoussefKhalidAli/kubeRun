package main

import (
	"kuberun.com/controller/concile"
	"kuberun.com/controller/informer"
	"kuberun.com/controller/store"
)

func main() {
	store.Targets = make(map[string]*store.TargetDto)

	go Alert()

	go concile.SyncLoop()
	informer.Connect()
}
