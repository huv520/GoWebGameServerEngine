package server

import (
	"net"
	"time"
)

type Session struct {
	IP net.IP

	// session related
	LoggedIn bool // flag for weather the user is logged in
	KickOut  bool // flag for player is kicked out

	// time related
	ConnectTime    time.Time // tcp connection establish time, in millsecond(ms)
	LastPacketTime int64     // last packet arrive time, in seconds(s)
	LastFlushTime  int64     // last flush to db time, in seconds(s)
	OpCount        int       // num of operations since last sync
}
