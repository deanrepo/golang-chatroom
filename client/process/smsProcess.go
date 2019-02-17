// Handle sms messages.
package process

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"projects/chatroom1/common/message"
	"projects/chatroom1/common/utils"
)

// SendGrpSmsMsg sends group sms message.
func SendGrpSmsMsg(conn net.Conn, content string) (err error) {
	var msg message.Message
	var grpMsg message.GroupSmsMsg
	msg.Type = message.GroupSmsMsgType
	grpMsg.Content = content

	d, err := json.Marshal(grpMsg)
	if err != nil {
		log.Printf("SendGrpSmsMsg -> Marshal GroupSmsMsg err: %v\n", err)
		return
	}

	msg.Data = string(d)
	m, err := json.Marshal(msg)
	if err != nil {
		log.Printf("SendGrpSmsMsg -> Marshal Message err: %v\n", err)
		return
	}

	err = utils.WritePkg(conn, m)
	if err != nil {
		log.Printf("SendGrpSmsMsg -> WritePkg err: %v\n", err)
		return
	}

	return
}

// ShowGrpMsg shows group message.
func ShowGrpMsg(msg *message.Message) {
	var grpMsg message.GroupSmsMsg
	err := json.Unmarshal([]byte(msg.Data), &grpMsg)
	if err != nil {
		log.Printf("ShowGrpMsg -> Unmarshal GroupSmsMsg err: %v\n", err)
		return
	}

	fmt.Println(grpMsg.Content)
}

// SendSmsMsg sends sms message.
func SendSmsMsg(conn net.Conn, userID int64, content string) (err error) {
	var msg message.Message
	var smsMsg message.SmsMsg
	msg.Type = message.SmsMsgType
	smsMsg.UserID = userID
	smsMsg.Content = content

	d, err := json.Marshal(smsMsg)
	if err != nil {
		log.Printf("SendSmsMsg -> Marshal SmsMsg err: %v\n", err)
		return
	}

	msg.Data = string(d)
	m, err := json.Marshal(msg)
	if err != nil {
		log.Printf("SendSmsMsg -> Marshal Message err: %v\n", err)
		return
	}

	err = utils.WritePkg(conn, m)
	if err != nil {
		log.Printf("SendSmsMsg -> WritePkg err: %v\n", err)
		return
	}

	return
}

// ShowSmsMsg shows sms message.
func ShowSmsMsg(msg *message.Message) {
	var smsMsg message.SmsMsg
	err := json.Unmarshal([]byte(msg.Data), &smsMsg)
	if err != nil {
		log.Printf("ShowSmsMsg -> Unmarshal SmsMsg err: %v\n", err)
		return
	}

	fmt.Println(smsMsg.Content)
}
