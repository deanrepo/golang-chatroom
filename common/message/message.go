package message

import (
	"projects/chatroom1/common/model"
)

const (
	RegisterMsgType             = "RegisterMsg"
	RegisterResMsgType          = "RegisterResMsg"
	LoginMsgType                = "LoginMsg"
	LoginResMsgType             = "LoginResMsg"
	NoticeMsgType               = "NoticeMsg"
	LogoutMsgType               = "LogoutMsg"
	OnlineUserListReqMsgType    = "OnlineUserListReqMsg"
	OnlineUserListReqResMsgType = "OnlineUserListReqResMsg"
	GroupSmsMsgType             = "GroupSmsMsg"
	SmsMsgType                  = "SmsMsg"
)

type Message struct {
	Type string `json:"type"` // message type
	Data string `json:"data"` // message data
}

type RegisterMsg struct {
	model.User
}

type RegisterResMsg struct {
	Code  int    `json:"code"` // result status code, 100 represents success, 200 represents fail.
	Error string `json:"error"`
}

type LoginMsg struct {
	model.User
}

type LoginResMsg struct {
	model.User
	Code  int    `json:"code"` // result status code, 300 represents success, 400 represents fail.
	Error string `json:"error"`
}

type NoticeMsg struct {
	Content string `json:"content"`
}

/* type LogoutMsg struct {

} */

/* type OnlineUserListReqMsg struct {

} */

type OnlineUserListReqResMsg struct {
	OnlineUsers []model.User `json:onlineUsers`
}

type GroupSmsMsg struct {
	Content string `json:"content"`
}

type SmsMsg struct {
	UserID  int64  `json:""userID`
	Content string `json:"content"`
}
