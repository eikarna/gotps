package functions

import (
	"fmt"
	"strconv"

	"github.com/codecat/go-libs/log"
	enet "github.com/eikarna/gotops"
	variant "github.com/eikarna/gotps/functions/variants"
	pkt "github.com/eikarna/gotps/packet"
	"github.com/eikarna/gotps/utils"
)

func SendWorldMenu(peer enet.Peer) {
	var world_packet string
	world_packet += "add_filter|\n"
	world_packet += "add_heading|Private Source<ROW2>|\n"
	world_packet += "add_floater|GOLANG|GOLANGPS|0|0.5|3529161471\n"
	world_packet += "add_heading|Credits<CR>|\n"
	world_packet += "add_floater|KIPASGTS|KIPASGTS|0|0.5|2147418367\n"
	world_packet += "add_floater|EIKARNA|EIKARNA|0|0.5|2147418367\n"
	world_packet += "add_floater|AKBARDEV|AKBARDEV|0|0.5|2147418367\n"
	world_packet += "add_floater|TEAMNEVOLUTION|TEAMNEVOLUTION|0|0.5|2147418367\n"

	variant := variant.NewVariant(0, -1) //delay netid
	variant.InsertString("OnRequestWorldSelectMenu")
	variant.InsertString(world_packet)
	variant.Send(peer)
}

func OnSuperMain(peer enet.Peer, itemHash uint32) {

	variant := variant.NewVariant(0, -1)
	variant.InsertString("OnSuperMainStartAcceptLogonHrdxs47254722215a")
	variant.InsertUnsignedInt(itemHash)
	variant.InsertString("ubistatic-a.akamaihd.net")
	variant.InsertString("0098/0704202400/cache/")
	variant.InsertString("cc.cz.madkite.freedom org.aqua.gg idv.aqua.bulldog com.cih.gamecih2 com.cih.gamecih com.cih.game_cih cn.maocai.gamekiller com.gmd.speedtime org.dax.attack com.x0.strai.frep com.x0.strai.free org.cheatengine.cegui org.sbtools.gamehack com.skgames.traffikrider org.sbtoods.gamehaca com.skype.ralder org.cheatengine.cegui.xx.multi1458919170111 com.prohiro.macro me.autotouch.autotouch com.cygery.repetitouch.free com.cygery.repetitouch.pro com.proziro.zacro com.slash.gamebuster")
	variant.InsertString("proto=206|choosemusic=audio/mp3/nusaverse_lobby.mp3|active_holiday=19|wing_week_day=0|ubi_week_day=2|server_tick=123665344|clash_active=0|drop_lavacheck_faster=1|isPayingUser=2|usingStoreNavigation=1|enableInventoryTab=1|bigBackpack=1|m_clientBits=0|eventButtons={\"EventButtonData\":[{\"Components\":[{\"Enabled\":false,\"Id\":\"Overlay\",\"Parameters\":\"target_child_entity_name:overlay_layer;var_name:alpha;target:0;interpolation:1;on_finish:1;duration_ms:1000;delayBeforeStartMS:1000\",\"Type\":\"InterpolateComponent\"}],\"DialogName\":\"openLnySparksPopup\",\"IsActive\":false,\"Name\":\"LnyButton\",\"Priority\":1,\"Text\":\"0/5\",\"TextOffset\":\"0.01,0.2\",\"Texture\":\"interface/large/event_button3.rttex\",\"TextureCoordinates\":\"0,2\"},{\"Components\":[{\"Enabled\":true,\"Parameters\":\"\",\"Type\":\"RenderDailyChallengeComponent\"}],\"DialogName\":\"dailychallengemenu\",\"IsActive\":false,\"Name\":\"DailyChallenge\",\"Priority\":2},{\"Components\":[{\"Enabled\":false,\"Id\":\"Overlay\",\"Parameters\":\"target_child_entity_name:overlay_layer;var_name:alpha;target:0;interpolation:1;on_finish:1;duration_ms:1000;delayBeforeStartMS:1000\",\"Type\":\"InterpolateComponent\"}],\"DialogName\":\"openStPatrickPiggyBank\",\"IsActive\":false,\"Name\":\"StPatrickPBButton\",\"Priority\":1,\"Text\":\"0/0\",\"TextOffset\":\"0.00,0.05\",\"Texture\":\"interface/large/event_button4.rttex\",\"TextureCoordinates\":\"0,0\"},{\"DialogName\":\"show_bingo_ui\",\"IsActive\":false,\"Name\":\"Bingo_Button\",\"Priority\":1,\"Texture\":\"interface/large/event_button4.rttex\"}]}")
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

	log.Info(spawnAvatar)
}

//variants
