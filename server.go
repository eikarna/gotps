package main

import (
	"sync"

	"github.com/codecat/go-libs/log"
	enet "github.com/eikarna/gotops"
	clients "github.com/eikarna/gotps/clients"
	pkt "github.com/eikarna/gotps/packet"
)

var (
	once          sync.Once
	GrowtopiaPort uint16 = 17091
)

func main() {
	// Initialize enet
	enet.Initialize()

	// Create a host listening on 0.0.0.0:17091
	host, err := enet.NewHost(enet.NewListenAddress(GrowtopiaPort), 1024, 1, 0, 0)
	if err != nil {
		log.Error("Couldn't create host: %s", err.Error())
		return
	}

	// GTPS Support
	host.EnableChecksum()
	host.CompressWithRangeCoder()

	// The event loop
	for true {
		// Wait until the next event
		ev := host.Service(100)

		if ev != nil {
			once.Do(func() { log.Info("Server Successfully started on 0.0.0.0:%d", GrowtopiaPort) })
		}

		switch ev.GetType() {
		case enet.EventConnect: // A new peer has connected
			log.Info("New peer connected: %s", ev.GetPeer().GetAddress())
			clients.OnConnect(ev.GetPeer(), host) //Handle Client OnConnect

		case enet.EventDisconnect: // A connected peer has disconnected
			log.Info("Peer disconnected: %s", ev.GetPeer().GetAddress())
			clients.OnDisConnect(ev.GetPeer(), host) //Handle Client OnDisConnect

		case enet.EventReceive: // A peer sent us some data
			// Get the packet
			packet := ev.GetPacket()
			// We must destroy the packet when we're done with it
			defer packet.Destroy()

			switch packet.GetData()[0] { //Net Message Type
			case 2:
				{
					log.Info("Packet Type %d: %s", packet.GetData()[0], pkt.GetMessageFromPacket(packet))
				}
			case 3:
				{
					log.Info("Packet Type %d: %s", packet.GetData()[0], pkt.GetMessageFromPacket(packet))
				}
			default:
				{
					log.Error("Unhandled type packet: %d, with data: %s", packet.GetData()[0], pkt.GetMessageFromPacket(packet))
				}
			}

		}
	}

	// Destroy the host when we're done with it
	host.Destroy()

	// Uninitialize enet
	enet.Deinitialize()
}
