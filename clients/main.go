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

func ParseUserData(peer enet.Peer, text string) {
	// Iterate over the lines to find the requestedName key
	lines := strings.Split(text, "\n")
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
				player.PlayerMap[peer].RequestedName = parts[1]
				break
			}
		case "TankIDName":
			{
				player.PlayerMap[peer].TankIDName = parts[1]
				break
			}
		case "TankIDPass":
			{
				player.PlayerMap[peer].TankIDPass = parts[1]
				break
			}
		case "protocol":
			{
				aa, err := strconv.ParseUint(parts[1], 10, 32)
				if err != nil {
					log.Error("Error Protocol:", err)
				}
				player.PlayerMap[peer].Protocol = uint32(aa)
				break
			}
		case "country":
			{
				player.PlayerMap[peer].Country = parts[1]
				break
			}
		case "PlatformID":
			{
				aa, err := strconv.ParseUint(parts[1], 10, 32)
				if err != nil {
					log.Error("Error PlatformID:", err)
				}
				player.PlayerMap[peer].PlatformID = uint32(aa)

				break
			}
		case "gid":
			{
				player.PlayerMap[peer].Gid = parts[1]
				break
			}
		case "rid":
			{
				player.PlayerMap[peer].Rid = parts[1]
				break
			}
		case "deviceVersion":
			{
				aa, err := strconv.ParseUint(parts[1], 10, 32)
				if err != nil {
					log.Error("Error DeviceVersion:", err)
				}
				player.PlayerMap[peer].DeviceVersion = uint32(aa)
				break
			}
		default:
			{
				continue
			}
		}
	}

}

func OnCommand(peer enet.Peer, host enet.Host, cmd string, isCommand bool) {
	lowerCmd := strings.ToLower(cmd)
	if isCommand {
		fn.LogMsg(peer, "`6"+cmd)
	}
	if strings.HasPrefix(lowerCmd, "/help") {
		fn.LogMsg(peer, "Help Command >> /help /ip")
	} else if strings.HasPrefix(lowerCmd, "/myip") {
		fn.LogMsg(peer, "Your IP Address: %s", player.GetPlayer(peer).IpAddress)
	} else {
		fn.LogMsg(peer, "`4Unknown command.``  Enter `$/?`` for a list of valid commands.")
	}
}

func OnChatInput(peer enet.Peer, host enet.Host, text string) {
	if strings.Contains(text, "player_chat=") || text == " " || len(strings.TrimSpace(text)) == 0 || (len(text) > 0 && text[0] == '`' && len(text) < 3) {
		return
	}

	chatPrefixBuble := player.GetChatPrefix(peer)
	if chatPrefixBuble == "`$" {
		chatPrefixBuble = ""
	}

	for _, currentPeer := range player.PlayerMap {
		if player.NotSafePlayer(currentPeer.Peer) {
			continue
		}
		if player.PlayerMap[peer].CurrentWorld == player.PlayerMap[currentPeer.Peer].CurrentWorld {
			fn.ConsoleMsg(currentPeer.Peer, 0, "CP:_PL:0_OID:_CT:[W]_ `6<`w"+player.GetPlayerName(peer)+"`6> "+player.GetChatPrefix(peer)+text)
			fn.TalkBubble(currentPeer.Peer, player.PlayerMap[peer].NetID, 1000, false, "CP:_PL:0_OID:_player_chat=%s", chatPrefixBuble+text)
		}
	}
}

func OnPlayerMove(peer enet.Peer, packet enet.Packet) {
	for _, currentPeer := range player.PlayerMap {
		if player.NotSafePlayer(currentPeer.Peer) {
			continue
		}
		if player.PlayerMap[peer].CurrentWorld == player.PlayerMap[currentPeer.Peer].CurrentWorld {
			movePacket := packet.GetData()
			binary.LittleEndian.PutUint16(movePacket[8:], uint16(player.PlayerMap[peer].NetID))
			packet, err := enet.NewPacket(movePacket, enet.PacketFlagReliable)
			if err != nil {
				panic(err)
			}
			currentPeer.Peer.SendPacket(packet, 0)
		}
	}
}

func OnPlayerExitWorld(peer enet.Peer) {
	if player.NotSafePlayer(peer) {
		return
	}
	if player.PlayerMap[peer].CurrentWorld == "" {
		return
	}
	world, err := worlds.GetWorld(player.PlayerMap[peer].CurrentWorld)
	if err != nil {
		log.Error(err.Error())
	}
	for _, currentPeer := range player.PlayerMap {
		if player.NotSafePlayer(currentPeer.Peer) {
			continue
		}
		if player.PlayerMap[peer].CurrentWorld == player.PlayerMap[currentPeer.Peer].CurrentWorld {
			fn.OnRemove(currentPeer.Peer, int(player.PlayerMap[peer].NetID))
			fn.ConsoleMsg(currentPeer.Peer, 0, "`5<`0%s`` left, `w%d`5 others here>``", player.GetPlayerName(peer), world.PlayersIn)
		}
	}

	player.PlayerMap[peer].CurrentWorld = ""
	player.PlayerMap[peer].SpawnX = 0
	player.PlayerMap[peer].SpawnY = 0
	world.PlayersIn--
	if world.PlayersIn < 0 {
		world.PlayersIn = 0
	}

	fn.SendWorldMenu(peer)
}

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
		ParseUserData(peer, text)
	} else if len(text) > 6 && text[:6] == "action" {

		if strings.HasPrefix(text[7:], "enter_game") {
			fn.SendWorldMenu(peer)
			// player.NewPlayer(peer)
			fn.LogMsg(peer, "Where would you like to go? (`w%d`` Online)", host.ConnectedPeers())
		} else if strings.HasPrefix(text[7:], "join_request") {
			log.Info("Invent Size: %d", byte(player.PlayerMap[peer].InventorySize))
			fn.SendInventory(player.Players, peer)
			worldName := strings.ToUpper(strings.Split(text[25:], "\n")[0])
			fn.LogMsg(peer, "Sending you to world (%s) (%d)", worldName, len(worldName))
			OnEnterGameWorld(peer, host, worldName)
		} else if strings.HasPrefix(text[7:], "input") {
			UserText := strings.Split(strings.Split(text[7:], "\n")[1], "|")[2]
			log.Info("User Input Text: %s", UserText)
			fn.ConsoleMsg(peer, 0, "CP:_PL:0_OID:_CT:[W]_ `6<`w%s`6> %s", player.PlayerMap[peer].RequestedName, UserText)
			fn.TalkBubble(peer, player.PlayerMap[peer].NetID, 100, true, UserText)
			if strings.HasPrefix(UserText, "get") {
				log.Info("GetPlayer Return: %v", player.GetPlayer(peer))
			}
		} else if strings.HasPrefix(text[7:], "quit_to_exit") {
			fn.ListActiveWorld[player.PlayerMap[peer].CurrentWorld]--
			player.PlayerMap[peer].CurrentWorld = ""
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
	if player.NotSafePlayer(peer) {
		return
	}
	if len(packet.GetData()) < 3 {
		return
	}

	Tank := &tankpacket.TankPacket{}
	Tank.SerializeFromMem(packet.GetData())

	switch Tank.PacketType {
	case 0:
		{ //player movement
			player.PlayerMap[peer].PosX = Tank.X
			player.PlayerMap[peer].PosY = Tank.Y
			fn.LogMsg(peer, "[Movement] X:%d, Y:%d", Tank.X, Tank.Y)
			OnPlayerMove(peer, packet)
			break
		}
	case 3:
		{ //punch / place
			switch Tank.Value {
			case 18:
				{
					fn.OnPunch(peer, Tank)
					break
				}
			case 32:
				{
					// fn.OnPlace(peer, Tank)
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
			if player.GetPlayer(peer).CurrentWorld != "" {
				OnPlayerExitWorld(peer)
			}
			break
			//fn.SendDoor(Tank, player.Players, peer)
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
			//player.PlayerMap[peer].NetID = Tank.NetID
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
	player.PlayerMap[peer].CurrentWorld = world.Name
	if int(world.PlayersIn) < 1 {
		world.PlayersIn = 1
	} else {
		world.PlayersIn++
	}
	// Simple Fix
	player.PlayerMap[peer].NetID = uint32(world.PlayersIn)
	for _, currentPlayer := range player.PlayerMap {
		if player.NotSafePlayer(currentPlayer.Peer) {
			continue
		}
		log.Info("%v", currentPlayer)
		if currentPlayer.CurrentWorld == world.Name {
			currentPeer := currentPlayer.Peer
			fn.ConsoleMsg(currentPeer, 0, "`5<`w%s ``entered, `w%d`` others here`5>", player.PlayerMap[peer].RequestedName, world.PlayersIn)
			fn.TalkBubble(currentPeer, player.PlayerMap[peer].NetID, 500, true, "`5<`w%s ``entered, `w%d`` others here`5>", player.PlayerMap[peer].RequestedName, world.PlayersIn)
			fn.OnSpawn(currentPeer, world.PlayersIn, world.PlayersIn, int32(SpawnX), int32(SpawnY), "`6@"+player.PlayerMap[peer].RequestedName, player.PlayerMap[peer].Country, false, true, true, true)
		}
	}
	fn.ListActiveWorld[world.Name] = int(world.PlayersIn)
}
