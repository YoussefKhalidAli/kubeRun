package main

import (
	"slices"
)

func Filter() {
	println("Listening for connections")

	for event := range EventChan {
		ip := event.Flow.TupleOrig.IP.DestinationAddress.Unmap().String()
		if slices.Contains(Config.Ips[:], ip) {
			println(ip)
			Alert(ip)
		}
	}
}
