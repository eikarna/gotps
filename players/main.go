package players

import enet "github.com/eikarna/gotops"

var Players Player

type Role int

const (
	PLAYER Role = iota
	VIP
	ADMIN
	MOD
	DEV
	FOUNDER
)

type ItemInfo struct {
	ID  int
	Qty int16
}

type Player struct {
	TankIDName    string
	TankIDPass    string
	RequestedName string
	IpAddress     string
	Country       string
	UserID        uint32
	NetID         uint32
	Protocol      uint32
	GameVersion   string
	PlatformID    uint32
	DeviceVersion uint32
	MacAddr       string
	Rid           string
	Gid           string
	PlayerAge     uint32
	CurrentWorld  string
	Peer          enet.Peer
	PosX          uint32
	PosY          uint32
	PunchX        uint32
	PunchY        uint32
	SpawnX        uint32
	SpawnY        uint32
	Inventory     []ItemInfo
	InventorySize uint16
	Roles         Role
}

var PlayerMap = make(map[enet.Peer]*Player)

func (p *Player) GetTankName() string {
	return p.TankIDName
}

func (p *Player) GetTankPass() string {
	return p.TankIDPass
}

func (p *Player) GetPeer() enet.Peer {
	return p.Peer
}

func (p *Player) GetCountry() string {
	return p.Country
}

func (p *Player) GetPlatformID() uint32 {
	return p.PlatformID
}

func (p *Player) GetAge() uint32 {
	return p.PlayerAge
}

func (p *Player) GetProtocol() uint32 {
	return p.Protocol
}

func (p *Player) GetMac() string {
	return p.MacAddr
}

func (p *Player) GetDeviceVersion() uint32 {
	return p.DeviceVersion
}

func (p *Player) GetRid() string {
	return p.Rid
}

func (p *Player) GetGid() string {
	return p.Gid
}

func (p *Player) GetIp() string {
	return p.IpAddress
}

func (p *Player) GetUserid() uint32 {
	return p.UserID
}

func NewPlayer(peer enet.Peer) *Player {
	player := &Player{}
	return player
}

func GetPlayer(peer enet.Peer) *Player {
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
	inventory := GetPlayer(peer).Inventory
	for _, item := range inventory {
		if item.ID == itemid {
			return true
		}
	}
	return false
}

func GetCountItemFromInventory(peer enet.Peer, itemid int) int16 {
	inventory := GetPlayer(peer).Inventory
	for _, item := range inventory {
		if item.ID == itemid {
			return item.Qty
		}
	}
	return 0
}
