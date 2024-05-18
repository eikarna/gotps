package clients

import (
	//	"encoding/binary"
	"encoding/binary"
	log "github.com/codecat/go-libs/log"
	enet "github.com/eikarna/gotops"
	fn "github.com/eikarna/gotps/functions"
	items "github.com/eikarna/gotps/items"
	pkt "github.com/eikarna/gotps/packet"
	tankpacket "github.com/eikarna/gotps/packet/TankPacket"
	. "github.com/eikarna/gotps/players"
	"github.com/eikarna/gotps/worlds"
	"runtime"
	"strconv"
	"strings"
)

var (
	SpawnX int
	SpawnY int
)

func OnTileUpdate(packet enet.Packet, peer enet.Peer, Tank *tankpacket.TankPacket, world *worlds.World) {
	switch Tank.Value {
	case 18:
		{
			fn.OnPunch(peer, Tank, world)
			break
		}
	case 20:
		{
			//fn.OnWrench(peer, Tank, world)
			break
		}
	case 7188:
		{
			// test, ok := worlds.Worlds[PInfo(peer).CurrentWorld]
			Coords := Tank.PunchX + (Tank.PunchY * uint32(world.SizeX))
			/*if test.Tiles[Coords].Fg == 6 {
			        fn.TalkBubble(peer, PInfo(peer).NetID, 100, false, "don't break/replace the white door!")
			}*/
			if world.Tiles[Coords].Fg == 0 {
				if worlds.Worlds[PInfo(peer).CurrentWorld].OwnerUid == 0 {
					lockPack := &tankpacket.TankPacket{
						PacketType:     15,
						PunchX:         Tank.PunchX,
						PunchY:         Tank.PunchY,
						CharacterState: Tank.CharacterState,
						NetID:          PInfo(peer).UserID,
						Value:          Tank.Value,
						/*
							X:              decodedWp.X,
							Y:              decodedWp.Y,							XSpeed:         decodedWp.XSpeed,
							YSpeed:         decodedWp.YSpeed,*/
					}
					lockPacket := lockPack.Serialize(56, true)
					packet, err := enet.NewPacket(lockPacket, enet.PacketFlagReliable)
					if err != nil {
						log.Fatal("Error packet 15: %s", err.Error())
					}
					for _, currentPeer := range GetPeers(PlayerMap) {
						if NotSafePlayer(peer) {
							continue
						}
						if PInfo(peer).CurrentWorld == PInfo(currentPeer).CurrentWorld {
							fn.PlayMsg(peer, 0, "audio/use_lock.wav")
							fn.TalkBubble(currentPeer, PInfo(peer).NetID, 100, true, "`5[`w%s ``has been `$World Locked ``by `2%s`5]``", PInfo(peer).CurrentWorld, strings.TrimSpace(PInfo(peer).Name))
							currentPeer.SendPacket(packet, 0)
						}
					}
					fn.ModifyInventory(peer, int(Tank.Value), -1, PInfo(peer))
					fn.AddTile(peer, Tank)
					worlds.Worlds[PInfo(peer).CurrentWorld].Tiles[Coords].Fg = int16(Tank.Value)
					worlds.Worlds[PInfo(peer).CurrentWorld].OwnerUid = int32(PInfo(peer).UserID)
					break
				} else {
					fn.TalkBubble(peer, PInfo(peer).NetID, 0, false, "Someone has been used locks!")
					break
				}
				//fn.UpdateInventory(peer)
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
				//go worlds.SaveWorld(PInfo(peer).CurrentWorld, *world)
				fn.TalkBubble(peer, PInfo(peer).NetID, 100, false, "Updating Block at %d", Coords)
			}
			break
		}
	default:
		{
			// fn.OnPlace(peer, Tank)
			decodedPack := &tankpacket.TankPacket{}
			decodedPack.SerializeFromMem(packet.GetData()[4:])
			log.Info("Got Unknown Packet with type %d: %#v", decodedPack.PacketType, decodedPack)
			break
		}

	}
	return

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
	lowerCmd = strings.TrimRight(lowerCmd, "\x00") // Trim any null character
	if isCommand {
		fn.LogMsg(peer, "`6"+cmd)
	}
	if strings.HasPrefix(lowerCmd, "/help") {
		fn.LogMsg(peer, "Help Command >> /help /ip")
	} else if strings.HasPrefix(lowerCmd, "/myip") {
		fn.LogMsg(peer, "Your IP Address: %s", GetPlayer(peer).IpAddress)
		/*else if strings.HasPrefix(lowerCmd, "/finditem") {
		errorMessage := ">> Usage: /finditem <`$item name``> - Searches for item."
		a_ := strings.Split(cmd, " ")
		if len(a_) <= 1 {
			fn.ConsoleMsg(peer, 0, errorMessage)
			return
		}
		if len(a_) >= 2 {
			a_ = a_[1:]
			targetSifinds := strings.ToLower(strings.Join(a_, " "))
			if len(targetSifinds) < 3 {
				fn.ConsoleMsg(peer, 0, ">> Enter atleast 3 letters!")
				return
			}
			searchList := ""

			for i := 0; i < len(items.); i++ {
				itemID := items[i].id
				if pInfo(peer).admin {
					if items[itemID].blockType == BlockTypes.SEED {
						continue
					}
					if !strings.Contains(items[itemID].oriName, "null_item") || !strings.Contains(items[itemID].oriName, "Guild Flag") || !strings.Contains(items[itemID].oriName, "Kranken") || !strings.Contains(items[itemID].oriName, "Sacrificial Well") {
						continue
					}
					if strings.Contains(strings.ToLower(items[i].oriName), targetSifinds) {
						searchList += "\nadd_button_with_icon|search_" + strconv.Itoa(int(itemID)) + "|`$" + items[i].oriName + "`2(" + strconv.Itoa(int(itemID)) + "`0)|staticYellowFrame | " + strconv.Itoa(int(itemID)) + " || "
					}
				}
				if pInfo(peer).supermod {
					if items[itemID].blockType == BlockTypes.SEED || items[itemID].blockType == BlockTypes.LOCK {
						continue
					}
					if !strings.Contains(items[itemID].oriName, "chest") || !strings.Contains(items[itemID].oriName, "legend") || !strings.Contains(items[itemID].oriName, "null_item") || !strings.Contains(items[itemID].oriName, "Guild Flag") || !strings.Contains(items[itemID].oriName, "Kranken") || !strings.Contains(items[itemID].oriName, "Sacrificial Well") {
						continue
					}
					if strings.Contains(strings.ToLower(items[i].oriName), targetSifinds) {
						searchList += "\nadd_button_with_icon|search_" + strconv.Itoa(int(itemID)) + "|`$" + items[i].oriName + "`2(" + strconv.Itoa(int(itemID)) + "`0)|staticYellowFrame | " + strconv.Itoa(int(itemID)) + " || "
					}
				}
			}
			if searchList == "" {
				packet_(peer, "action|log\nmsg| `4Oops: `oThere is no items found starting with `w"+targetSifinds+"`o.", "")
				return
			}
			p := gamepacket.NewPacket()
			p.Insert("OnDialogRequest")
			p.Insert("add_label_with_icon|big|`wSearch results for ``\"" + targetSifinds + "\"``|left|6016|\nadd_spacer|small|\nembed_data|search|" + targetSifinds + "\nend_dialog|search_option|Cancel|\nadd_spacer|big|\n" + searchList + "add_quick_exit|\n")
			p.CreatePacket(peer)
		}*/
	} else if strings.HasPrefix(lowerCmd, "/sb") {
		parsedSb := strings.Fields(lowerCmd)
		if len(parsedSb) != 2 {
			fn.ConsoleMsg(peer, 0, "Usage: /sb teks")
			return
		}
		for _, currentPeer := range GetPeers(PlayerMap) {
			if NotSafePlayer(currentPeer) {
				continue
			}
			fn.ConsoleMsg(currentPeer, 0, "CP:_PL:0_OID:_CT:[SB]_ `5** `1Super Broadcast`` from (`0%s`5) in [`$%s`5] ** : %s", GetPlayerName(peer), PInfo(peer).CurrentWorld, GetChatPrefix(peer)+parsedSb[1])
			fn.PlayMsg(currentPeer, 0, "audio/beep.wav")
		}
	} else if strings.HasPrefix(lowerCmd, "/give") {
		parsedGive := strings.Fields(lowerCmd)
		if len(parsedGive) != 3 {
			fn.ConsoleMsg(peer, 0, "Usage: /give ItemID Quantity")
		}
		log.Info("Give Cmd: %v, 0: %s, 1: %s, 2: %s", parsedGive, parsedGive[0], parsedGive[1], parsedGive[2])
		itemId, err := strconv.Atoi(parsedGive[1])
		if err != nil {
			fn.ConsoleMsg(peer, 0, err.Error())
			return
		}
		itemQty, err := strconv.Atoi(parsedGive[2])
		if err != nil {
			fn.ConsoleMsg(peer, 0, err.Error())
			return
		}
		inventory := append(PInfo(peer).Inventory, ItemInfo{ID: itemId, Qty: int16(itemQty)})
		PInfo(peer).Inventory = inventory
		fn.UpdateInventory(peer)
		fn.ConsoleMsg(peer, 0, "Successfully added `w%d`` (%d) to Inventory", itemId, itemQty)
	} else if strings.HasPrefix(lowerCmd, "/info") {
		cpuUsage := runtime.NumCPU()
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		memoryUsage := m.Sys / 1024 / 1024
		fn.ConsoleMsg(peer, 0, "CPU: %d Core(s)\nAlloc: %d MB\nTotalAlloc: %d MB\nSys: %d MB\nNumGC: %d Thread(s)", cpuUsage, m.Alloc/1024/1024, m.TotalAlloc/1024/1024, memoryUsage, m.NumGC)
		fn.TalkBubble(peer, PInfo(peer).NetID, 100, false, "CPU: %d Core(s)\nAlloc: %d MB\nTotalAlloc: %d MB\nSys: %d MB\nNumGC: %d Thread(s)", cpuUsage, m.Alloc/1024/1024, m.TotalAlloc/1024/1024, memoryUsage, m.NumGC)
	} else if strings.HasPrefix(lowerCmd, "/dialog") {
		fn.OnDialogRequest(peer, "text_scaling_string|Dirttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttt|\nset_default_color|`o\nadd_label_with_icon|big|`wGet a growlandID``|left|206|\nadd_spacer|small|\nadd_textbox|By choosing a `wgrowlandID``, you can use a name and password to logon from any device.Your `wname`` will be shown to other players!|left|\nadd_spacer|small|\nadd_text_input|logon|Name|Akbar Awor|18|\nadd_textbox|Your `wpassword`` must contain `w8 to 18 characters, 1 letter, 1 number`` and `w1 special character: @#!$^&*.,``|left|\nadd_text_input_password|password|Password|Anton Malang 69|18|\nadd_text_input_password|password_verify|Password Verify|Anton Malang 69 Verif|18|\nadd_textbox|Your `wemail`` will only be used for account verification and support. If you enter a fake email, you can't verify your account, recover or change your password.|left|\nadd_text_input|email|Email|Akbar Faisal|64|\nadd_textbox|We will never ask you for your password or email, never share it with anyone!|left|\nend_dialog|growid_apply|Cancel|Get My growlandID!|\n", 0)
	} else {
		fn.LogMsg(peer, "`4Unknown command.``  Enter `$/?`` for a list of valid commands.")
	}
	return
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
	return
}

func OnPlayerMove(peer enet.Peer, packet enet.Packet) {

	for _, currentPeer := range GetPeers(PlayerMap) {
		if NotSafePlayer(currentPeer) {
			continue
		}
		if PInfo(peer).CurrentWorld == PInfo(currentPeer).CurrentWorld {
			movePacket := packet.GetData()
			binary.LittleEndian.PutUint16(movePacket[8:], uint16(PInfo(peer).NetID))
			packet, err := enet.NewPacket(movePacket, enet.PacketFlagReliable)
			if err != nil {
				panic(err)
			}
			currentPeer.SendPacket(packet, 0)
		}
	}
	return
}

func OnPlayerExitWorld(peer enet.Peer, world *worlds.World) {
	if NotSafePlayer(peer) {
		return
	}
	/*if PInfo(peer).CurrentWorld == "" {
		return
	}*/
	/*world, err := worlds.LoadWorld(PInfo(peer).CurrentWorld)
	if err != nil {
		peer.DisconnectLater(0)
		log.Error("[OnPlayerExitWorld] Worlds with name: %s is not found in our database!", world.Name)
	}*/
	for _, currentPeer := range GetPeers(PlayerMap) {
		if NotSafePlayer(currentPeer) {
			continue
		}
		if PInfo(peer).CurrentWorld == PInfo(currentPeer).CurrentWorld && PInfo(currentPeer).PeerID != PInfo(peer).PeerID {
			fn.PlayMsg(currentPeer, 0, "audio/door_shut.wav")
			fn.OnRemove(currentPeer, int(PInfo(peer).NetID))
			fn.TalkBubble(currentPeer, PInfo(peer).NetID, 0, true, "`5<`0%s`` left, `w%d`5 others here>``", GetPlayerName(peer), world.PlayersIn)
			fn.ConsoleMsg(currentPeer, 0, "`5<`0%s`` left, `w%d`5 others here>``", GetPlayerName(peer), world.PlayersIn)
		}
	}
	PInfo(peer).SpawnX = 0
	PInfo(peer).SpawnY = 0
	if PInfo(peer).CurrentWorld != "" && (world.Name != "" || world.Name != "EXIT") {
		if world.PlayersIn < 1 {
			world.PlayersIn = 0
			//worlds.UpsertWorld(world.Name)
			//delete(worlds.Worlds, world.Name)
		}
		world.PlayersIn--
		worlds.Worlds[world.Name].PlayersIn = world.PlayersIn
		fn.ListActiveWorld[PInfo(peer).CurrentWorld] = int(world.PlayersIn)
	}
	PInfo(peer).CurrentWorld = ""
	fn.UpdateInventory(peer)
	fn.SendWorldMenu(peer)
	//codedWorld := worlds.AutoTagMsgpackStruct(world)
	/*if PInfo(peer).TankIDName != "" {
		UpsertPlayer(peer, PInfo(peer).TankIDName)
	} else {
		UpsertPlayer(peer, PInfo(peer).Rid)
	}
	log.Info("Saving World with name: %s", world.Name)*/
	return
}

func OnConnect(peer enet.Peer, host enet.Host, items *items.ItemInfo, globalPeer []enet.Peer) {
	log.Info("New Client Connected %s", peer.GetAddress().String())
	/*PlayerMap[peer] = &Player{
		IpAddress: peer.GetAddress().String(),
		//Roles:     6,
	}
	/*PlayerConnect := &Player{
		IpAddress: peer.GetAddress().String(),
		Peer:      peer,
		Roles:     6,
	}
	SavePlayer(*PlayerConnect)*/
	pkt.SendPacket(peer, 1, "") //hello response
	return
}

func OnDisConnect(peer enet.Peer, host enet.Host, items *items.ItemInfo, globalPeer []enet.Peer) {
	log.Info("New Client Disconnected %s", peer.GetAddress().String())
	if PInfo(peer).CurrentWorld != "" {
		OnPlayerExitWorld(peer, worlds.Worlds[PInfo(peer).CurrentWorld])
	}
	/*PlayerMapBackup := PlayerMap
	if NotSafePlayer(peer) {
		return
	} else {
		currentWorld := PInfo(peer).CurrentWorld
		fn.ListActiveWorld[currentWorld] = int(worlds.Worlds[currentWorld].PlayersIn)
		if currentWorld != "" || currentWorld != "EXIT" {
			for _, currentPeer := range GetPeers(PlayerMapBackup) {
				if NotSafePlayer(currentPeer) {
					return
				}
				if PInfo(currentPeer).CurrentWorld == currentWorld {
					fn.OnRemove(currentPeer, int(PInfo(peer).NetID))
					fn.TalkBubble(currentPeer, PInfo(peer).NetID, 0, true, "`5<`0%s`` left, `w%d`5 others here>``", GetPlayerName(peer), worlds.Worlds[currentWorld].PlayersIn)
				}
			}
			/*if PInfo(peer).TankIDName != "" {
				UpsertPlayer(peer, PInfo(peer).TankIDName)
			} else {
				UpsertPlayer(peer, PInfo(peer).Rid)
			}
			delete(PlayerMap, peer)
		}
	}*/
	return
}

func OnTextPacket(peer enet.Peer, host enet.Host, text string, items *items.ItemInfo, globalPeer []enet.Peer) {
	log.Info("TextPacket: %s", text)
	if strings.Contains(text, "requestedName|") {
		fn.OnSuperMain(peer, items.GetItemHash())
		ParseUserData(text, peer)
	} else if len(text) > 6 && text[:6] == "action" {
		lengthText := 7 + len(strings.Split(text[7:], "\n")[0])
		switch text[7:lengthText] {
		case "enter_game":
			{
				fn.UpdateName(peer, PInfo(peer).Name)
				if PInfo(peer).TankIDName != "" {
					fn.SetHasGrowID(peer)
					// fn.SetAccountHasSecured(peer)
				}
				log.Info("Loaded Skin: %d", PInfo(peer).SkinColor)
				fn.UpdateClothes(0, peer)
				fn.TextOverlay(peer, "`2Welcome To GotPS!``")
				fn.SendWorldMenu(peer)
				// NewPlayer(peer)
				fn.LogMsg(peer, "Where would you like to go? (`w%d`` Online)", host.ConnectedPeers())
				break
			}
		case "join_request":
			{
				log.Info("Invent Size: %d", byte(PInfo(peer).InventorySize))
				log.Info("%s", text[7:])
				fn.UpdateInventory(peer)
				worldName := strings.ToUpper(strings.Split(text[25:], "\n")[0])
				fn.LogMsg(peer, "Sending you to world (%s) (%d)", worldName, len(worldName))
				OnEnterGameWorld(peer, host, worldName)
				break
			}
		case "input":
			{
				text := strings.Split(text[19:], "\n")[0]
				if text[0] == '/' {
					OnCommand(peer, host, text, true)
				} else {
					OnChatInput(peer, host, text)
				}
				break
			}
		case "quit_to_exit":
			{
				/*fn.ListActiveWorld[PInfo(peer).CurrentWorld]--
				PInfo(peer).CurrentWorld = ""*/
				OnPlayerExitWorld(peer, worlds.Worlds[PInfo(peer).CurrentWorld])
				break
			}
		case "quit":
			{
				peer.DisconnectLater(0)
				break
			}
		case "refresh_item_data":
			{
				fn.LogMsg(peer, "One moment, Updating item data...")
				packet, err := enet.NewPacket(items.FileBufferPacket, enet.PacketFlagReliable)
				if err != nil {
					panic(err)
				}
				peer.SendPacket(packet, 0)
				break
			}
		case "setSkin":
			{
				splitSkin := strings.Split(strings.Split(text[7:], "\n")[1], "|")[1]
				log.Info("Skin ID: %s", splitSkin)
				skinId, err := strconv.Atoi(splitSkin)
				if err != nil {
					fn.ConsoleMsg(peer, 0, "Error when trying to change your skin color to %d!", skinId)
					return
				}
				PInfo(peer).SkinColor = skinId
				for _, currentPeer := range GetPeers(PlayerMap) {
					if NotSafePlayer(currentPeer) {
						return
					}
					if PInfo(peer).CurrentWorld == PInfo(currentPeer).CurrentWorld {
						fn.UpdateClothes(0, currentPeer)
					}
				}
				break
			}
		default:
			{
				if strings.HasPrefix(text[7:], "quit_to_exit") {
					OnPlayerExitWorld(peer, worlds.Worlds[PInfo(peer).CurrentWorld])
				} else {
					//fn.LogMsg(peer, "Unhandled Action Packet type: %s", text[7:])
				}
				break
			}
		}
	} else {
		fn.LogMsg(peer, "Unhandled TextPacket, msg: %v", text)
		log.Info("Unhandled TextPacket, msg: %v", text)
	}
	return
}

func OnTankPacket(peer enet.Peer, host enet.Host, packet enet.Packet, items *items.ItemInfo, globalPeer []enet.Peer) {
	if NotSafePlayer(peer) {
		return
	}
	if len(packet.GetData()) < 3 {
		return
	}

	/*world, err := worlds.LoadWorld(PInfo(peer).CurrentWorld)
	if err != nil {
		fn.ConsoleMsg(peer, 0, "Some data is missing in this world! exiting..")
		OnPlayerExitWorld(peer, world)
	}*/
	if PInfo(peer).CurrentWorld != "" {
		world, err := worlds.GetWorld(PInfo(peer).CurrentWorld)
		if err != nil {
			log.Fatal(err.Error())
		}
		//log.Info("[OnTankPacket] Successfully Loaded World from Memory:%v", world)
		Tank := &tankpacket.TankPacket{}
		Tank.SerializeFromMem(packet.GetData()[4:])

		switch Tank.PacketType {
		case 0:
			{ //player movement
				//fn.LogMsg(peer, "[Movement] X:%f, Y:%f", Tank.X, Tank.Y)
				OnPlayerMove(peer, packet)
				break
			}
		case 3:
			{
				OnTileUpdate(packet, peer, Tank, world)
			}
		case 7:
			{
				// Door
				if GetPlayer(peer).CurrentWorld != "" {
					OnPlayerExitWorld(peer, world)
				}
				break
			}
		case 10:
			{
				switch Tank.Value {
				case 5480:
					{
						PInfo(peer).Clothes.Hand = float32(Tank.Value)
						for _, currentPeer := range GetPeers(PlayerMap) {
							fn.UpdateClothes(0, currentPeer)
						}
						break
					}
				}
				log.Info("Packet type: %d, val: %d (%#v)", Tank.PacketType, Tank.Value, Tank)
				break
			}
		case 22:
			{
				pkt.SendPacket(peer, 21, "")
				break
			}
		default:
			{
				//PInfo(peer).NetID = Tank.NetID
				log.Info("Packet type: %d, val: %d, struct: %#v", Tank.PacketType, Tank.Value, Tank)
				break
			}
		}
	}
	return
}

func OnEnterGameWorld(peer enet.Peer, host enet.Host, name string) {
	/*log.Info("[OnEnterGameWorld] Player Data: %v", PInfo(peer))
	if NotSafePlayer(peer) {
		fn.LogMsg(peer, "`4Invalid Player Data!``")
		return nil
	}
	world := worlds.Worlds[name]
	if world == nil {
		/*var err error
		world, err = worlds.LoadWorld(name)
		if err != nil {
		worlds.Worlds[name] = worlds.GenerateWorld(name, 100, 60)
		//codedWorld := AutoTagMsgpackStruct(world)
		world = worlds.Worlds[name]
		//worlds.UpsertWorld(name)
		log.Error("Worlds with name: %s is not found in our database!", name)
		//}
	}
	/*nameLen := len(world.Name)
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
				/*worldPacket[extraDataPos+8] = 3 //block types
				binary.LittleEndian.PutUint16(worldPacket[extraDataPos+9:], uint16(len(world.Tiles[i].Label)))
				copy(worldPacket[extraDataPos+11:], []byte(world.Tiles[i].Label))

				//SpawnX = (i % int(world.SizeX)) * 32
				// SpawnY = (i / int(world.SizeX)) * 32
				extraDataPos += 4 + len(world.Tiles[i].Label)
				worldPacket[extraDataPos+8] = 3
				//musicBpm := 100
				binary.LittleEndian.PutUint16(worldPacket[extraDataPos+2:], uint16(PInfo(peer).UserID))
				/*binary.LittleEndian.PutUint16(worldPacket[extraDataPos+18:], uint16(musicBpm))
				extraDataPos += 4

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
	peer.SendPacket(packet, 0)*/
	fn.UpdateWorld(peer, name)
	world := worlds.Worlds[name]
	SpawnX = int(PInfo(peer).SpawnX)
	SpawnY = int(PInfo(peer).SpawnY)
	// Avoid minus
	if int(world.PlayersIn) < 1 {
		world.PlayersIn = 0
	}
	world.PlayersIn++
	worlds.Worlds[world.Name].PlayersIn = world.PlayersIn
	// Simple Fix
	PInfo(peer).NetID = uint32(world.PlayersIn)
	// BotList.Load(peer)
	if PInfo(peer).CurrentWorld == "" {
		fn.OnSpawn(peer, int16(PInfo(peer).NetID), PInfo(peer).PeerID, int32(PInfo(peer).SpawnX), int32(PInfo(peer).SpawnY), GetPlayerName(peer), PInfo(peer).Country, false, true, true, true)
	}
	PInfo(peer).CurrentWorld = world.Name
	fn.PlayMsg(peer, 0, "audio/door_open.wav")
	fn.ConsoleMsg(peer, 0, "`5<`w%s ``entered, `w%d`` others here`5>", GetPlayerName(peer), world.PlayersIn)
	fn.TalkBubble(peer, PInfo(peer).NetID, 300, true, "`5<`w%s ``entered, `w%d`` others here`5>", GetPlayerName(peer), world.PlayersIn)
	for _, currentPeer := range GetPeers(PlayerMap) {
		log.Info("%v", currentPeer)
		if PInfo(currentPeer).CurrentWorld == PInfo(peer).CurrentWorld && PInfo(currentPeer).PeerID != PInfo(peer).PeerID {
			//fn.OnSpawn(peer, world.PlayersIn, world.PlayersIn, int32(SpawnX), int32(SpawnY), GetPlayerName(currentPeer), PInfo(currentPeer).Country, false, true, true, isLocal)
			// Spawn Another Player Avatar to You
			fn.PlayMsg(currentPeer, 0, "audio/door_open.wav")
			fn.ConsoleMsg(currentPeer, 0, "`5<`w%s ``entered, `w%d`` others here`5>", GetPlayerName(peer), world.PlayersIn)
			fn.TalkBubble(currentPeer, PInfo(peer).NetID, 300, true, "`5<`w%s ``entered, `w%d`` others here`5>", GetPlayerName(peer), world.PlayersIn)
			fn.OnSpawn(currentPeer, int16(PInfo(peer).NetID), PInfo(peer).PeerID, int32(PInfo(peer).SpawnX), int32(PInfo(peer).SpawnY), GetPlayerName(peer), PInfo(peer).Country, false, true, true, false)
			fn.OnSpawn(peer, int16(PInfo(currentPeer).NetID), PInfo(currentPeer).PeerID, int32(PInfo(currentPeer).SpawnX), int32(PInfo(currentPeer).SpawnY), GetPlayerName(currentPeer), PInfo(currentPeer).Country, false, true, true, false)

		}
		break
	}
	fn.ListActiveWorld[world.Name] = int(world.PlayersIn)
	return
}
