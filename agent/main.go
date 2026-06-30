package main

import (
	"fmt"

	"github.com/ti-mo/netfilter"
	"kuberun.com/agent/store"
)

func main() {
	watcher := store.LoadVariables()
	defer watcher.Close()

	fmt.Println(store.Config)
	eventGroups := []netfilter.NetlinkGroup{
		netfilter.GroupCTNew,
	}

	if store.Config.Update {
		eventGroups = append(eventGroups, netfilter.GroupCTUpdate)
	}

	c := KernelListener(eventGroups)
	defer c.Close()
	Filter()
}
