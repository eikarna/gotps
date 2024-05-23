package players

import (
	"encoding/gob"
	"errors"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/codecat/go-libs/log"
	enet "github.com/eikarna/gotops"
	// fn "github.com/eikarna/gotps/functions"
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
	ID  int   `clover:"ID"`
	Qty int16 `clover:"Qty"`
}

type Cloth struct {
	Hair     float32 `clover:"Hair"`
	Necklace float32 `clover:"Necklace"`
	Pants    float32 `clover:"Pants"`
	Shirt    float32 `clover:"Shirt"`
	Feet     float32 `clover:"Feet"`
	Back     float32 `clover:"Back"`
	Mask     float32 `clover:"Mask"`
	Face     float32 `clover:"Face"`
	Hand     float32 `clover:"Hand"`
}

//var PlayerGuestNum = make([]int, 900)

type Player struct {
	TankIDName    string     `clover:"TankIDName"`
	TankIDPass    string     `clover:"TankIDPass"`
	RequestedName string     `clover:"RequestedName"`
	Name          string     `clover:"Name"`
	IpAddress     string     `clover:"IpAddress"`
	Country       string     `clover:"Country"`
	UserID        uint32     `clover:"UserID"`
	NetID         uint32     `clover:"NetID"`
	Protocol      uint32     `clover:"Protocol"`
	GameVersion   string     `clover:"GameVersion"`
	PlatformID    string     `clover:"PlatformID"`
	DeviceVersion uint32     `clover:"DeviceVersion"`
	MacAddr       string     `clover:"MacAddr"`
	Rid           string     `clover:"Rid"`
	Gid           string     `clover:"Gid"`
	PlayerAge     uint32     `clover:"PlayerAge"`
	CurrentWorld  string     `clover:"CurrentWorld"`
	PosX          float32    `clover:"PosX"`
	PosY          float32    `clover:"PosY"`
	SpawnX        uint32     `clover:"SpawnX"`
	SpawnY        uint32     `clover:"SpawnY"`
	Inventory     []ItemInfo `clover:"Inventory"`
	InventorySize uint16     `clover:"InventorySize"`
	Roles         Role       `clover:"Roles"`
	Clothes       Cloth      `clover:"Clothes"`
	SkinColor     int        `clover:"SkinColor"`
	Gems          int        `clover:"Gems"`
	RotatedLeft   bool       `clover:"RotatedLeft"`
	PeerID        uint32     `clover:"PeerID"`
	IsOnline      bool       `clover:"IsOnline"`
}

var PlayerMap = make(map[enet.Peer]*Player)

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

	if GetPlayer(peer).TankIDName != "" {
		displayName += GetPlayer(peer).TankIDName
	} else {
		displayName += GetPlayer(peer).RequestedName
	}
	displayName += "``"

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

func ParseUserData(text string, host enet.Host, peer enet.Peer, ConsoleMsg func(peer enet.Peer, delay int, a ...interface{})) {
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
	peerId := peer.GetConnectID()
	for _, currentPeer := range host.ConnectedPeers() {
		if !NotSafePlayer(currentPeer) {
			if strings.EqualFold(PInfo(currentPeer).TankIDName, userData["tankIDName"]) && strings.EqualFold(PInfo(currentPeer).TankIDPass, userData["tankIDPass"]) {
				if PInfo(currentPeer).IsOnline {
					ConsoleMsg(peer, 0, "`4Already Logged In?`` It seems that this account already logged in by somebody else.")
					currentPeer.DisconnectLater(0)
				}
			}
		}
	}

	if NotSafePlayer(peer) {
		var loadedPlayer *Player
		var err error
		if _, ok := userData["tankIDName"]; ok {
			loadedPlayer, err = LoadPlayer(userData["tankIDName"])
		} else if _, ok := userData["requestedName"]; ok {
			loadedPlayer, err = LoadPlayer(userData["requestedName"])
		} else {
			log.Warn("Got invalid login packet from %s", peer.GetAddress().String())
			peer.DisconnectNow(0)
			return
		}
		if err != nil {
			// Now you can access the parsed key-value pairs from the userData map
			// var isGuest bool

			// Convert protocol and platformID to uint32
			protocol := parseUint(userData["protocol"])

			// Convert deviceVersion to uint32
			deviceVersion := parseUint(userData["deviceVersion"])

			userData["requestedName"] = userData["requestedName"] + "_" + strconv.Itoa(100+rand.Intn(899))

			// Create a player struct
			NewPlayerData := Player{
				RequestedName: userData["requestedName"],
				Protocol:      uint32(protocol),
				Country:       userData["country"],
				PlatformID:    userData["platformID"],
				Gid:           userData["gid"],
				Rid:           userData["rid"],
				DeviceVersion: uint32(deviceVersion),
				Roles:         5,
				PeerID:        peerId,
				SkinColor:     2864971775,
			}

			// Check if TankIDName exists
			if _, ok := userData["tankIDName"]; ok {
				// TankIDName exists, parse and save as registered user
				IsRegistered := PlayerMap[peer]
				if IsRegistered == nil {
					NewPlayerData.TankIDName = userData["tankIDName"]
					NewPlayerData.TankIDPass = userData["tankIDPass"]
					NewPlayerData.Name = userData["tankIDName"]
					log.Info("Growid player loaded & saved: %s", NewPlayerData.Name)
					PlayerMap[peer] = &NewPlayerData
				} else {
					PlayerMap[peer] = IsRegistered
					log.Info("Growid player loaded: %s", PlayerMap[peer].Name)
				}
			} else {
				// TankIDName does not exist, save as guest user
				IsRegistered := PlayerMap[peer]
				if IsRegistered == nil {
					NewPlayerData.Name = userData["requestedName"]
					PlayerMap[peer] = &NewPlayerData
					log.Info("Guest player loaded & saved: %s", NewPlayerData.Name)
				} else {
					PlayerMap[peer] = IsRegistered
					log.Info("Guest player loaded: %s", PlayerMap[peer].Name)
				}
			}
			return
		} else {
			loadedPlayer.PeerID = peerId
			loadedPlayer.IpAddress = peer.GetAddress().String()
			PlayerMap[peer] = loadedPlayer
			log.Info("%#v", loadedPlayer)
			return
		}
	} else {
		PlayerMap[peer].PeerID = peerId
		PlayerMap[peer].IpAddress = peer.GetAddress().String()
		return
	}
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

func SavePlayer(player *Player) error {
	if player == nil {
		return errors.New("SavePlayer: player is nil")
	}
	filePath := "database/players/_" + player.Name + ".bin"
	f, err := os.Create(filePath)
	if err != nil {
		return errors.New("Couldn't open file: " + err.Error())
	}
	defer f.Close()
	encoder := gob.NewEncoder(f)
	err = encoder.Encode(player)
	if err != nil {
		return errors.New("Encoding failed: " + err.Error())
	}
	return nil
}

func LoadPlayer(name string) (*Player, error) {
	filePath := "database/players/_" + name + ".bin"
	f, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New("Couldn't open file: " + err.Error())
	}
	defer f.Close()

	player := &Player{}
	decoder := gob.NewDecoder(f)
	err = decoder.Decode(player)
	if err != nil {
		return nil, errors.New("Decoding failed: " + err.Error())
	}
	return player, nil
}
