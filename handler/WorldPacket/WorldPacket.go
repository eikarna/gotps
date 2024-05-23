package WorldPacket

import (
	"encoding/binary"
)

type WorldPacket struct {
	PacketType     uint32
	NetID          uint32
	CharacterState uint32
	PunchX         uint32
	PunchY         uint32
	PlantingTree   uint32
}

func (wp *WorldPacket) Serialize(packetLength uint32, createPacket bool) []byte {
	if createPacket {
		packetLength += 4
	}
	data := make([]byte, packetLength)
	if createPacket {
		binary.LittleEndian.PutUint32(data[0:], uint32(4))
	}
	binary.LittleEndian.PutUint32(data[4:], wp.PacketType)
	binary.LittleEndian.PutUint32(data[8:], wp.NetID)
	binary.LittleEndian.PutUint32(data[16:], wp.CharacterState)
	binary.LittleEndian.PutUint32(data[24:], wp.PunchX)
	binary.LittleEndian.PutUint32(data[28:], wp.PunchY)
	return data
}

func (wp *WorldPacket) SerializeFromMem(data []byte) *WorldPacket {
	wp.PacketType = binary.LittleEndian.Uint32(data[:4])
	wp.NetID = binary.LittleEndian.Uint32(data[4:8])
	wp.CharacterState = binary.LittleEndian.Uint32(data[12:16])
	wp.PunchX = binary.LittleEndian.Uint32(data[20:24])
	wp.PunchY = binary.LittleEndian.Uint32(data[24:28])
	return wp
}
