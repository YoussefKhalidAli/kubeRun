package main

import (
	"github.com/ti-mo/netfilter"
)

var KubeRunController string = "localhost:4444"

func main() {
	// updated := os.Getenv.("UPDATED")
	updated := false

	eventGroups := []netfilter.NetlinkGroup{
		netfilter.GroupCTNew,
	}

	if updated {
		eventGroups = append(eventGroups, netfilter.GroupCTUpdate)
	}

	c := KernelListener(eventGroups)
	defer c.Close()
	Filter()
}
