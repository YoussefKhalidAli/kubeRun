package main

import (
	"fmt"

	"github.com/ti-mo/netfilter"
)

func main() {
	watcher := LoadVariables()
	defer watcher.Close()

	fmt.Println(Config)
	eventGroups := []netfilter.NetlinkGroup{
		netfilter.GroupCTNew,
	}

	if Config.Update {
		eventGroups = append(eventGroups, netfilter.GroupCTUpdate)
	}

	c := KernelListener(eventGroups)
	defer c.Close()
	Filter()
}
