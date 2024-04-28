package clients

import (
	"encoding/binary"
	"strings"

	log "github.com/codecat/go-libs/log"
	enet "github.com/eikarna/gotops"
	fn "github.com/eikarna/gotps/functions"
	items "github.com/eikarna/gotps/items"
	pkt "github.com/eikarna/gotps/packet"
	tankpacket "github.com/eikarna/gotps/packet/TankPacket"
	player "github.com/eikarna/gotps/players"
	"github.com/eikarna/gotps/worlds"
	"strconv"
)

var (
	SpawnX int
	SpawnY int
)

func OnConnect(peer enet.Peer, host enet.Host, items *items.ItemInfo, globalPeer []enet.Peer) {
	log.Info("New Client Connected %s", peer.GetAddress().String())
	player.PlayerMap[peer] = &player.Player{
		IpAddress: peer.GetAddress().String(),
		Peer:      peer,
	}
	pkt.SendPacket(peer, 1, "") //hello response
}

func OnDisConnect(peer enet.Peer, host enet.Host, items *items.ItemInfo, globalPeer []enet.Peer) {
	log.Info("New Client Disconnected %s", peer.GetAddress().String())
	delete(player.PlayerMap, peer)
}

func OnTextPacket(peer enet.Peer, host enet.Host, text string, items *items.ItemInfo, globalPeer []enet.Peer) {
	if strings.Contains(text, "requestedName|") {
		fn.OnSuperMain(peer, items.GetItemHash())
		lines := strings.Split(text, "\n")

		// Iterate over the lines to find the requestedName key
		for _, line := range lines {
			// Split each line into key and value parts
			parts := strings.Split(line, "|")
			// Check if the key is "requestedName"
			if len(parts) != 2 {
				continue
			}
			switch parts[0] {
			case "requestedName":
				{
					player.Players.RequestedName = parts[1]
					break
				}
			case "protocol":
				{
					aa, err := strconv.ParseUint(parts[1], 10, 32)
					if err != nil {
						log.Error("Error Protocol:", err)
					}
					player.Players.Protocol = uint32(aa)
					break
				}
			case "country":
				{
					player.Players.Country = parts[1]
					break
				}
			case "PlatformID":
				{
					aa, err := strconv.ParseUint(parts[1], 10, 32)
					if err != nil {
						log.Error("Error PlatformID:", err)
					}
					player.Players.PlatformID = uint32(aa)
					break
				}
			case "gid":
				{
					player.Players.Gid = parts[1]
					break
				}
			case "rid":
				{
					player.Players.Rid = parts[1]
					break
				}
			case "deviceVersion":
				{
					aa, err := strconv.ParseUint(parts[1], 10, 32)
					if err != nil {
						log.Error("Error DeviceVersion:", err)
					}
					player.Players.DeviceVersion = uint32(aa)
					break
				}
			default:
				{
					continue
				}
			}
		}
	} else if len(text) > 6 && text[:6] == "action" {

		if strings.HasPrefix(text[7:], "enter_game") {
			fn.SendWorldMenu(peer)
			// player.NewPlayer(peer)
			fn.LogMsg(peer, "Where would you like to go? (`w%d`` Online)", host.ConnectedPeers())
		} else if strings.HasPrefix(text[7:], "join_request") {
			log.Info("Invent Size: %d", byte(player.Players.InventorySize))
			fn.SendInventory(player.Players, peer)
			worldName := strings.ToUpper(strings.Split(text[25:], "\n")[0])
			fn.LogMsg(peer, "Sending you to world (%s) (%d)", worldName, len(worldName))
			OnEnterGameWorld(peer, host, worldName)
		} else if strings.HasPrefix(text[7:], "input") {
			UserText := strings.Split(strings.Split(text[7:], "\n")[1], "|")[2]
			log.Info("User Input Text: %s", UserText)
			fn.ConsoleMsg(peer, "CP:_PL:0_OID:_CT:[W]_ `6<`w%s`6> %s", player.Players.RequestedName, UserText)
			fn.TalkBubble(peer, 1, UserText)
			if strings.HasPrefix(UserText, "get") {
				log.Info("GetPlayer Return: %v", player.GetPlayer(peer))
			}
		} else if strings.HasPrefix(text[7:], "quit_to_exit") {
			player.Players.CurrentWorld = ""
			fn.SendWorldMenu(peer)
		} else if strings.HasPrefix(text[7:], "quit") {
			peer.DisconnectLater(0)
		} else {
			fn.LogMsg(peer, "Unhandled Action Packet type: %s", text[7:])
		}
	} else {
		fn.LogMsg(peer, "Unhandled TextPacket, msg: %v", text)
		log.Info("Unhandled TextPacket, msg: %v", text)
	}
}

func OnTankPacket(peer enet.Peer, host enet.Host, packet enet.Packet, items *items.ItemInfo, globalPeer []enet.Peer) {
	if len(packet.GetData()) < 3 {
		return
	}

	Tank := &tankpacket.TankPacket{}
	Tank.SerializeFromMem(packet.GetData())

	switch Tank.PacketType {
	case 0:
		{ //player movement
			player.Players.PosX = Tank.X
			player.Players.PosY = Tank.Y
			fn.LogMsg(peer, "[Movement] X:%d, Y:%d", Tank.X, Tank.Y)
			break
		}
	case 3:
		{ //punch / place
			switch Tank.Value {
			case 18:
				{
					player.Players.PunchX = Tank.PunchX
					player.Players.PunchY = Tank.PunchY
					testt := &tankpacket.TankPacket{
						PacketType:     3,
						NetID:          player.Players.NetID,
						CharacterState: Tank.CharacterState,
						Value:          Tank.Value,
						X:              player.Players.PosX,
						Y:              player.Players.PosY,
						XSpeed:         Tank.XSpeed,
						YSpeed:         Tank.YSpeed,
						PunchX:         player.Players.PunchX,
						PunchY:         player.Players.PunchY,
					}
					bbb := testt.Serialize(56, true)
					aaa, err := enet.NewPacket(bbb, enet.PacketFlagReliable)
					if err != nil {
						log.Error("Error Packet 3:", err)
					}
					peer.SendPacket(aaa, 0)
					fn.LogMsg(peer, "[Punch/Place] X:%d, Y:%d, Value:%d, NetID:%d", Tank.PunchX, Tank.PunchY, Tank.Value, Tank.NetID)
					break
				}
			default:
				{
					break
				}
			}
		}
	case 7:
		{
			// Door
			fn.SendDoor(*Tank, player.Players, peer)
		}
	case 18:
		{
			// Break
			pkt.SendPacket(peer, 3, "")
			fn.LogMsg(peer, "[Break] X: %d, Y:%d", Tank.PunchX, Tank.PunchY)
		}
	case 24:
		{
			// Check Client?
			log.Info("Check Client Msg: %v", Tank)
			player.Players.NetID = Tank.NetID
			//player.Players.UserID = Tank.UserID
			fn.LogMsg(peer, "[Client] Client Msg: %v (Value:%d)", Tank, Tank.Value)
		}
	case 32:
		{
			// Unknpwn
			pkt.SendPacket(peer, 3, "")
			fn.LogMsg(peer, "[Break2??] X: %d, Y:%d", Tank.PunchX, Tank.PunchY)

		}
	default:
		{
			//player.Players.NetID = Tank.NetID
			log.Info("Packet type: %d, val: %d", Tank.PacketType, Tank.Value)
			break
		}
	}

}

func OnEnterGameWorld(peer enet.Peer, host enet.Host, name string) {

	world, err := worlds.GetWorld(name)
	if err != nil {
		log.Error(err.Error())
	}
	nameLen := len(world.Name)
	totalPacketLen := 78 + nameLen + len(world.Tiles) + 24 + (8*len(world.Tiles) + (0 * 16))
	worldPacket := make([]byte, totalPacketLen)
	worldPacket[0] = 4  //game message
	worldPacket[4] = 4  //world packet type
	worldPacket[16] = 8 //char state
	worldPacket[66] = byte(len(world.Name))
	copy(worldPacket[68:], []byte(world.Name))

	worldPacket[nameLen+68] = byte(world.SizeX)
	worldPacket[nameLen+72] = byte(world.SizeY)
	binary.LittleEndian.PutUint16(worldPacket[nameLen+76:], uint16(world.TotalTiles))
	extraDataPos := 85 + nameLen

	for i := 0; i < int(world.TotalTiles); i++ {
		binary.LittleEndian.PutUint16(worldPacket[extraDataPos:], uint16(world.Tiles[i].Fg))
		binary.LittleEndian.PutUint16(worldPacket[extraDataPos+2:], uint16(world.Tiles[i].Bg))
		binary.LittleEndian.PutUint32(worldPacket[extraDataPos+4:], uint32(world.Tiles[i].Flags))

		switch world.Tiles[i].Fg {
		case 6:
			{
				worldPacket[extraDataPos+8] = 1 //block types
				binary.LittleEndian.PutUint16(worldPacket[extraDataPos+9:], uint16(len(world.Tiles[i].Label)))
				copy(worldPacket[extraDataPos+11:], []byte(world.Tiles[i].Label))

				SpawnX = (i % int(world.SizeX)) * 32
				SpawnY = (i / int(world.SizeX)) * 32
				extraDataPos += 4 + len(world.Tiles[i].Label)
			}
		default:
			{
				break
			}
		}

		extraDataPos += 8
	}

	packet, err := enet.NewPacket(worldPacket, enet.PacketFlagReliable)
	if err != nil {
		panic(err)
	}
	peer.SendPacket(packet, 0)
	player.Players.CurrentWorld = name
	//for i := 0; i < int(world.PlayersIn); i++ {
	fn.ConsoleMsg(peer, "`5<`w%s ``entered, `w%d`` others here`5>", player.Players.RequestedName, world.PlayersIn)
	fn.TalkBubble(peer, 1, "`5<`w%s ``entered, `w%d`` others here`5>", player.Players.RequestedName, world.PlayersIn)
	if int(world.PlayersIn) < 1 {
		world.PlayersIn = 1
	} else {
		world.PlayersIn++
	}
	fn.OnSpawn(peer, world.PlayersIn, world.PlayersIn, int32(SpawnX), int32(SpawnY), "`6@"+player.Players.RequestedName, player.Players.Country, false, true, true, true)
	//}
}
