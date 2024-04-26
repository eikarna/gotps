package clients

import (
	log "github.com/codecat/go-libs/log"
	enet "github.com/eikarna/gotops"
	pkt "github.com/eikarna/gotps/packet"
)

func OnConnect(peer enet.Peer, host enet.Host) {
	log.Info("New Client Connected %s", peer.GetAddress().String())
	pkt.SendPacket(peer, 1, "") //hello response
}

func OnDisConnect(peer enet.Peer, host enet.Host) {
	log.Info("New Client Disconnected %s", peer.GetAddress().String())
}
