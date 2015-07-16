package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/weaveworks/weave/ipam"
	"github.com/weaveworks/weave/nameserver"
	"github.com/weaveworks/weave/net/address"
	weave "github.com/weaveworks/weave/router"
	"net/http"
	"strings"
	"text/template"
)

// Strip escaped newlines from template
func escape(template string) string {
	return strings.Replace(template, "\\\n", "", -1)
}

var statusTemplate = escape(`\
          Service: router
             Name: {{.Router.Name}}
         NickName: {{.Router.NickName}}
       Encryption: {{.Router.Encryption}}
    PeerDiscovery: {{.Router.PeerDiscovery}}
            Peers: {{len .Router.Peers}}
             MACs: {{len .Router.MACs}}
    UnicastRoutes: {{len .Router.Routes.Unicast}}
  BroadcastRoutes: {{len .Router.Routes.Broadcast}}
      DirectPeers: {{len .Router.ConnectionMaker.DirectPeers}}
     Reconnecting: {{len .Router.ConnectionMaker.Reconnects}}

          Service: ipam
{{if .IPAM.Paxos}}\
        Consensus: {{.IPAM.Paxos.Consensus}}
           Quorum: {{.IPAM.Paxos.Quorum}}
       KnownNodes: {{.IPAM.Paxos.KnownNodes}}
{{end}}\
            Range: {{.IPAM.Range}}
    DefaultSubnet: {{.IPAM.DefaultSubnet}}

          Service: dns
           Domain: {{.DNS.Domain}}
             Port: {{.DNS.Port}}
              TTL: {{.DNS.TTL}}\
`)

func HandleHTTP(muxRouter *mux.Router,
	router *weave.Router,
	allocator *ipam.Allocator,
	defaultSubnet address.CIDR,
	ns *nameserver.Nameserver,
	dnsserver *nameserver.DNSServer) {

	muxRouter.Methods("GET").Path("/status").Headers("Accept", "application/json").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			json, _ := router.StatusJSON(version)
			w.Header().Set("Content-Type", "application/json")
			w.Write(json)
		})

	muxRouter.Methods("GET").Path("/status").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			tmpl, err := template.New("status").Parse(statusTemplate)
			if err != nil {
				fmt.Println(err)
				return
			}

			status := Status(router, allocator, defaultSubnet, ns, dnsserver)

			err = tmpl.Execute(w, status)
			if err != nil {
				fmt.Println(err)
				return
			}

		})

}
