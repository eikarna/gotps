package variants

import (
	"encoding/binary"

	enet "github.com/eikarna/gotops"
)

// variant
type Variant struct {
	index      int
	len        int
	packetData []byte
}

func NewVariant(delay int, NetID int) *Variant {
	packetData := make([]byte, 61)
	len := 61

	binary.LittleEndian.PutUint32(packetData[0:4], uint32(4))       //message type
	binary.LittleEndian.PutUint32(packetData[4:8], uint32(1))       //packet type variant
	binary.LittleEndian.PutUint32(packetData[8:12], uint32(NetID))  //netid
	binary.LittleEndian.PutUint32(packetData[16:20], uint32(8))     //characterState
	binary.LittleEndian.PutUint32(packetData[24:28], uint32(delay)) //delay

	return &Variant{
		index:      0,
		len:        len,
		packetData: packetData,
	}
}

func NewVariantRaw() *Variant {
	packetData := make([]byte, 61)
	len := 61
	binary.LittleEndian.PutUint32(packetData[0:4], uint32(4)) //message type
	binary.LittleEndian.PutUint32(packetData[4:8], uint32(1)) //packet type variant
	return &Variant{
		index:      0,
		len:        len,
		packetData: packetData,
	}
}

func (v *Variant) InsertInt(a1 int) {
	data := make([]byte, v.len+2+4)
	copy(data, v.packetData)
	data[v.len] = byte(v.index)
	data[v.len+1] = 0x9
	binary.LittleEndian.PutUint32(data[v.len+2:], uint32(a1))
	v.index++
	v.packetData = data
	v.len += 2 + 4
	v.packetData[60] = byte(v.index)
}

func (v *Variant) InsertUnsignedInt(a1 uint32) {
	data := make([]byte, v.len+2+4)
	copy(data, v.packetData)
	data[v.len] = byte(v.index)
	data[v.len+1] = 0x5
	binary.LittleEndian.PutUint32(data[v.len+2:], a1)
	v.index++
	v.packetData = data
	v.len += 2 + 4
	v.packetData[60] = byte(v.index)
}
func (v *Variant) InsertFloat(a1 float32) {
	data := make([]byte, v.len+2+4)
	copy(data, v.packetData)
	data[v.len] = byte(v.index)
	data[v.len+1] = 0x1
	binary.LittleEndian.PutUint32(data[v.len+2:], uint32(a1))
	v.index++
	v.packetData = data
	v.len += 2 + 4
	v.packetData[60] = byte(v.index)
}
func (v *Variant) InsertDoubleFloat(a1 float32, a2 float32) {
	data := make([]byte, v.len+2+8)
	copy(data, v.packetData)
	data[v.len] = byte(v.index)
	data[v.len+1] = 0x3
	binary.LittleEndian.PutUint32(data[v.len+2:], uint32(a1))
	binary.LittleEndian.PutUint32(data[v.len+6:], uint32(a2))
	v.index++
	v.packetData = data
	v.len += 2 + 8
	v.packetData[60] = byte(v.index)
}
func (v *Variant) InsertTripleFloat(a1 float32, a2 float32, a3 float32) {
	data := make([]byte, v.len+2+12)
	copy(data, v.packetData)
	data[v.len] = byte(v.index)
	data[v.len+1] = 0x4
	binary.LittleEndian.PutUint32(data[v.len+2:], uint32(a1))
	binary.LittleEndian.PutUint32(data[v.len+6:], uint32(a2))
	binary.LittleEndian.PutUint32(data[v.len+10:], uint32(a3))
	v.index++
	v.packetData = data
	v.len += 2 + 12
	v.packetData[60] = byte(v.index)
}
func (v *Variant) InsertString(a1 string) {
	data := make([]byte, v.len+2+len(a1)+4)
	copy(data, v.packetData)
	data[v.len] = byte(v.index)
	data[v.len+1] = 0x2
	binary.LittleEndian.PutUint32(data[v.len+2:], uint32(len(a1)))
	copy(data[v.len+6:], []byte(a1))
	v.index++
	v.packetData = data
	v.len += 2 + len(a1) + 4
	v.packetData[60] = byte(v.index)
}

func (v *Variant) Send(peer enet.Peer) {
	packet, err := enet.NewPacket(v.packetData, enet.PacketFlagReliable)
	if err != nil {
		panic(err)
	}
	peer.SendPacket(packet, 0)
}

func (v *Variant) SendBroadcast(host enet.Host) {
	packet, err := enet.NewPacket(v.packetData, enet.PacketFlagReliable)
	if err != nil {
		panic(err)
	}
	host.BroadcastPacket(packet, 0)
}
