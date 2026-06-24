package main

import (
	"fmt"
	"slices"
)

func Filter() {
	println("Listening for connections")

	ips := [2]string{"10.101.138.168", "10.106.160.35"}
	for event := range EventChan {
		ip := event.Flow.TupleOrig.IP.DestinationAddress.Unmap().String()
		if slices.Contains(ips[:], ip) {
			fmt.Println(ip)

		}
	}
}
