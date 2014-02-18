package protos


type user_login_info_req struct {
	F_client_version int32
	F_new_user       bool
	F_user_name      string
}

type user_login_info_ack struct {
	Name string
}


