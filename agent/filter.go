package main

import (
	"slices"

	"kuberun.com/agent/server"
	"kuberun.com/agent/store"
)

func Filter() {
	println("Listening for connections")

	for event := range EventChan {
		ip := event.Flow.TupleOrig.IP.DestinationAddress.Unmap().String()
		if slices.Contains(store.Config.Ips[:], ip) {
			println(ip)
			server.Alert(ip)
		}
	}
}
