package main

import (
	"github.com/weaveworks/weave/ipam"
	"github.com/weaveworks/weave/nameserver"
	"github.com/weaveworks/weave/net/address"
	"github.com/weaveworks/weave/router"
)

func Status(
	r *router.Router,
	allocator *ipam.Allocator,
	defaultSubnet address.CIDR,
	ns *nameserver.Nameserver,
	dnsserver *nameserver.DNSServer) interface{} {

	return struct {
		Router interface{}
		IPAM   interface{}
		DNS    interface{}
	}{
		router.RouterStatus(r),
		ipam.IPAMStatus(allocator, defaultSubnet),
		nameserver.DNSStatus(ns, dnsserver),
	}

}
