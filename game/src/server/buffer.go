package server

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
)


type Buffer struct {
	ctrl    chan bool   // receive exit signal
	pending chan []byte // pending Packet
	max     int         // max queue size
	conn    net.Conn    // connection
	player  Player     // player
}

const (
	DEFAULT_QUEUE_SIZE = 15
)

//------------------------------------------------ send packet
func (buf *Buffer) Send(data []byte) (err error) {
	// len of Channel: the number of elements queued (un-sent) in the channel buffer
	if len(buf.pending) < buf.max {
		buf.pending <- data
		return nil
	} else {
		Ban(buf.player.IP)
		return errors.New(fmt.Sprintf("Send Buffer Overflow, possible DoS attack. Remote: %v", buf.conn.RemoteAddr()))
	}
}

//------------------------------------------------ packet sender goroutine
func (buf *Buffer) Start() {
	defer func() {
		recover()
	}()

	for {
		select {
		case data := <-buf.pending:
			buf.raw_send(data)
		case _, ok := <-buf.ctrl:
			if !ok {
				return
			}
		}
	}
}

//------------------------------------------------ send packet online
func (buf *Buffer) raw_send(data []byte) {
	size := uint16(2 + len(data))
	log.Printf("conn out size: (%v)", size)

	pid := make([]byte, 2)
	binary.BigEndian.PutUint16(pid, size)

	out := append(pid, data...)

	log.Printf("conn out data: (%v)\n", out)

	_, err := buf.conn.Write(out)
	if err != nil {
		log.Println("Error send reply :", err)
		return
	}
}

//------------------------------------------------ create a new write buffer
func NewBuffer(player *Player, conn net.Conn, ctrl chan bool) *Buffer {
	max := DEFAULT_QUEUE_SIZE

	buf := Buffer{conn: conn}
	buf.pending = make(chan []byte, max)
	buf.ctrl = ctrl
	buf.max = max
	return &buf
}
