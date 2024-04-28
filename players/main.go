package players

import (
	"encoding/binary"

	enet "github.com/eikarna/gotops"
)

type Role int

const (
	PLAYER Role = iota
	VIP
	ADMIN
	MOD
	DEV
	FOUNDER
)

type Inventory struct {
	ItemID int16
	Amount int16
}

type Players struct {
	TankIDName    string
	TankIDPass    string
	RequestedName string
	IpAddress     string
	CurrentWorld  string
	Roles         Role
	Country       string
	NickUse       string

	Userid int32
	Netid  int32

	SpawnX int32
	SpawnY int32

	UpdateClothes bool

	Peer enet.Peer
	Inv  []Inventory
}

var PlayerMap = make(map[enet.Peer]*Players)

func GetPlayer(peer enet.Peer) *Players {
	player, exists := PlayerMap[peer]
	if !exists {
		return nil
	}
	return player
}

func NotSafePlayer(peer enet.Peer) bool {
	return GetPlayer(peer) == nil
}

func GetRoleNick(peer enet.Peer) string {
	if NotSafePlayer(peer) {
		return ""
	}
	switch GetPlayer(peer).Roles {
	case VIP:
		return "`w[`cVVIP`w]"
	case ADMIN:
		return "`w[`4ADMIN`w]"
	case MOD:
		return "`#@"
	case DEV:
		return "`6@"
	case FOUNDER:
		return "`b@"
	default:
		return ""
	}
}

func GetChatPrefix(peer enet.Peer) string {
	switch GetPlayer(peer).Roles {
	case PLAYER:
		return "`$"
	case VIP:
		return "`c"
	case ADMIN:
		return "`4"
	case MOD:
		return "`^"
	case DEV:
		return "`5"
	case FOUNDER:
		return "`5"
	}
	return "`$"
}

func GetPlayerName(peer enet.Peer) string {
	if NotSafePlayer(peer) {
		return ""
	}
	displayName := GetRoleNick(peer)
	if len(displayName) != 0 {
		displayName += " "
	}

	if GetPlayer(peer).TankIDName != "" {
		displayName += GetPlayer(peer).TankIDName
	} else {
		displayName += GetPlayer(peer).RequestedName
	}

	return displayName
}

func HasItem(peer enet.Peer, itemid int) bool {
	inventory := GetPlayer(peer).Inv
	for _, item := range inventory {
		if item.ItemID == int16(itemid) {
			return true
		}
	}
	return false
}

func GetCountItemFromInventory(peer enet.Peer, itemid int) int16 {
	inventory := GetPlayer(peer).Inv
	for _, item := range inventory {
		if item.ItemID == int16(itemid) {
			return item.Amount
		}
	}
	return 0
}

func UpdateInventory(peer enet.Peer) {
	if NotSafePlayer(peer) {
		return
	}
	inv := GetPlayer(peer).Inv
	invSize := uint32(len(inv))
	netid := -1
	packetLen := 66 + (invSize * 4) + 4
	buffer := make([]byte, packetLen)
	binary.LittleEndian.PutUint32(buffer[0:], 4)                        //net message type
	binary.LittleEndian.PutUint32(buffer[4:], 9)                        //packet type
	binary.LittleEndian.PutUint32(buffer[8:], uint32(netid))            //netid default -1
	binary.LittleEndian.PutUint32(buffer[16:], uint32(8))               //char state
	binary.LittleEndian.PutUint32(buffer[56:], uint32(6+(invSize*4)+4)) //payload inv size
	binary.LittleEndian.PutUint32(buffer[60:], uint32(1))               //payload inv size
	binary.LittleEndian.PutUint32(buffer[61:], uint32(invSize))         //payload inv size
	binary.LittleEndian.PutUint32(buffer[65:], uint32(invSize))         //payload inv size
	memPos := 67
	for _, Inven := range inv {
		binary.LittleEndian.PutUint32(buffer[memPos:], uint32(Inven.ItemID))
		memPos += 2
		buffer[memPos] = byte(Inven.Amount)
		memPos += 2
	}

	packet, err := enet.NewPacket(buffer, enet.PacketFlagReliable)
	if err != nil {
		panic(err)
	}
	peer.SendPacket(packet, 0)
}
