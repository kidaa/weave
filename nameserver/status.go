package nameserver

func DNSStatus(ns *Nameserver, dnsServer *DNSServer) interface{} {
	return struct {
		Domain string
		Port   int
		TTL    uint32
	}{
		dnsServer.domain,
		dnsServer.port,
		dnsServer.ttl,
	}
}
