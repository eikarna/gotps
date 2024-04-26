package functions

import (
	"github.com/eikarna/gotops"
	pkt "github.com/eikarna/gotps/packet"
)

func SendLogonFail(p enet.Peer) {
	pkt.SendPacket(p, 3, "action|play_sfx\nfile|audio/piano_nice.wav\ndelayMS|0\n")
	pkt.SendPacket(p, 3, "action|log\nmsg|`2Hello World from gotops!")
	pkt.SendPacket(p, 3, "action|set_url\nurl|https://github.com/eikarna/gotops\nlabel|GOTOPS Repo")
	pkt.SendPacket(p, 3, "action|logon_fail\n")
}
