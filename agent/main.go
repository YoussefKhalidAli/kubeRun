package main

import (
	"github.com/ti-mo/netfilter"
	"kuberun.com/agent/server"
	"kuberun.com/agent/store"
	"kuberun.com/agent/utils"
)

var logger = utils.Logger

func main() {
	go server.Updates()

	watcher := store.LoadVariables()
	defer watcher.Close()

	logger.Info("loaded agent config", "config", store.Config)
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
