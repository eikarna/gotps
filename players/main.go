package players

import (
	"errors"
	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"github.com/codecat/go-libs/log"
	"github.com/eikarna/gotops"
	"github.com/vmihailenco/msgpack/v5"
	"strconv"
	"strings"
)

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
	Name          string
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
	PosX          float32
	PosY          float32
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

/*func (p *Player) GetPeer() enet.Peer {
	return p.
}*/

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
	inventory := PlayerMap[peer].Inventory
	for _, item := range inventory {
		if item.ID == itemid {
			return item.Qty
		}
	}
	return 0
}

func PInfo(peer enet.Peer) *Player {
	return PlayerMap[peer]
}

// SaveWorld saves a single World struct to the database with the given name
func SavePlayer(db *sqlite3.Conn, player Player, name string) error {
	// Serialize World struct to MessagePack binary forma
	/*name := ""
	if len(player.TankIDName) > 0 && len(player.TankIDPass) > 0 {
		name = player.TankIDName
	} else {
		name = player.RequestedName
	}*/
	playerBytes, err := msgpack.Marshal(player)
	if err != nil {
		return err
	}

	// Prepare statement for insertion
	stmt, err := db.Prepare("INSERT INTO players (name, data) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Insert name and binary data into the database
	stmt.Exec(name, playerBytes)
	return nil
}

func LoadPlayer(db *sqlite3.Conn, name string) (*Player, error) {
	var player Player

	// Query the data
	row, err := db.Prepare("SELECT data FROM players WHERE name = ?", name)

	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	defer row.Close()

	for {
		hasRow, err := row.Step()
		if err != nil {
			log.Fatal(err.Error())

			// panic("Error Parsing World: " + name)
		}
		if !hasRow {
			// The query is finished
			log.Fatal("Player %s Not Found in our database!", name)
			return nil, errors.New("Player " + name + " Not Found in our database!")
			break
		}

		// Deserialize MessagePack binary data into World struct
		var playerBytes []byte
		if err := row.Scan(&playerBytes); err != nil {
			log.Fatal(err.Error())

			return nil, err
			break
		}
		if err := msgpack.Unmarshal(playerBytes, &player); err != nil {
			log.Fatal(err.Error())
			return nil, err
			break
		}
	}
	return &player, nil
}

func ParseUserData(db *sqlite3.Conn, text string, peer enet.Peer) {
	// Initialize a map to store key-value pairs
	userData := make(map[string]string)

	// Split the text into lines
	lines := strings.Split(text, "\n")

	// Iterate over the lines
	for _, line := range lines {
		// Split each line into key and value parts using the delimiter "|"
		parts := strings.Split(line, "|")
		if len(parts) != 2 {
			// Skip lines that don't contain a key-value pair
			continue
		}
		// Store the key-value pair in the userData map
		userData[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	print(text)

	// Now you can access the parsed key-value pairs from the userData map
	// var isGuest bool

	// Convert protocol and platformID to uint32
	protocol, err := strconv.ParseUint(userData["protocol"], 10, 32)
	if err != nil {
		log.Error("Error Protocol:", err)
		return
	}
	platformID, err := strconv.ParseUint(userData["platformID"], 10, 32)
	if err != nil {
		log.Error("Error platformID:", err)
		return
	}
	// Convert deviceVersion to uint32
	deviceVersion, err := strconv.ParseUint(userData["deviceVersion"], 10, 32)
	if err != nil {
		log.Error("Error DeviceVersion:", err)
		return
	}

	// Create a player struct
	NewPlayerData := Player{
		RequestedName: userData["requestedName"],
		Protocol:      uint32(protocol),
		Country:       userData["country"],
		PlatformID:    uint32(platformID),
		Gid:           userData["gid"],
		Rid:           userData["rid"],
		DeviceVersion: uint32(deviceVersion),
	}

	// Check if TankIDName exists
	if _, ok := userData["tankIDName"]; ok {
		// TankIDName exists, parse and save as registered user
		// isGuest = false
		IsRegistered, err := LoadPlayer(db, userData["tankIDName"])
		if err != nil {
			NewPlayerData.TankIDName = userData["tankIDName"]
			NewPlayerData.TankIDPass = userData["tankIDPass"]
			NewPlayerData.Name = userData["tankIDName"]
			SavePlayer(db, NewPlayerData, NewPlayerData.TankIDName)
			// PlayerMap[peer].Peer = peer
			PlayerMap[peer] = &NewPlayerData
		} else {
			// PlayerMap[peer].Peer = peer
			PlayerMap[peer] = IsRegistered
		}
	} else {
		// TankIDName does not exist, save as guest user
		// isGuest = true
		IsRegistered, err := LoadPlayer(db, userData["rid"])
		if err != nil {
			NewPlayerData.Name = userData["requestedName"]
			SavePlayer(db, NewPlayerData, NewPlayerData.Rid)
			//PlayerMap[peer].Peer = peer
			PlayerMap[peer] = &NewPlayerData
		} else {
			// PlayerMap[peer].Peer = peer
			PlayerMap[peer] = IsRegistered
		}
	}

	// Optionally, you can log whether the player is registered or a guest
	/*if isGuest {
		log.Info("Guest player saved: %s", NewPlayerData.RequestedName)
	} else {
		log.Info("Registered player saved: %s", NewPlayerData.TankIDName)
	}
	PlayerMap[peer] = &NewPlayerData*/
}

// parseUint parses a uint32 from a string and returns 0 if parsing fails
func parseUint(s string) uint32 {
	val, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0
	}
	return uint32(val)
}

func GetPeers(playerMap map[enet.Peer]*Player) []enet.Peer {
	peers := make([]enet.Peer, 0, len(playerMap))
	for peer := range playerMap {
		peers = append(peers, peer)
	}
	return peers
}
