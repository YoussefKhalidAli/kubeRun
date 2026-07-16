package main

import (
	"slices"

	"kuberun.com/agent/server"
	"kuberun.com/agent/store"
)

func Filter() {
	logger.Info("listening for connections")

	for event := range EventChan {
		ip := event.Flow.TupleOrig.IP.DestinationAddress.Unmap().String()
		if slices.Contains(store.Config.Ips[:], ip) {
			logger.Info("matched tracked ip", "ip", ip)
			server.Alert(ip)
		}
	}
}
