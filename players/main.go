package players

import enet "github.com/eikarna/gotops"

var Players Player

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
	rid           string
	gid           string
	PlayerAge     uint32
	CurrentWorld  string
	Peer          enet.Peer
	PosX          uint32
	PosY          uint32
	PunchX        uint32
	PunchY        uint32
}

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
	return p.rid
}

func (p *Player) GetGid() string {
	return p.gid
}

func (p *Player) GetIp() string {
	return p.IpAddress
}

func (p *Player) GetUserid() uint32 {
	return p.UserID
}

func NewPlayer() *Player {
	player := &Player{}
	return player
}
