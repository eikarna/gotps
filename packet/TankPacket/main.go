package tankpacket

import "encoding/binary"

type TankPacket struct {
	PacketType     uint32
	NetID          uint32
	CharacterState uint32
	Value          uint32
	X              uint32
	Y              uint32
	XSpeed         uint32
	YSpeed         uint32
	PunchX         uint32
	PunchY         uint32
}

func (packet *TankPacket) Serialize(packetLength int32, createPacket bool) []byte {
	if createPacket {
		packetLength += 4
	}
	tank := make([]byte, packetLength) //default 56

	if createPacket {
		binary.LittleEndian.PutUint32(tank[0:], uint32(4))
	}
	binary.LittleEndian.PutUint32(tank[4:], uint32(packet.PacketType))
	binary.LittleEndian.PutUint32(tank[4+4:], uint32(packet.NetID))
	binary.LittleEndian.PutUint32(tank[4+12:], uint32(packet.CharacterState))
	binary.LittleEndian.PutUint32(tank[4+20:], uint32(packet.Value))
	binary.LittleEndian.PutUint32(tank[4+24:], uint32(packet.X))
	binary.LittleEndian.PutUint32(tank[4+28:], uint32(packet.Y))
	binary.LittleEndian.PutUint32(tank[4+32:], uint32(packet.XSpeed))
	binary.LittleEndian.PutUint32(tank[4+36:], uint32(packet.YSpeed))
	binary.LittleEndian.PutUint32(tank[4+44:], uint32(packet.PunchX))
	binary.LittleEndian.PutUint32(tank[4+48:], uint32(packet.PunchY))
	return tank
}

func (tank *TankPacket) SerializeFromMem(data []byte) *TankPacket {
	tank.PacketType = binary.LittleEndian.Uint32(data[4+0:])
	tank.NetID = binary.LittleEndian.Uint32(data[4+4:])
	tank.CharacterState = binary.LittleEndian.Uint32(data[4+12:])
	tank.Value = binary.LittleEndian.Uint32(data[4+20:])
	tank.X = binary.LittleEndian.Uint32(data[4+24:])
	tank.Y = binary.LittleEndian.Uint32(data[4+28:])
	tank.XSpeed = binary.LittleEndian.Uint32(data[4+32:])
	tank.YSpeed = binary.LittleEndian.Uint32(data[4+36:])
	tank.PunchX = binary.LittleEndian.Uint32(data[4+44:])
	tank.PunchY = binary.LittleEndian.Uint32(data[4+48:])
	return tank
}
