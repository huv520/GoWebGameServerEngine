package main

import (
	"log"
	"server"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	log.Println("main")
	server.start()

}
