package protos

import (
	//"bytes"
	
	"log"
)

import (
	//"../helper"
	. "../types"
)

type user_login_info_req struct {
	F_client_version int32
	F_new_user       bool
	F_user_name      string
}

type user_login_info_ack struct {
	Name string
}

// Object Type
type Object map[string]interface{}

func P_heart_beat_req(sess *Session, tbl []interface{}) []interface{} {
	// nothing should be done
	return nil
}

func P_user_login_req(sess *Session, tbl []interface{}) []interface{} {

	log.Printf("get: (%v)\n", tbl)

	ack := []interface{}{1, 222, "ddd", "中午", "111", 55, []int{1, 2, 3}, []string{"a", "b", "c"}, []interface{}{"a", 1, "c"}}


	return ack

}
