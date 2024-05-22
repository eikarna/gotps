package main

import (
	"sync"

	//	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"github.com/codecat/go-libs/log"
	enet "github.com/eikarna/gotops"
	clients "github.com/eikarna/gotps/clients"
	"github.com/eikarna/gotps/items"
	pkt "github.com/eikarna/gotps/packet"
	// "github.com/vmihailenco/msgpack/v5"
	"time"
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

	// Print IP Server Address and it's Port
	log.Info("IP Address Server: %s, Port: %d", host.GetAddress().String(), host.GetAddress().GetPort())

	// GTPS Support
	host.EnableChecksum()
	host.CompressWithRangeCoder()
	log.Warn("Loading \"items.dat\"..")
	startTimestamp := time.Now()
	itemInfo, err := items.SerializeItemsDat("items.dat", startTimestamp)
	if err != nil {
		log.Error("Itemsdat: %s", err.Error())
	}
	// The event loop
	for true {
		// Wait until the next event
		ev := host.Service(1000)

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
				clients.OnConnect(ev.GetPeer(), host, itemInfo) //Handle Client OnConnect
				break
			}
		case enet.EventDisconnect:
			{
				clients.OnDisConnect(ev.GetPeer(), host, itemInfo) //Handle Client OnDisConnect
				break
			}

		case enet.EventReceive: // A peer sent us some data
			// Get the packet
			packet := ev.GetPacket()
			// We must destroy the packet when we're done with it1
			defer packet.Destroy()

			switch packet.GetData()[0] { //Net Message Type
			case 2:
				{
					clients.OnTextPacket(ev.GetPeer(), host, pkt.GetMessageFromPacket(packet), itemInfo)
					break
				}
			case 3:
				{
					clients.OnTextPacket(ev.GetPeer(), host, pkt.GetMessageFromPacket(packet), itemInfo)
					break
				}
			case 4:
				{
					clients.OnTankPacket(ev.GetPeer(), host, packet, itemInfo)
					break
				}
			case 22:
				{
					pkt.SendPacket(ev.GetPeer(), 21, "")
					break
				}
			default:
				{
					log.Error("Unhandled type packet: %d", packet.GetData()[0])
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
