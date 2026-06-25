package main

import (
	"github.com/ti-mo/netfilter"
)

func main() {
	// kubeRunNamespace := os.Getenv("NAMESPACE")
	// updated := os.Getenv.("UPDATED")
	updated := true

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
