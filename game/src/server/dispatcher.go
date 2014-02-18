package server

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"time"
)

import (
	"github.com/bitly/go-simplejson"
	"protos"
	"utils"
)

func call(player *Player, js *simplejson.Json) []byte {
	defer utils.PrintPanicStack()
	now := time.Now()

	log.Println("now:", now)

	// read protocol id
	pidByte := p[:2]
	b := binary.BigEndian.Uint16(pidByte)

	log.Printf("protocol id:%v\n", b)
	handle := protos.Handler[b]
	//log.Printf("handle:%v\n", handle)
	if handle != nil {

		tbl := []interface{}{}

		err := json.Unmarshal(p[2:], &tbl)
		if err != nil {
			player.KickOut = true
			log.Printf("json error: (%s)\n", err.Error())
		}

		log.Printf("read data: (%v)\n", tbl)

		ret := handle(player, tbl)
		if len(ret) != 0 {

			outBuf, err := json.Marshal(ret)
			if err != nil {
				player.KickOut = true
				log.Printf("json error: (%s)\n", err.Error())
			}

			log.Printf("write data: (%v)\n", outBuf)

			return append(pidByte, outBuf...)
		}
	}

	return nil
}
