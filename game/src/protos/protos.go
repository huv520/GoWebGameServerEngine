package protos

import (
    "log"
    "github.com/bitly/go-simplejson"
)

type commandHandler func(this *player, param *simplejson.Json)

// handler map for "Command"
var cmHandlers = map[string]commandHandler{
    "CM_REGISTER" : cmRegisterHander,
    "CM_LOGIN" : cmLoginHander,
    "CM_CHAR_CREATE" : cmCharCreateHander,
    "CM_CHAR_GET" : cmCharGetHander,
    "CM_CARDS_GET" : cmCardsGetHander,
    "CM_RAID" : cmRaidHander,
}

func commandDispatcher(this *player, js *simplejson.Json) {
    rtnCode := 0
    command := ""

    // defer need to be placed at header of func
    defer func() {
        // only send error message here
        if 0 != rtnCode {
            rtnMsg, _ := errCodes[rtnCode]
            rtnJson := responseJson(command, rtnCode, rtnMsg)
            if err := websocket.Message.Send(this.ws, rtnJson); err != nil {
                log.Printf("Send fail for commandDispatcher")
            }
        }
    }()

    command, err := js.Get("Command").String()
    if  err != nil {
        rtnCode = 1
        return
    }

    param, ok := js.CheckGet("Param")
    if ok {
        handler, ok := cmHandlers[command]
        if  ok {
            handler(this, command, param)
            return
        }
    }

    rtnCode = 1
}