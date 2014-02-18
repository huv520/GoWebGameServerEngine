package server

import (
	"encoding/binary"
	"net"
	"time"
)

import (
	"utils"
)

const (
	DEFAULT_BAN_TIME = 5
)

//---------------------------------------------------------- IP->UnBan time
var _banned_ips map[uint32]int64

func init() {
	_banned_ips = make(map[uint32]int64)
}

//---------------------------------------------------------- ban an ip
func Ban(_ip net.IP) {
	ban_time := DEFAULT_BAN_TIME

	intip := utils.Ip2Uint32(_ip)

	// randomize the timeout, for effective DoS protection
	ban := uint32(ban_time)
	_banned_ips[intip] = time.Now().Unix() + int64(ban+utils.LCG()%ban)
}

//---------------------------------------------------------- test whether the ip is banned
func IsBanned(_ip net.IP) bool {
	intip := utils.Ip2Uint32(_ip)
	timeout, exists := _banned_ips[intip]

	if !exists {
		return false
	} else if timeout < time.Now().Unix() {
		delete(_banned_ips, intip)
		return false
	} else {
		return true
	}
}
