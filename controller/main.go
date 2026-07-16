package main

import (
	"kuberun.com/controller/concile"
	"kuberun.com/controller/informer"
	"kuberun.com/controller/server"
	"kuberun.com/controller/store"
	"kuberun.com/controller/utils"
)

var logger = utils.Logger

func main() {
	store.Targets = make(map[string]*store.TargetDto)

	server.MarkerMux.Lock()
	server.Reserve(4444)
	server.MarkerMux.Unlock()

	go Alert()

	go concile.SyncLoop()
	informer.Connect()
}
