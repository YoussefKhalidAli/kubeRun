package main

import (
	"errors"
	"syscall"

	"github.com/ti-mo/conntrack"
	"github.com/ti-mo/netfilter"
)

var EventChan chan conntrack.Event

// This function opens a connection with the kernel and scrapes all new and updated connections.
// All captured data is send to "EventChan" channel that is then filtered
func KernelListener() *conntrack.Conn {
	c, err := conntrack.Dial(nil)
	if err != nil {
		if errors.Is(err, syscall.EPERM) {
			HandelError(err, "KR0403", "_")
		} else if errors.Is(err, syscall.EPROTONOSUPPORT) {
			HandelError(err, "KR0404", "the nf_conntrack module isn't loaded on the host node")
		}
	}

	EventChan = make(chan conntrack.Event)
	eventGroups := []netfilter.NetlinkGroup{
		netfilter.GroupCTNew,
		netfilter.GroupCTUpdate,
	}

	go func() {
		if _, err := c.Listen(EventChan, 2, eventGroups); err != nil {
			HandelError(err, "KR0403", "_")
		}
	}()

	return c
}
