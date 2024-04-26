package main

import (
	"github.com/codecat/go-libs/log"
	"github.com/eikarna/gotops"
	pkt "github.com/eikarna/gotps/packet"
	"sync"
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

		// Do nothing if we didn't get any event
		if ev.GetType() == enet.EventNone {
			continue
		}

		switch ev.GetType() {
		case enet.EventConnect: // A new peer has connected
			log.Info("New peer connected: %s", ev.GetPeer().GetAddress())
			/*packet := ev.GetPacket()
			defer packet.Destroy()
			log.Info("Got Login Packet from %s: %s", ev.GetPeer().GetAddress(), GetMessageFromPacket(packet))*/
			if pkt.SendPacket(ev.GetPeer(), 1, "") == 1 {
				pkt.SendPacket(ev.GetPeer(), 3, "action|play_sfx\nfile|audio/piano_nice.wav\ndelayMS|0\n")
				pkt.SendPacket(ev.GetPeer(), 3, "action|log\nmsg|`2Hello World from gotops!")
				pkt.SendPacket(ev.GetPeer(), 3, "action|set_url\nurl|https://github.com/eikarna/gotops\nlabel|GOTOPS Repo")
				pkt.SendPacket(ev.GetPeer(), 3, "action|logon_fail\n")
			}

		case enet.EventDisconnect: // A connected peer has disconnected
			log.Info("Peer disconnected: %s", ev.GetPeer().GetAddress())

		case enet.EventReceive: // A peer sent us some data
			// Get the packet
			packet := ev.GetPacket()

			// We must destroy the packet when we're done with it
			defer packet.Destroy()

			// Get the bytes in the packet and handle the packet
			switch packet.GetData()[0] {
			// On Connect
			case 1:
				{
					log.Info("Packet Type %d: %s", packet.GetData()[0], pkt.GetMessageFromPacket(packet))

				}
			// On Change
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