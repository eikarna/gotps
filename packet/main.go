package packet

import (
	"encoding/binary"
	"github.com/eikarna/gotops"
)

// By: Haikal (Kipas)
func GetMessageFromPacket(packet enet.Packet) string {
	packet.GetData()[len(packet.GetData())-1] = 0
	return string(packet.GetData()[4:])
}

// By: Haikal (Kipas)
func SendPacket(peer enet.Peer, gameMessageType int32, strData string) int {
	packetSize := 5 + len(strData)
	netPacket := make([]byte, packetSize)

	binary.LittleEndian.PutUint32(netPacket[0:4], uint32(gameMessageType))
	if strData != "" {
		copy(netPacket[4:4+len(strData)], []byte(strData))
	}
	netPacket[4+len(strData)] = 0
	packet, err := enet.NewPacket(netPacket, enet.PacketFlagReliable)
	if err != nil {
		panic(err)
	}
	peer.SendPacket(packet, 0)
	return 1
}
