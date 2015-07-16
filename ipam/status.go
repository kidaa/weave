package ipam

import (
	"github.com/weaveworks/weave/ipam/paxos"
	"github.com/weaveworks/weave/net/address"
)

func IPAMStatus(allocator *Allocator, defaultSubnet address.CIDR) interface{} {
	resultChan := make(chan interface{})
	allocator.actionChan <- func() {
		resultChan <- struct {
			Paxos         interface{}
			Range         string
			DefaultSubnet string
		}{
			paxos.PaxosStatus(allocator.paxos),
			allocator.universe.String(),
			defaultSubnet.String(),
		}
	}
	return <-resultChan
}
