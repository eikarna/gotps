package clients

import (
	"encoding/binary"
	"github.com/bvinc/go-sqlite-lite/sqlite3"
	log "github.com/codecat/go-libs/log"
	enet "github.com/eikarna/gotops"
	fn "github.com/eikarna/gotps/functions"
	items "github.com/eikarna/gotps/items"
	pkt "github.com/eikarna/gotps/packet"
	tankpacket "github.com/eikarna/gotps/packet/TankPacket"
	. "github.com/eikarna/gotps/players"
	"github.com/eikarna/gotps/worlds"
	"strings"
	// worldpacket "github.com/eikarna/gotps/worlds/WorldPacket"
	"runtime"
	// "strconv"
)

var (
	SpawnX int
	SpawnY int
)

func OnTileUpdate(packet enet.Packet, peer enet.Peer, Tank *tankpacket.TankPacket) {
	world, err := worlds.LoadWorld(db, PInfo(peer).CurrentWorld)
	if err != nil {
		fn.ConsoleMsg(peer, 0, "Some data is missing in this world! exiting..")
		OnPlayerExitWorld(peer, db)
	}
	switch Tank.Value {
	case 18:
		{
			fn.OnPunch(peer, Tank, PInfo(peer).CurrentWorld, world)
			break
		}
	case 32:
		{
			// fn.OnPlace(peer, Tank)
			break
		}
	default:
		{
			// test, ok := worlds.Worlds[PInfo(peer).CurrentWorld]
			Coords := Tank.PunchX + (Tank.PunchY * uint32(world.SizeX))
			/*if test.Tiles[Coords].Fg == 6 {
			        fn.TalkBubble(peer, PInfo(peer).NetID, 100, false, "don't break/replace the white door!")
			}*/
			if world.Tiles[Coords].Fg == 0 {
				decodedWp := &tankpacket.TankPacket{}
				decodedWp.SerializeFromMem(packet.GetData()[4:])
				log.Info("Got Place Packet with type %d: %v", decodedWp.PacketType, decodedWp)
				PlacePack := &tankpacket.TankPacket{
					PacketType:     3,
					NetID:          PInfo(peer).NetID,
					CharacterState: decodedWp.CharacterState,
					Value:          decodedWp.Value,
					X:              decodedWp.X,
					Y:              decodedWp.Y,
					XSpeed:         decodedWp.XSpeed,
					YSpeed:         decodedWp.YSpeed,
					PunchX:         decodedWp.PunchX,
					PunchY:         decodedWp.PunchY,
				}
				PlacePacket := PlacePack.Serialize(56, true)
				world.Tiles[Coords].Fg = int16(decodedWp.Value)
				inventory := PInfo(peer).Inventory
				fn.TalkBubble(peer, PInfo(peer).NetID, 100, false, "ID: %d, Qty: %d", decodedWp.Value, GetCountItemFromInventory(peer, int(decodedWp.Value)))
				for i := range inventory {
					if inventory[i].ID == int(decodedWp.Value) {
						inventory[i].Qty--
						ReducePack := &tankpacket.TankPacket{
							PacketType: 13,
							Value:      uint32(inventory[i].ID),
						}
						ReducedPacket := ReducePack.Serialize(56, true)
						Packet, err := enet.NewPacket(ReducedPacket, enet.PacketFlagReliable)
						if err != nil {
							log.Error("Packet type 13:", err.Error())
						}
						peer.SendPacket(Packet, 0)
						break

					} else {
						continue
					}
				}
				PInfo(peer).Inventory = inventory
				Packet, err := enet.NewPacket(PlacePacket, enet.PacketFlagReliable)
				if err != nil {
					log.Error("Error Packet 3:", err)
				}
				for _, currentPeer := range GetPeers(PlayerMap) {
					if NotSafePlayer(currentPeer) {
						continue
					}
					if PInfo(peer).CurrentWorld == PInfo(currentPeer).CurrentWorld {
						currentPeer.SendPacket(Packet, 0)
					}
				}
				/*WorldPack := &worldpacket.WorldPacket{
				          PacketType:   15,
				          NetID:        PInfo(peer).NetID,
				          PunchX:       decodedWp.PunchX,
				          PunchY:       decodedWp.PunchY,
				          PlantingTree: uint32(test.Tiles[Coords].Fg),
				  }
				  //WorldPack.Serialize(56, true)
				  bbb := WorldPack.Serialize(56, true)
				  aaa, err := enet.NewPacket(bbb, enet.PacketFlagReliable)
				  if err != nil {
				          log.Error("Error Packet 15:", err)
				  }
				  for _, currentPeer := range PlayerMap {
				          if NotSafePlayer(currentPeer.Peer) {
				                  continue
				          }
				          if PInfo(peer).CurrentWorld == currentPeer.CurrentWorld {
				                  currentPeer.Peer.SendPacket(aaa, 0)
				          }
				  }*/
				//go worlds.SaveWorld(db, PInfo(peer).CurrentWorld, *world)
				fn.TalkBubble(peer, PInfo(peer).NetID, 100, false, "Updating Block at %d", Coords)
			}
			break
		}
	}

	/*
	   case 7:
	           {
	                   // Door
	                   if GetPlayer(peer).CurrentWorld != "" {
	                           OnPlayerExitWorld(peer, db)
	                   }
	                   break
	                   //fn.SendDoor(Tank, Players, peer)
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
	                   //PInfo(peer).NetID = Tank.NetID
	                   log.Info("Packet type: %d, val: %d", Tank.PacketType, Tank.Value)
	                   break
	           }
	   }*/
}

func OnCommand(peer enet.Peer, host enet.Host, cmd string, isCommand bool) {
	lowerCmd := strings.ToLower(cmd)
	if isCommand {
		fn.LogMsg(peer, "`6"+cmd)
	}
	if strings.HasPrefix(lowerCmd, "/help") {
		fn.LogMsg(peer, "Help Command >> /help /ip")
	} else if strings.HasPrefix(lowerCmd, "/myip") {
		fn.LogMsg(peer, "Your IP Address: %s", GetPlayer(peer).IpAddress)
	} else if strings.HasPrefix(lowerCmd, "/info") {
		cpuUsage := runtime.NumCPU()
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		memoryUsage := m.Sys / 1024 / 1024
		fn.ConsoleMsg(peer, 0, "CPU: %d Core(s)\nAlloc: %d MB\nTotalAlloc: %d MB\nSys: %d MB\nNumGC: %d Thread(s)", cpuUsage, m.Alloc/1024/1024, m.TotalAlloc/1024/1024, memoryUsage, m.NumGC)
		fn.TalkBubble(peer, PInfo(peer).NetID, 100, false, "CPU: %d Core(s)\nAlloc: %d MB\nTotalAlloc: %d MB\nSys: %d MB\nNumGC: %d Thread(s)", cpuUsage, m.Alloc/1024/1024, m.TotalAlloc/1024/1024, memoryUsage, m.NumGC)
	} else {
		fn.LogMsg(peer, "`4Unknown command.``  Enter `$/?`` for a list of valid commands.")
	}
}

func OnChatInput(peer enet.Peer, host enet.Host, text string) {
	if strings.Contains(text, "player_chat=") || text == " " || len(strings.TrimSpace(text)) == 0 || (len(text) > 0 && text[0] == '`' && len(text) < 3) {
		return
	}

	chatPrefixBuble := GetChatPrefix(peer)
	if chatPrefixBuble == "`$" {
		chatPrefixBuble = ""
	}

	for _, currentPeer := range GetPeers(PlayerMap) {
		if NotSafePlayer(currentPeer) {
			continue
		}
		if PInfo(peer).CurrentWorld == PInfo(currentPeer).CurrentWorld {
			fn.ConsoleMsg(currentPeer, 0, "CP:_PL:0_OID:_CT:[W]_ `6<`w"+GetPlayerName(peer)+"`6> "+GetChatPrefix(peer)+text)
			fn.TalkBubble(currentPeer, PInfo(peer).NetID, 100, false, "CP:_PL:0_OID:_player_chat=%s", chatPrefixBuble+text)
		}
	}
}

func OnPlayerMove(peer enet.Peer, packet enet.Packet) {
	movePacket := packet.GetData()
	Tank := tankpacket.TankPacket{}
	Tank.SerializeFromMem(packet.GetData()[4:])
	PInfo(peer).PosX = Tank.X
	PInfo(peer).PosY = Tank.Y
	PInfo(peer).SpawnX = uint32(Tank.X)
	PInfo(peer).SpawnY = uint32(Tank.Y)
	packet, err := enet.NewPacket(movePacket, enet.PacketFlagReliable)
	if err != nil {
		panic(err)
	}
	for _, currentPeer := range GetPeers(PlayerMap) {
		if NotSafePlayer(currentPeer) {
			continue
		}
		if PInfo(peer).CurrentWorld == PInfo(currentPeer).CurrentWorld {

			currentPeer.SendPacket(packet, 0)
		}
	}
}

func OnPlayerExitWorld(peer enet.Peer, db *sqlite3.Conn) {
	if NotSafePlayer(peer) {
		return
	}
	if PInfo(peer).CurrentWorld == "" {
		return
	}
	world, err := worlds.LoadWorld(db, PInfo(peer).CurrentWorld)
	if err != nil {
		peer.DisconnectLater(0)
		log.Error("[OnPlayerExitWorld] Worlds with name: %s is not found in our database!", world.Name)
	}
	for _, currentPeer := range GetPeers(PlayerMap) {
		if NotSafePlayer(currentPeer) {
			continue
		}
		if PInfo(peer).CurrentWorld == PInfo(currentPeer).CurrentWorld {
			fn.OnRemove(currentPeer, int(PInfo(peer).NetID))
			fn.ConsoleMsg(currentPeer, 0, "`5<`0%s`` left, `w%d`5 others here>``", GetPlayerName(peer), world.PlayersIn)
		}
	}

	PInfo(peer).CurrentWorld = ""
	PInfo(peer).SpawnX = 0
	PInfo(peer).SpawnY = 0
	world.PlayersIn--
	if world.PlayersIn < 0 {
		world.PlayersIn = 0
	}

	fn.SendWorldMenu(peer)
	//codedWorld := worlds.AutoTagMsgpackStruct(world)
	if PInfo(peer).TankIDName != "" {
		go SavePlayer(db, *PInfo(peer), PInfo(peer).TankIDName)
	} else {
		go SavePlayer(db, *PInfo(peer), PInfo(peer).Rid)
	}
	go worlds.SaveWorld(db, world.Name, *world)
	log.Error("Saving Worlds with name: %s", world.Name)
}

func OnConnect(peer enet.Peer, host enet.Host, items *items.ItemInfo, globalPeer []enet.Peer, db *sqlite3.Conn) {
	log.Info("New Client Connected %s", peer.GetAddress().String())
	PInfo(peer) = &Player{
		IpAddress: peer.GetAddress().String(),
		//Roles:     6,
	}
	/*PlayerConnect := &Player{
		IpAddress: peer.GetAddress().String(),
		Peer:      peer,
		Roles:     6,
	}
	SavePlayer(db, *PlayerConnect)*/
	pkt.SendPacket(peer, 1, "") //hello response
}

func OnDisConnect(peer enet.Peer, host enet.Host, items *items.ItemInfo, globalPeer []enet.Peer, db *sqlite3.Conn) {
	log.Info("New Client Disconnected %s", peer.GetAddress().String())
	delete(PlayerMap, peer)
}

func OnTextPacket(peer enet.Peer, host enet.Host, text string, items *items.ItemInfo, globalPeer []enet.Peer, db *sqlite3.Conn) {
	if strings.Contains(text, "requestedName|") {
		fn.OnSuperMain(peer, items.GetItemHash())
		ParseUserData(db, text, peer)
	} else if len(text) > 6 && text[:6] == "action" {

		if strings.HasPrefix(text[7:], "enter_game") {
			fn.SendWorldMenu(peer)
			// NewPlayer(peer)
			fn.LogMsg(peer, "Where would you like to go? (`w%d`` Online)", host.ConnectedPeers())
		} else if strings.HasPrefix(text[7:], "join_request") {
			log.Info("Invent Size: %d", byte(PInfo(peer).InventorySize))
			fn.SendInventory(Players, peer)
			worldName := strings.ToUpper(strings.Split(text[25:], "\n")[0])
			fn.LogMsg(peer, "Sending you to world (%s) (%d)", worldName, len(worldName))
			OnEnterGameWorld(peer, host, worldName, db)
		} else if strings.HasPrefix(text[7:], "input") {
			text := strings.Split(text[19:], "\n")[0]
			if text[0] == '/' {
				OnCommand(peer, host, text, true)
			} else {
				OnChatInput(peer, host, text)
			}
		} else if strings.HasPrefix(text[7:], "quit_to_exit") {
			/*fn.ListActiveWorld[PInfo(peer).CurrentWorld]--
			PInfo(peer).CurrentWorld = ""*/
			OnPlayerExitWorld(peer, db)
		} else if strings.HasPrefix(text[7:], "quit") {
			peer.DisconnectLater(0)
		} else if strings.HasPrefix(text[7:], "refresh_item_data") {
			fn.LogMsg(peer, "One moment, Updating item data...")
			packet, err := enet.NewPacket(items.FileBufferPacket, enet.PacketFlagReliable)
			if err != nil {
				panic(err)
			}
			peer.SendPacket(packet, 0)
		} else {
			fn.LogMsg(peer, "Unhandled Action Packet type: %s", text[7:])
		}
	} else {
		fn.LogMsg(peer, "Unhandled TextPacket, msg: %v", text)
		log.Info("Unhandled TextPacket, msg: %v", text)
	}
}

func OnTankPacket(peer enet.Peer, host enet.Host, packet enet.Packet, items *items.ItemInfo, globalPeer []enet.Peer, db *sqlite3.Conn) {
	if NotSafePlayer(peer) {
		return
	}
	if len(packet.GetData()) < 3 {
		return
	}

	Tank := &tankpacket.TankPacket{}
	Tank.SerializeFromMem(packet.GetData()[4:])

	switch Tank.PacketType {
	case 0:
		{ //player movement
			fn.LogMsg(peer, "[Movement] X:%f, Y:%f", Tank.X, Tank.Y)
			OnPlayerMove(peer, packet)
			break
		}
	case 3:
		{
			OnTileUpdate(packet, peer, Tank)
		}
	case 7:
		{
			// Door
			if GetPlayer(peer).CurrentWorld != "" {
				OnPlayerExitWorld(peer, db)
			}
			break
			//fn.SendDoor(Tank, Players, peer)
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
			//PInfo(peer).NetID = Tank.NetID
			log.Info("Packet type: %d, val: %d", Tank.PacketType, Tank.Value)
			break
		}
	}

}

func OnEnterGameWorld(peer enet.Peer, host enet.Host, name string, db *sqlite3.Conn) {
	//log.Info("[OnEnterGameWorld] Player Data: %v", PInfo(peer))
	if NotSafePlayer(peer) {
		fn.LogMsg(peer, "`4Invalid Player Data!``")
		return
	}
	/*world, ok := worlds.Worlds[name]
	if !ok {
		world = *worlds.GenerateWorld(name, 100, 60)
	}*/
	world, err := worlds.LoadWorld(db, name)
	if err != nil {
		world = worlds.GenerateWorld(name, 100, 60)
		//codedWorld := AutoTagMsgpackStruct(world)
		go worlds.SaveWorld(db, name, *world)
		log.Error("Worlds with name: %s is not found in our database!", name)
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
		// log.Info("Loaded Tiles: %v", world.Tiles[i])
		if world.Tiles[i].Fg == 7188 {
			log.Error("Found BGL! %d", i)
		}
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
		case 7188:
			{
				worldPacket[extraDataPos+8] = 3
				musicBpm := 100 * -1
				binary.LittleEndian.PutUint16(worldPacket[extraDataPos+9:], uint16(0))
				binary.LittleEndian.PutUint16(worldPacket[extraDataPos+10:], uint16(PInfo(peer).UserID))
				binary.LittleEndian.PutUint16(worldPacket[extraDataPos+18:], uint16(musicBpm))
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
	PInfo(peer).CurrentWorld = world.Name
	if int(world.PlayersIn) < 1 {
		world.PlayersIn = 1
	} else {
		world.PlayersIn++
	}
	// Simple Fix
	PInfo(peer).NetID = uint32(world.PlayersIn)
	for _, currentPeer := range GetPeers(PlayerMap) {
		log.Info("%v", currentPeer)
		if PInfo(currentPeer).CurrentWorld == world.Name {
			fn.ConsoleMsg(currentPeer, 0, "`5<`w%s ``entered, `w%d`` others here`5>", PInfo(peer).Name, world.PlayersIn)
			fn.TalkBubble(currentPeer, PInfo(peer).NetID, 500, true, "`5<`w%s ``entered, `w%d`` others here`5>", PInfo(peer).Name, world.PlayersIn)
			isLocal := PInfo(currentPeer).NetID == PInfo(peer).NetID
			fn.OnSpawn(currentPeer, world.PlayersIn, world.PlayersIn, int32(SpawnX), int32(SpawnY), PInfo(peer).Name, "ccBadge", false, true, true, isLocal)
			break
		} else {
			continue
		}
	}
	fn.ListActiveWorld[world.Name] = int(world.PlayersIn)
}
