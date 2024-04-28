package clients

import (
	"encoding/binary"
	"strconv"
	"strings"

	log "github.com/codecat/go-libs/log"
	enet "github.com/eikarna/gotops"
	fn "github.com/eikarna/gotps/functions"
	items "github.com/eikarna/gotps/items"
	pkt "github.com/eikarna/gotps/packet"
	tankpacket "github.com/eikarna/gotps/packet/TankPacket"
	players "github.com/eikarna/gotps/players"
	"github.com/eikarna/gotps/utils"
	"github.com/eikarna/gotps/worlds"
)

func OnConnect(peer enet.Peer, host enet.Host, items *items.ItemInfo) {
	log.Info("New Client Connected %s", peer.GetAddress().String())
	players.PlayerMap[peer] = &players.Players{
		IpAddress: peer.GetAddress().String(),
		Peer:      peer,
	}
	pkt.SendPacket(peer, 1, "") //hello response
}

func OnPlayerMove(peer enet.Peer, packet enet.Packet) {
	for _, currentPeer := range players.PlayerMap {
		if players.NotSafePlayer(currentPeer.Peer) {
			continue
		}
		if players.GetPlayer(peer).CurrentWorld == players.GetPlayer(currentPeer.Peer).CurrentWorld {
			movePacket := packet.GetData()
			binary.LittleEndian.PutUint16(movePacket[8:], uint16(players.GetPlayer(peer).Netid))
			packet, err := enet.NewPacket(movePacket, enet.PacketFlagReliable)
			if err != nil {
				panic(err)
			}
			currentPeer.Peer.SendPacket(packet, 0)
		}
	}
}

func OnPlayerExitWorld(peer enet.Peer) {
	if players.NotSafePlayer(peer) {
		return
	}
	if players.GetPlayer(peer).CurrentWorld == "" {
		return
	}
	world, err := worlds.GetWorld(players.GetPlayer(peer).CurrentWorld)
	if err != nil {
		log.Error(err.Error())
	}
	for _, currentPeer := range players.PlayerMap {
		if players.NotSafePlayer(currentPeer.Peer) {
			continue
		}
		if players.GetPlayer(peer).CurrentWorld == players.GetPlayer(currentPeer.Peer).CurrentWorld {
			fn.OnRemove(currentPeer.Peer, int(players.GetPlayer(peer).Netid))
			fn.OnConsoleMessage(currentPeer.Peer, "`5<`0"+players.GetPlayerName(peer)+"`` left, `w"+strconv.Itoa(int(world.TotalPlayer))+" `5others here>``", 0)
		}
	}

	players.GetPlayer(peer).CurrentWorld = ""
	world.TotalPlayer--
	if world.TotalPlayer < 0 {
		world.TotalPlayer = 0
	}

	fn.SendWorldMenu(peer)
}

func OnDisConnect(peer enet.Peer, host enet.Host, items *items.ItemInfo) {
	log.Info("New Client Disconnected %s", peer.GetAddress().String())
	if players.GetPlayer(peer).CurrentWorld != "" {
		OnPlayerExitWorld(peer)
	}
	delete(players.PlayerMap, peer)
	peer.SetData(nil)
}

func OnCommand(peer enet.Peer, host enet.Host, cmd string, isCommand bool) {
	lowerCmd := strings.ToLower(cmd)
	if isCommand {
		fn.LogMsg(peer, "`6"+cmd)
	}
	if strings.HasPrefix(lowerCmd, "/help") {
		fn.LogMsg(peer, "Help Command >> /help /ip")
	} else if strings.HasPrefix(lowerCmd, "/myip") {
		fn.LogMsg(peer, "Your IP Addres: %s", players.GetPlayer(peer).IpAddress)
	} else {
		fn.LogMsg(peer, "`4Unknown command.``  Enter `$/?`` for a list of valid commands.")
	}
}

func OnChatInput(peer enet.Peer, host enet.Host, text string) {
	if strings.Contains(text, "player_chat=") || text == " " || len(strings.TrimSpace(text)) == 0 || (len(text) > 0 && text[0] == '`' && len(text) < 3) {
		return
	}

	chatPrefixBuble := players.GetChatPrefix(peer)
	if chatPrefixBuble == "`$" {
		chatPrefixBuble = ""
	}

	for _, currentPeer := range players.PlayerMap {
		if players.NotSafePlayer(currentPeer.Peer) {
			continue
		}
		if players.GetPlayer(peer).CurrentWorld == players.GetPlayer(currentPeer.Peer).CurrentWorld {
			fn.OnConsoleMessage(currentPeer.Peer, "CP:_PL:0_OID:_CT:[W]_ `6<`w"+players.GetPlayerName(peer)+"`6> "+players.GetChatPrefix(peer)+text, 0)
			fn.OnTalkBubble(currentPeer.Peer, int(players.GetPlayer(peer).Netid), "CP:_PL:0_OID:_player_chat="+chatPrefixBuble+text, false)
		}
	}
}

func OnTextPacket(peer enet.Peer, host enet.Host, text string, items *items.ItemInfo) {
	log.Info(text)
	if strings.Contains(text, "requestedName|") {
		fn.OnSuperMain(peer, items.GetItemHash())
		players.GetPlayer(peer).Inv = make([]players.Inventory, 18)
		players.GetPlayer(peer).Inv[0] = players.Inventory{ItemID: 18, Amount: 1}
		players.GetPlayer(peer).Inv[1] = players.Inventory{ItemID: 32, Amount: 1}
		players.GetPlayer(peer).Inv[2] = players.Inventory{ItemID: 7188, Amount: 200}
		players.GetPlayer(peer).Inv[3] = players.Inventory{ItemID: 242, Amount: 200}
		tp := utils.TextPacket{}
		tp.Parse(text)
		if tp.HasKey("tankIDName") && tp.HasKey("tankIDPass") {
			players.GetPlayer(peer).TankIDName = tp.GetFromKey("tankIDName")
			players.GetPlayer(peer).TankIDPass = tp.GetFromKey("tankIDPass")
		} else {
			players.GetPlayer(peer).RequestedName = tp.GetFromKey("requestedName")
		}
	} else if len(text) > 6 && text[:6] == "action" {
		if strings.HasPrefix(text[7:], "enter_game") {
			players.UpdateInventory(peer)
			fn.SendWorldMenu(peer)
		} else if strings.HasPrefix(text[7:], "join_request") {
			worldName := strings.ToUpper(strings.Split(text[25:], "\n")[0])
			fn.LogMsg(peer, "Sending you to world (%s) (%d)", worldName, len(worldName))
			OnEnterGameWorld(peer, host, worldName)
		} else if strings.HasPrefix(text[7:], "quit_to_exit") {
			OnPlayerExitWorld(peer)
		} else if strings.HasPrefix(text[7:], "quit") {
			peer.DisconnectLater(0)
		} else if strings.HasPrefix(text[7:], "input") {
			text := strings.Split(text[19:], "\n")[0]
			if text[0] == '/' {
				OnCommand(peer, host, text, true)
			} else {
				OnChatInput(peer, host, text)
			}
		} else if strings.HasPrefix(text[7:], "refresh_item_data") {
			//items dat update
			// fn.LogMsg(peer, "Updating item: %s, %d", players.GetPlayer(peer).IpAddress, len(items.FileBufferPacket))
			pkt.SendPacket(peer, 3, "action|log\nmsg|One moment, updating item data...")
			packet, err := enet.NewPacket(items.FileBufferPacket, enet.PacketFlagReliable)
			if err != nil {
				panic(err)
			}
			peer.SendPacket(packet, 0)
		}
	}
}

func OnTankPacket(peer enet.Peer, host enet.Host, packet enet.Packet, items *items.ItemInfo) {
	if len(packet.GetData()) < 3 || players.NotSafePlayer(peer) {
		return
	}

	var Tank = &tankpacket.TankPacket{}
	Tank.SerializeFromMem(packet.GetData())

	switch Tank.PacketType {
	case 0:
		{ //player movement
			OnPlayerMove(peer, packet)
			break
		}
	case 7:
		{
			if players.GetPlayer(peer).CurrentWorld != "" {
				OnPlayerExitWorld(peer)
			}
			break
		}
	case 3:
		{ //punch / place
			switch Tank.Value {
			case 18:
				{ //fist
					break
				}
			case 32:
				{
					break
				}
			default:
				{
					// world.Tiles[Tank.PunchX+(Tank.PunchY*uint32(world.SizeX))].Fg = 2
					break
				}
			}
			log.Info("Packet type: %d, val: %d", Tank.PacketType, Tank.Value)
			break
		}
	default:
		{
			// log.Info("Packet type: %d, val: %d", Tank.PacketType, Tank.Value)
			break
		}
	}

}

func OnEnterGameWorld(peer enet.Peer, host enet.Host, name string) {
	// players.GetPlayer(peer).UpdateClothes = true
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
				players.GetPlayer(peer).SpawnX = int32(i % int(world.SizeX) * 32)
				players.GetPlayer(peer).SpawnY = int32(i / int(world.SizeX) * 32)
				// players.GetPlayer(peer).SetSpawn(int32(i%int(world.SizeX)*32), int32(i/int(world.SizeX)*32))
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
	players.GetPlayer(peer).CurrentWorld = name
	players.GetPlayer(peer).Netid = world.TotalPlayer + 1
	world.TotalPlayer++
	fn.LogMsg(peer, "x: %d, y: %d", players.GetPlayer(peer).SpawnX, players.GetPlayer(peer).SpawnY)
	fn.OnSpawn(peer, players.GetPlayer(peer).Netid, 1, int32(players.GetPlayer(peer).SpawnX), int32(players.GetPlayer(peer).SpawnY), players.GetPlayerName(peer), "ccBadge", false, true, true, true)

	for _, currentPeer := range players.PlayerMap {
		if players.NotSafePlayer(currentPeer.Peer) {
			continue
		}
		if players.GetPlayer(currentPeer.Peer).TankIDName == players.GetPlayer(peer).TankIDName {
			continue
		}
		if players.GetPlayer(peer).CurrentWorld == players.GetPlayer(currentPeer.Peer).CurrentWorld {
			fn.OnSpawn(currentPeer.Peer, players.GetPlayer(peer).Netid, 1, int32(players.GetPlayer(peer).SpawnX), int32(players.GetPlayer(peer).SpawnY), players.GetPlayerName(peer), "ccBadge", false, true, true, false)
			fn.OnSpawn(peer, players.GetPlayer(currentPeer.Peer).Netid, 1, int32(players.GetPlayer(currentPeer.Peer).SpawnX), int32(players.GetPlayer(currentPeer.Peer).SpawnY), players.GetPlayerName(currentPeer.Peer), "ccBadge", false, true, true, false)

			fn.OnConsoleMessage(currentPeer.Peer, "`5<`0"+players.GetPlayerName(peer)+" `` entered, `w"+strconv.Itoa(int(world.TotalPlayer))+" `5others here>``", 0)
			fn.OnTalkBubble(currentPeer.Peer, int(players.GetPlayer(peer).Netid), "`5<`0"+players.GetPlayerName(peer)+" `` entered, `w"+strconv.Itoa(int(world.TotalPlayer))+" `5others here>``", true)
		}
	}
}
