package tankpacket

import (
	"encoding/binary"
	"math"
)

type TankPacket struct {
	PacketType     uint32
	NetID          uint32
	CharacterState uint32
	Value          uint32
	X              float32
	Y              float32
	XSpeed         float32
	YSpeed         float32
	PunchX         uint32
	PunchY         uint32
}

func (tank *TankPacket) Serialize(packetLength int32, createPacket bool) []byte {
	if createPacket {
		packetLength += 4
	}
	data := make([]byte, packetLength) //default 56

	if createPacket {
		binary.LittleEndian.PutUint32(data[0:], uint32(4))
	}
	binary.LittleEndian.PutUint32(data[4:], tank.PacketType)
	binary.LittleEndian.PutUint32(data[8:], tank.NetID)
	binary.LittleEndian.PutUint32(data[16:], tank.CharacterState)
	binary.LittleEndian.PutUint32(data[24:], tank.Value)
	binary.LittleEndian.PutUint32(data[28:], math.Float32bits(tank.X))
	binary.LittleEndian.PutUint32(data[32:], math.Float32bits(tank.Y))
	binary.LittleEndian.PutUint32(data[36:], math.Float32bits(tank.XSpeed))
	binary.LittleEndian.PutUint32(data[44:], math.Float32bits(tank.YSpeed))
	binary.LittleEndian.PutUint32(data[48:], tank.PunchX)
	binary.LittleEndian.PutUint32(data[52:], tank.PunchY)
	return data
}

func (tank *TankPacket) SerializeFromMem(data []byte) *TankPacket {
	tank.PacketType = binary.LittleEndian.Uint32(data[:4])
	tank.NetID = binary.LittleEndian.Uint32(data[4:8])
	tank.CharacterState = binary.LittleEndian.Uint32(data[12:16])
	tank.Value = binary.LittleEndian.Uint32(data[20:24])
	tank.X = math.Float32frombits(binary.LittleEndian.Uint32(data[24:28]))
	tank.Y = math.Float32frombits(binary.LittleEndian.Uint32(data[28:32]))
	tank.XSpeed = math.Float32frombits(binary.LittleEndian.Uint32(data[32:36]))
	tank.YSpeed = math.Float32frombits(binary.LittleEndian.Uint32(data[36:40]))
	tank.PunchX = binary.LittleEndian.Uint32(data[44:48])
	tank.PunchY = binary.LittleEndian.Uint32(data[48:52])
	return tank
}
