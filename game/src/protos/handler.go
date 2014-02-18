package protos

import (
	. "../types"
)

var Handler map[uint16]func(*Session, []interface{}) []interface{} = map[uint16]func(*Session, []interface{}) []interface{}{
	1000: P_heart_beat_req,
	1001: P_user_login_req,
	1100: P_key_exchange_req,
}
