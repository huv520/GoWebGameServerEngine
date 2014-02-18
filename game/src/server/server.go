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

func start() {
	log.Println("Starting the server.")
	// Listen
	service := ":8080"
	log.Println("Service:", service)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	log.Println("Game Server OK.")

	for {
		conn, err := listener.Accept()

		if err != nil {
			continue
		}

		// DoS prevention  流量攻击预防
		IP := net.ParseIP(conn.RemoteAddr().String())
		if !utils.IsBanned(IP) {
			go handleClient(conn)
		} else {
			conn.Close()
		}
	}
}

//----------------------------------------------- start a goroutine when a new connection is accepted
func handleClient(conn net.Conn) {
	defer conn.Close()

	header := make([]byte, 2)
	ch := make(chan []byte, 10)

	go StartAgent(ch, conn)

	for {
		// header
		n, err := io.ReadFull(conn, header)
		if n == 0 && err == io.EOF {
			break
		} else if err != nil {
			log.Println("error receiving header:", err)
			break
		}

		// data
		size := binary.BigEndian.Uint16(header)
		//log.Println("Receive string:", string(header))
		if size > 1024 {
			log.Println("error size:", size)
			break
		}

		log.Printf("Receive Size: (%v)\n", size)

		data := make([]byte, size)
		n, err = io.ReadFull(conn, data)

		log.Printf("Receive Data: (%v)\n", data)

		if err != nil {
			log.Println("error receiving msg:", err)
			break
		}

		ch <- data
	}

	close(ch)
}

func checkError(err error) {
	if err != nil {
		log.Println("Fatal error: %v", err)
		os.Exit(-1)
	}
}

func init() {
	log.SetPrefix("[GS]")
}
