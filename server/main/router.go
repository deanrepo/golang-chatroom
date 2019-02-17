package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"projects/chatroom1/common/message"
	"projects/chatroom1/common/model"
	"projects/chatroom1/common/utils"
	"projects/chatroom1/server/process"
)

func processConn(conn net.Conn) {
	defer conn.Close()

	// Define user processor.
	userPro := &process.UserProcessor{
		Conn: conn,
	}

	for {
		// Read message from client.
		msg, err := utils.ReadPkg(conn)
		// Handle OpError and EOF.
		if _, ok := err.(*net.OpError); ok || err == io.EOF {
			log.Printf("client close, so the server side close, too...\n")
			// Delete user and notify others
			if userPro.User != (model.User{}) {

				process.DeleteUser(userPro.UserID)

				// Notice other users the user logout status.
				var msg message.Message
				var noticeMsg message.NoticeMsg

				msg.Type = message.NoticeMsgType
				noticeMsg.Content = userPro.UserName + "下线了"

				data, err := json.Marshal(noticeMsg)
				if err != nil {
					fmt.Printf("processConn -> marshal notice message err: %v\n", err)
					return
				}
				msg.Data = string(data)
				err = process.NotifyOthersMsg(userPro.UserID, &msg)
				if err != nil {
					fmt.Printf("processConn -> notice others message err: %v\n", err)
					return
				}
			}
			return
		}

		if err != nil {
			log.Printf("err: %+v, continue to process...\n", err)
			continue
		}

		// Process received message.
		switch msg.Type {
		case message.RegisterMsgType:
			// Process register.
			err = userPro.ProcessRegister(msg)
			if err != nil {
				log.Printf("register fail, err: %v\n", err)
			}

			return
		case message.LoginMsgType:
			// Process login.
			user, err := userPro.ProcessLogin(msg)
			if err != nil {
				log.Printf("login fail, err: %v\n", err)
				return
			}

			// If login successfully, record the user info.
			if user != nil {
				userPro.User = *user
				process.RecordUser(userPro)

				// Notice other users the user login status.
				var msg message.Message
				msg.Type = message.NoticeMsgType
				var noticeMsg message.NoticeMsg
				noticeMsg.Content = user.UserName + "上线了"

				data, err := json.Marshal(noticeMsg)
				if err != nil {
					fmt.Printf("marshal notice message err: %v\n", err)
					continue
				}
				msg.Data = string(data)
				err = process.NotifyOthersMsg(userPro.UserID, &msg)
				if err != nil {
					log.Printf("notice others message err: %v\n", err)
					continue
				}
			}
		case message.OnlineUserListReqMsgType:
			err := userPro.ProcessReqOnlineUserList()
			if err != nil {
				log.Printf("process request online user list err: %v\n", err)
			}

		case message.LogoutMsgType:
			err := userPro.ProcessLogout()
			if err != nil {
				fmt.Printf("process logout err: %v\n", err)
			}
			return
		case message.GroupSmsMsgType:
			err := process.NotifyAllMsg(msg)
			if err != nil {
				fmt.Printf("send group sms message err: %v\n", err)
			}
		case message.SmsMsgType:
			err := process.SendSmsMsg(msg)
			if err != nil {
				fmt.Printf("send sms message err: %v\n", err)
			}
		default:
			err = fmt.Errorf("unkonwn message type")
			log.Println(err)
			return
		}
	}
}
