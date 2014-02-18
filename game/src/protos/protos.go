package protos

import (
    "log"
    "github.com/bitly/go-simplejson"
	"server"
)

type commandHandler func(player *Player, param *simplejson.Json)

// handler map for "Command"
var cmHandlers = map[string]commandHandler{
    //"CM_REGISTER" : cmRegisterHander,
  //  "CM_LOGIN" : cmLoginHander,
   // "CM_CHAR_CREATE" : cmCharCreateHander,
   // "CM_CHAR_GET" : cmCharGetHander,
   // "CM_CARDS_GET" : cmCardsGetHander,
   // "CM_RAID" : cmRaidHander,
}