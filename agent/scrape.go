package main

import (
	"errors"
	"syscall"

	"github.com/ti-mo/conntrack"
	"github.com/ti-mo/netfilter"
	"kuberun.com/agent/utils"
)

var EventChan chan conntrack.Event

// This function opens a connection with the kernel and scrapes all new and updated connections.
// All captured data is send to "EventChan" channel that is then filtered
func KernelListener(eventGroups []netfilter.NetlinkGroup) *conntrack.Conn {
	c, err := conntrack.Dial(nil)
	if err != nil {
		if errors.Is(err, syscall.EPERM) {
			utils.HandelError(err, "KRA0403H", "conntrack Dial failed due to insufficient permissions (CAP_NET_ADMIN required)")
		} else if errors.Is(err, syscall.EPROTONOSUPPORT) {
			utils.HandelError(err, "KRA0404H", "the nf_conntrack module isn't loaded on the host node")
		} else {
			utils.HandelError(err, "KRA9022H", "unexpected conntrack Dial failure")
		}
		return nil
	}

	EventChan = make(chan conntrack.Event)

	go func() {
		if _, err := c.Listen(EventChan, 2, eventGroups); err != nil {
			utils.HandelError(err, "KRA9014H", "failed to listen to Netlink group events")
		}
	}()

	return c
}
