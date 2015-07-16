package paxos

import ()

func PaxosStatus(node *Node) interface{} {
	if node == nil {
		return nil
	}

	consensus, _ := node.Consensus()

	return struct {
		Consensus  bool
		KnownNodes int
		Quorum     uint
	}{
		consensus,
		len(node.knows),
		node.quorum,
	}
}
