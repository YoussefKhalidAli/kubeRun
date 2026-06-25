package main

import (
	"fmt"
	"slices"
)

func Filter() {
	println("Listening for connections")

	ips := [2]string{"10.101.138.168", "10.110.91.113"}
	for event := range EventChan {
		ip := event.Flow.TupleOrig.IP.DestinationAddress.Unmap().String()
		if slices.Contains(ips[:], ip) {
			fmt.Println(ip)
		}
	}
}
