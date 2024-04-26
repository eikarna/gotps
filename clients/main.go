package clients

import (
	log "github.com/codecat/go-libs/log"
	enet "github.com/eikarna/gotops"
)

func OnConnect(peer enet.Peer, host enet.Host) {
	log.Info("New Client Connected %s", peer.GetAddress().String())

}

func OnDisConnect(peer enet.Peer, host enet.Host) {
	log.Info("New Client Disconnected %s", peer.GetAddress().String())
}
