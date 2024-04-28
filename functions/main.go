package functions

import (
	"encoding/binary"
	"fmt"
	"github.com/codecat/go-libs/log"
	enet "github.com/eikarna/gotops"
	variant "github.com/eikarna/gotps/functions/variants"
	pkt "github.com/eikarna/gotps/packet"
	tankpacket "github.com/eikarna/gotps/packet/TankPacket"
	player "github.com/eikarna/gotps/players"
	"github.com/eikarna/gotps/utils"
	"strconv"
)

var ListActiveWorld = make(map[string]int)

// TODO:
func SendItemsData(peer enet.Peer) {

}

func OnRemove(peer enet.Peer, netid int) {
	variant := variant.NewVariant(0, -1)
	variant.InsertString("OnRemove")
	variant.InsertString("netID|" + strconv.Itoa(netid) + "\n")
	variant.Send(peer)
}

func OnPunch(peer enet.Peer, Tank *tankpacket.TankPacket) {
	player.PlayerMap[peer].PunchX = Tank.PunchX
	player.PlayerMap[peer].PunchY = Tank.PunchY
	testt := &tankpacket.TankPacket{
		PacketType:     3,
		NetID:          player.PlayerMap[peer].NetID,
		CharacterState: Tank.CharacterState,
		Value:          Tank.Value,
		X:              player.PlayerMap[peer].PosX,
		Y:              player.PlayerMap[peer].PosY,
		XSpeed:         Tank.XSpeed,
		YSpeed:         Tank.YSpeed,
		PunchX:         player.PlayerMap[peer].PunchX,
		PunchY:         player.PlayerMap[peer].PunchY,
	}
	bbb := testt.Serialize(56, true)
	aaa, err := enet.NewPacket(bbb, enet.PacketFlagReliable)
	if err != nil {
		log.Error("Error Packet 3:", err)
	}
	for _, currentPeer := range player.PlayerMap {
		if player.PlayerMap[peer].CurrentWorld == currentPeer.CurrentWorld {
			currentPeer.Peer.SendPacket(aaa, 0)
		} else {
			continue
		}
	}
	LogMsg(peer, "[Punch/Place] X:%d, Y:%d, Value:%d, NetID:%d", Tank.PunchX, Tank.PunchY, Tank.Value, Tank.NetID)
}

func SendWorldMenu(peer enet.Peer) {
	var world_packet string
	// Add World Start as default
	ListActiveWorld["START"] = 255
	world_packet += "add_filter|\n"
	world_packet += "add_heading|Top Worlds<ROW2>|\n"
	for listworld, count := range ListActiveWorld {
		if count > 0 {
			if listworld == "START" {
				world_packet += "add_floater|" + listworld + "|" + listworld + "|" + strconv.Itoa(count) + "|0.8|3529161471\n"

			} else {
				world_packet += "add_floater|" + listworld + "|" + listworld + "|" + strconv.Itoa(count) + "|0.5|3529161471\n"
			}
		} else {
			delete(ListActiveWorld, listworld)
		}
	}
	world_packet += "add_heading|Credits<CR>|\n"
	world_packet += "add_floater|KIPASGTS|KIPASGTS|0|0.5|2147418367\n"
	world_packet += "add_floater|EIKARNA|EIKARNA|0|0.5|2147418367\n"
	world_packet += "add_floater|AKBARDEV|AKBARDEV|0|0.5|2147418367\n"
	world_packet += "add_floater|TEAMNEVOLUTION|TEAMNEVOLUTION|0|0.5|2147418367\n"
	world_packet += "add_heading|Based On: https://github.com/eikarna/gotops<CR>|\n"

	variant := variant.NewVariant(0, -1) //delay netid
	variant.InsertString("OnRequestWorldSelectMenu")
	variant.InsertString(world_packet)
	variant.Send(peer)
}

func SendInventory(pl player.Player, peer enet.Peer) {
	if player.NotSafePlayer(peer) {
		return
	}
	if len(pl.Inventory) < 1 || pl.InventorySize < 1 {
		//NewInvent := pl.Inventory
		pl.InventorySize = 30
		itemsToAdd := []player.ItemInfo{
			{ID: 18, Qty: 1},
			{ID: 32, Qty: 1},
			{ID: 7188, Qty: 200},
		}
		for _, item := range itemsToAdd {
			pl.Inventory = append(pl.Inventory, item)
		}
	}
	packetLen := 66 + (pl.InventorySize * 4) + 4
	d_ := make([]byte, packetLen)
	binary.LittleEndian.PutUint16(d_[0:], 4)
	binary.LittleEndian.PutUint16(d_[4:], 9)
	binary.LittleEndian.PutUint16(d_[8:], 255)
	binary.LittleEndian.PutUint16(d_[16:], 8)
	binary.LittleEndian.PutUint16(d_[56:], 6+(pl.InventorySize*4)+4)
	binary.LittleEndian.PutUint16(d_[60:], 1)
	binary.LittleEndian.PutUint16(d_[61:], pl.InventorySize)
	binary.LittleEndian.PutUint16(d_[65:], pl.InventorySize)
	offset := 67
	for _, Inven := range pl.Inventory {
		binary.LittleEndian.PutUint16(d_[offset:], uint16(Inven.ID))
		offset += 2
		binary.LittleEndian.PutUint16(d_[offset:], uint16(Inven.Qty))
		offset += 2
	}
	//}
	log.Info("SendInventory Byte: %b | String: %s", d_, d_)
	packet, err := enet.NewPacket(d_, enet.PacketFlagReliable)
	if err != nil {
		log.Error(err.Error())
	}
	peer.SendPacket(packet, 0)
}

/*func SendDoor(peer enet.Peer) {
	if player.GetPlayer(peer).CurrentWorld != "" {
		OnPlayerExitWorld(peer)
	}
	break
}*/

func ConsoleMsg(peer enet.Peer, delay int, a ...interface{}) {
	msg := fmt.Sprintf(a[0].(string), a[1:]...)
	variant := variant.NewVariant(delay, -1)
	variant.InsertString("OnConsoleMessage")
	variant.InsertString(msg)
	variant.Send(peer)
}

func BroadcastConsoleMsg(host enet.Host, a ...interface{}) {
	msg := fmt.Sprintf(a[0].(string), a[1:]...)
	variant := variant.NewVariant(0, -1)
	variant.InsertString("OnConsoleMessage")
	variant.InsertString(msg)
	variant.SendBroadcast(host)
}

func TalkBubble(peer enet.Peer, netID uint32, delay int, isOverlay bool, a ...interface{}) {
	msg := fmt.Sprintf(a[0].(string), a[1:]...)
	variant := variant.NewVariant(delay, -1)
	variant.InsertString("OnTalkBubble")
	variant.InsertUnsignedInt(netID)
	variant.InsertString(msg)
	variant.InsertInt(utils.BoolToInt(isOverlay))
	variant.InsertInt(utils.BoolToInt(isOverlay))
	variant.Send(peer)
}

func BroadcastTalkBubble(host enet.Host, netID uint32, isOverlay bool, a ...interface{}) {
	msg := fmt.Sprintf(a[0].(string), a[1:]...)
	variant := variant.NewVariant(0, -1)
	variant.InsertString("OnTalkBubble")
	variant.InsertUnsignedInt(netID)
	variant.InsertInt(utils.BoolToInt(isOverlay))
	variant.InsertInt(utils.BoolToInt(isOverlay))
	variant.InsertString(msg)
	variant.SendBroadcast(host)
}

func OnSuperMain(peer enet.Peer, itemHash uint32) {

	variant := variant.NewVariant(0, -1)
	variant.InsertString("OnSuperMainStartAcceptLogonHrdxs47254722215a")
	variant.InsertUnsignedInt(itemHash)
	variant.InsertString("ubistatic-a.akamaihd.net")
	variant.InsertString("0098/0704202400/cache/")
	variant.InsertString("cc.cz.madkite.freedom org.aqua.gg idv.aqua.bulldog com.cih.gamecih2 com.cih.gamecih com.cih.game_cih cn.maocai.gamekiller com.gmd.speedtime org.dax.attack com.x0.strai.frep com.x0.strai.free org.cheatengine.cegui org.sbtools.gamehack com.skgames.traffikrider org.sbtoods.gamehaca com.skype.ralder org.cheatengine.cegui.xx.multi1458919170111 com.prohiro.macro me.autotouch.autotouch com.cygery.repetitouch.free com.cygery.repetitouch.pro com.proziro.zacro com.slash.gamebuster")
	variant.InsertString("proto=206|choosemusic=audio/mp3/lobby.mp3|active_holiday=19|wing_week_day=0|ubi_week_day=2|server_tick=123665344|clash_active=0|drop_lavacheck_faster=1|isPayingUser=2|usingStoreNavigation=1|enableInventoryTab=1|bigBackpack=1|m_clientBits=0|eventButtons={\"EventButtonData\":[{\"Components\":[{\"Enabled\":false,\"Id\":\"Overlay\",\"Parameters\":\"target_child_entity_name:overlay_layer;var_name:alpha;target:0;interpolation:1;on_finish:1;duration_ms:1000;delayBeforeStartMS:1000\",\"Type\":\"InterpolateComponent\"}],\"DialogName\":\"openLnySparksPopup\",\"IsActive\":false,\"Name\":\"LnyButton\",\"Priority\":1,\"Text\":\"0/5\",\"TextOffset\":\"0.01,0.2\",\"Texture\":\"interface/large/event_button3.rttex\",\"TextureCoordinates\":\"0,2\"},{\"Components\":[{\"Enabled\":true,\"Parameters\":\"\",\"Type\":\"RenderDailyChallengeComponent\"}],\"DialogName\":\"dailychallengemenu\",\"IsActive\":false,\"Name\":\"DailyChallenge\",\"Priority\":2},{\"Components\":[{\"Enabled\":false,\"Id\":\"Overlay\",\"Parameters\":\"target_child_entity_name:overlay_layer;var_name:alpha;target:0;interpolation:1;on_finish:1;duration_ms:1000;delayBeforeStartMS:1000\",\"Type\":\"InterpolateComponent\"}],\"DialogName\":\"openStPatrickPiggyBank\",\"IsActive\":false,\"Name\":\"StPatrickPBButton\",\"Priority\":1,\"Text\":\"0/0\",\"TextOffset\":\"0.00,0.05\",\"Texture\":\"interface/large/event_button4.rttex\",\"TextureCoordinates\":\"0,0\"},{\"DialogName\":\"show_bingo_ui\",\"IsActive\":false,\"Name\":\"Bingo_Button\",\"Priority\":1,\"Texture\":\"interface/large/event_button4.rttex\"}]}")
	//p.Insert("654171113"); //tribute_data
	variant.Send(peer)
}
func LogMsg(peer enet.Peer, a ...interface{}) {
	msg := fmt.Sprintf(a[0].(string), a[1:]...)
	pkt.SendPacket(peer, 3, "action|log\nmsg|"+msg)
}

func OnSpawn(peer enet.Peer, netid int32, userid int32, posX int32, posY int32, username string, country string, invis bool, mstate bool, smsate bool, local bool) {
	spawnAvatar := "spawn|avatar\n"
	spawnAvatar += "netID|" + strconv.Itoa(int(netid)) + "\n"
	spawnAvatar += "userID|" + strconv.Itoa(int(userid)) + "\n"
	spawnAvatar += "colrect|0|0|20|30\n"
	spawnAvatar += "posXY|" + strconv.Itoa(int(posX)) + "|" + strconv.Itoa(int(posY)) + "\n"
	spawnAvatar += "name|" + username + "\n"
	spawnAvatar += "country|" + country + "\n"
	spawnAvatar += "invis|" + utils.BoolToIntString(invis) + "\n"    //1 0
	spawnAvatar += "mstate|" + utils.BoolToIntString(mstate) + "\n"  //1 0
	spawnAvatar += "smstate|" + utils.BoolToIntString(smsate) + "\n" //1 0
	if local {
		spawnAvatar += "onlineID|\ntype|local\n"
	}

	variant := variant.NewVariant(0, -1)
	variant.InsertString("OnSpawn")
	variant.InsertString(spawnAvatar)
	variant.Send(peer)

	//log.Info(spawnAvatar)
}

//variants
