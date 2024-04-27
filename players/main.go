package players

import enet "github.com/eikarna/gotops"

type Players struct {
	TankIDName    string
	TankIDPass    string
	RequestedName string
	IpAddress     string

	Country string

	Userid int32
	Netid  int32

	Peer enet.Peer
}

func (p *Players) GetTankName() string {
	return p.TankIDName
}

func (p *Players) GetTankPass() string {
	return p.TankIDPass
}

func (p *Players) GetPeer() enet.Peer {
	return p.Peer
}

func NewPlayer() *Players {
	player := &Players{}

	return player
}
