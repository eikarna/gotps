package main

import (
	"sync"

	"github.com/codecat/go-libs/log"

	enet "github.com/eikarna/gotops"
	clients "github.com/eikarna/gotps/clients"
	"github.com/eikarna/gotps/items"
	pkt "github.com/eikarna/gotps/packet"
)

var (
	once          sync.Once
	GrowtopiaPort uint16 = 17091

	globalPeer []enet.Peer
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
	itemInfo, err := items.SerializeItemsDat("items.dat")
	if err != nil {
		log.Error("Itemsdat: %s", err.Error())
	}
	// The event loop
	for true {
		// Wait until the next event
		ev := host.Service(100)

		if ev != nil {
			once.Do(func() { log.Info("Server Successfully started on 0.0.0.0:%d", GrowtopiaPort) })
		}

		switch ev.GetType() {
		case enet.EventNone:
			{
				break
			}
		case enet.EventConnect:
			{
				clients.OnConnect(ev.GetPeer(), host, itemInfo, globalPeer) //Handle Client OnConnect
				break
			}
		case enet.EventDisconnect:
			{
				clients.OnDisConnect(ev.GetPeer(), host, itemInfo, globalPeer) //Handle Client OnDisConnect
				break
			}

		case enet.EventReceive: // A peer sent us some data
			// Get the packet
			packet := ev.GetPacket()
			// We must destroy the packet when we're done with it
			defer packet.Destroy()

			switch packet.GetData()[0] { //Net Message Type
			case 2:
				{
					clients.OnTextPacket(ev.GetPeer(), host, pkt.GetMessageFromPacket(packet), itemInfo, globalPeer)
					break
				}
			case 3:
				{
					clients.OnTextPacket(ev.GetPeer(), host, pkt.GetMessageFromPacket(packet), itemInfo, globalPeer)
					break
				}
			default:
				{
					clients.OnTankPacket(ev.GetPeer(), host, ev.GetPacket(), itemInfo, globalPeer)
					// log.Error("Unhandled type packet: %d", packet.GetData()[0])
					break
				}
			}
			break
		}
	}

	// Destroy the host when we're done with it
	host.Destroy()
	enet.Deinitialize()
}
