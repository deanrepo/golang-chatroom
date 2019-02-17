package process

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"projects/chatroom1/common/message"
	"projects/chatroom1/common/model"
	"projects/chatroom1/common/utils"
	"sync"
)

var users map[int64]*UserProcessor
var mu sync.RWMutex

func init() {
	users = make(map[int64]*UserProcessor)
}

// RecordUser records user when login.
func RecordUser(userProc *UserProcessor) {
	mu.Lock()
	_, ok := users[userProc.UserID]
	if !ok {
		users[userProc.UserID] = userProc
	}
	mu.Unlock()
}

// DeleteUser deletes a user.
func DeleteUser(userID int64) {
	mu.Lock()
	delete(users, userID)
	mu.Unlock()
}

// ShowAllUsers shows all online users.
func ShowAllUsers() {
	for id, user := range users {
		fmt.Printf("user ID: %d, user info: %+v\n", id, user.User)
	}
}

// NotifyOthersMsg notify other users the message.
func NotifyOthersMsg(userID int64, msg *message.Message) (err error) {
	m, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Marshal NotifyOthersMsg data err: %v\n", err)
		return
	}

	mu.RLock()
	for _, user := range users {
		if user.UserID != userID {
			err = utils.WritePkg(user.Conn, m)
			if err != nil {
				log.Printf("err: %v, fail to send message to %s\n", err, user.UserAcc)
				continue
			}
		}
	}
	mu.RUnlock()
	return
}

// NotifyAllMsg notifies other users the message.
func NotifyAllMsg(msg *message.Message) (err error) {
	m, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Marshal NotifyAllMsg data err: %v\n", err)
		return
	}

	mu.RLock()
	for _, user := range users {
		err = utils.WritePkg(user.Conn, m)
		if err != nil {
			log.Printf("err: %v, fail to send message to %s\n", err, user.UserAcc)
			continue
		}
	}
	mu.RUnlock()
	return
}

// SendSmsMsg sends sms message to client.
func SendSmsMsg(msg *message.Message) (err error) {
	var smsMsg message.SmsMsg

	err = json.Unmarshal([]byte(msg.Data), &smsMsg)
	if err != nil {
		log.Printf("SendSmsMsg -> Unmarshal SmsMsg err:%v\n", err)
		return
	}
	userID := smsMsg.UserID

	m, err := json.Marshal(msg)
	if err != nil {
		log.Printf("SendSmsMsg -> Marshal Message err:%v\n", err)
		return
	}
	mu.RLock()
	for _, user := range users {
		if user.UserID == userID {
			err = utils.WritePkg(user.Conn, m)
			if err != nil {
				log.Printf("err: %v, fail to send message to %s\n", err, user.UserAcc)
			}
			break
		}
	}
	mu.RUnlock()

	return
}

// UserProcessor handles requests related to users, like login, register, logout,
// user list mangement and so on.
type UserProcessor struct {
	model.User
	Conn net.Conn
}

// ProcessRegister processes client register.
func (userPro *UserProcessor) ProcessRegister(msg *message.Message) (err error) {
	// Unmarshall the data
	var registerMsg message.RegisterMsg
	err = json.Unmarshal([]byte(msg.Data), &registerMsg)
	if err != nil {
		log.Printf("ProcessRegister -> Unmarshall RegisterMsg fail, err: %v\n", err)
		return
	}

	var replyMsg message.Message
	var registerResMsg message.RegisterResMsg

	replyMsg.Type = message.RegisterResMsgType

	// TODO: Register user to database
	newUser := &model.User{
		UserAcc:  registerMsg.UserAcc,
		UserName: registerMsg.UserName,
		UserPwd:  registerMsg.UserPwd,
	}

	_, err = model.CreateUser(newUser)
	if err == nil {
		registerResMsg.Code = 100
	} else {
		registerResMsg.Code = 200
		registerResMsg.Error = err.Error()
	}

	// serialize loginResMsg
	data, err := json.Marshal(registerResMsg)
	if err != nil {
		log.Printf("ProcessRegister -> Marshal RegisterResMsg fail, err: %v\n", err)
		return
	}

	replyMsg.Data = string(data)

	msgData, err := json.Marshal(replyMsg)
	if err != nil {
		log.Printf("ProcessRegister -> Marshal ReplyMsg fail, err: %v\n", err)
		return
	}

	// Send msgData
	err = utils.WritePkg(userPro.Conn, msgData)
	if err != nil {
		log.Printf("ProcessRegister -> WritePkg fail, err: %v\n", err)
		return
	}

	return
}

// ProcessLogin processes client login.
func (userPro *UserProcessor) ProcessLogin(msg *message.Message) (user *model.User, err error) {
	// Unmarshall the data
	var loginMsg message.LoginMsg
	err = json.Unmarshal([]byte(msg.Data), &loginMsg)
	if err != nil {
		fmt.Printf("ProcessLogin -> Unmarshall fail, err: %v\n", err)
		return
	}

	var replyMsg message.Message
	var loginResMsg message.LoginResMsg

	replyMsg.Type = message.LoginResMsgType

	newUser, _ := model.GetUserByAcc(loginMsg.UserAcc)
	if newUser == nil {
		loginResMsg.Code = 400
		loginResMsg.Error = model.ErrUserNotExists.Error()
	} else {
		// Validate the user password.
		if newUser.UserPwd == loginMsg.UserPwd {
			loginResMsg.Code = 300
			loginResMsg.User = *newUser
			user = newUser
		} else {
			loginResMsg.Code = 400
			loginResMsg.Error = model.ErrWrongPwd.Error()
		}
	}

	// serialize loginResMsg
	data, err := json.Marshal(loginResMsg)
	if err != nil {
		log.Printf("ProcessLogin -> Marshal LoginResMsg fail, err: %v\n", err)
		return
	}

	replyMsg.Data = string(data)

	msgData, err := json.Marshal(replyMsg)
	if err != nil {
		log.Printf("ProcessLogin -> Marshal ReplyMsg fail, err: %v\n", err)
		return
	}

	// Send msgData
	err = utils.WritePkg(userPro.Conn, msgData)
	if err != nil {
		log.Printf("ProcessLogin -> WritePkg fail, err: %v\n", err)
	}

	return
}

// ProcessLogout processes the logout business logic.
func (userPro *UserProcessor) ProcessLogout() (err error) {

	DeleteUser(userPro.UserID)

	// Notice other users the user logout status.
	var msg message.Message
	var noticeMsg message.NoticeMsg
	msg.Type = message.NoticeMsgType

	noticeMsg.Content = userPro.UserName + "下线了"

	data, err := json.Marshal(noticeMsg)
	if err != nil {
		fmt.Printf("ProcessLogout -> marshal notice message err: %v\n", err)
		return
	}
	msg.Data = string(data)
	err = NotifyOthersMsg(userPro.UserID, &msg)
	if err != nil {
		fmt.Printf("ProcessLogout -> notice others message err: %v\n", err)
		return
	}

	return
}

// ProcessReqOnlineUserList processes online user list request from client.
func (userPro *UserProcessor) ProcessReqOnlineUserList() error {
	var msg message.Message
	var reqRes message.OnlineUserListReqResMsg
	msg.Type = message.OnlineUserListReqResMsgType

	u := model.User{}
	mu.RLock()
	for _, user := range users {
		u.UserID = user.UserID
		u.UserName = user.UserName
		reqRes.OnlineUsers = append(reqRes.OnlineUsers, u)
	}
	mu.RUnlock()

	data, err := json.Marshal(reqRes)
	if err != nil {
		log.Printf("ProcessReqOnlineUserList -> Marshal OnlineUserListReqResMsg err: %v\n", err)
		return err
	}

	msg.Data = string(data)
	m, err := json.Marshal(msg)
	if err != nil {
		log.Printf("ProcessReqOnlineUserList -> Marshal Message err: %v\n", err)
		return err
	}

	err = utils.WritePkg(userPro.Conn, m)
	if err != nil {
		log.Printf("ProcessReqOnlineUserList -> WritePkg err: %v\n", err)
		return err
	}

	return nil
}
