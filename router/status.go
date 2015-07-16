package router

import (
	"time"
)

func RouterStatus(router *Router) interface{} {
	type routerStatus struct {
		Encryption      bool
		PeerDiscovery   bool
		Name            string
		NickName        string
		Interface       interface{}
		MACs            interface{}
		Peers           interface{}
		Routes          interface{}
		ConnectionMaker interface{}
	}

	return &routerStatus{
		router.UsingPassword(),
		router.PeerDiscovery,
		router.Ourself.Name.String(),
		router.Ourself.NickName,
		router.Iface,
		macsStatus(router.Macs),
		peersStatus(router.Peers),
		routesStatus(router.Routes),
		connectionMakerStatus(router.ConnectionMaker)}
}

func connectionMakerStatus(cm *ConnectionMaker) interface{} {
	type connectionMakerStatus struct {
		DirectPeers []string
		Reconnects  map[string]interface{}
	}

	// We need to Refresh first in order to clear out any 'attempting'
	// connections from cm.targets that have been established since
	// the last run of cm.checkStateAndAttemptConnections. These
	// entries are harmless but do represent stale state that we do
	// not want to report.
	cm.Refresh()
	resultChan := make(chan connectionMakerStatus, 0)
	cm.actionChan <- func() bool {
		status := connectionMakerStatus{
			DirectPeers: []string{},
			Reconnects:  make(map[string]interface{}),
		}
		for peer := range cm.directPeers {
			status.DirectPeers = append(status.DirectPeers, peer)
		}

		for address, target := range cm.targets {
			status.Reconnects[address] = targetStatus(target)
		}
		resultChan <- status
		return false
	}
	return <-resultChan
}

func targetStatus(target *Target) interface{} {
	t := struct {
		Attempting  bool      `json:"Attempting,omitempty"`
		TryAfter    time.Time `json:"TryAfter,omitempty"`
		TryInterval string    `json:"TryInterval,omitempty"`
		LastError   string    `json:"LastError,omitempty"`
	}{
		Attempting:  target.attempting,
		TryAfter:    target.tryAfter,
		TryInterval: target.tryInterval.String(),
	}
	if target.lastError != nil {
		t.LastError = target.lastError.Error()
	}
	return t
}

func macsStatus(cache *MacCache) interface{} {
	type macStatus struct {
		Mac      string
		Name     string
		NickName string
		LastSeen time.Time
	}

	var macStatuses []interface{}
	for key, entry := range cache.table {
		macStatuses = append(macStatuses,
			&macStatus{
				intmac(key).String(),
				entry.peer.Name.String(),
				entry.peer.NickName,
				entry.lastSeen})
	}

	return macStatuses
}

func peersStatus(peers *Peers) []interface{} {
	type peerStatus struct {
		Name        string
		NickName    string
		UID         PeerUID
		Version     uint64
		Connections []interface{}
	}

	var peerStatuses []interface{}
	peers.ForEach(func(peer *Peer) {
		var connections []interface{}
		if peer == peers.ourself.Peer {
			for conn := range peers.ourself.Connections() {
				connections = append(connections, connectionStatus(conn))
			}
		} else {
			// Modifying peer.connections requires a write lock on
			// Peers, and since we are holding a read lock (due to the
			// ForEach), access without locking the peer is safe.
			for _, conn := range peer.connections {
				connections = append(connections, connectionStatus(conn))
			}
		}
		peerStatuses = append(peerStatuses,
			&peerStatus{
				peer.Name.String(),
				peer.NickName,
				peer.UID,
				peer.version,
				connections})
	})

	return peerStatuses
}

func connectionStatus(conn Connection) interface{} {
	return struct {
		Name        string
		NickName    string
		TCPAddr     string
		Outbound    bool
		Established bool
	}{
		conn.Remote().Name.String(),
		conn.Remote().NickName,
		conn.RemoteTCPAddr(),
		conn.Outbound(),
		conn.Established()}
}

func routesStatus(routes *Routes) interface{} {
	routes.RLock()
	defer routes.RUnlock()

	type uni struct {
		Dest, Via string
	}
	type broad struct {
		Source string
		Via    []string
	}
	var r struct {
		Unicast   []*uni
		Broadcast []*broad
	}
	for name, hop := range routes.unicast {
		r.Unicast = append(r.Unicast, &uni{name.String(), hop.String()})
	}
	for name, hops := range routes.broadcast {
		var hopStrings []string
		for _, hop := range hops {
			hopStrings = append(hopStrings, hop.String())
		}
		r.Broadcast = append(r.Broadcast, &broad{name.String(), hopStrings})
	}

	return r
}
