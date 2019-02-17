package process

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"projects/chatroom1/common/message"
	"projects/chatroom1/common/utils"
)

// Register processes client side register business logic.
func Register(userAcc, userName, userPwd, confirmPwd string) (err error) {
	// Collected user's information
	fmt.Print("请输入用户账号：")
	fmt.Scanf("%s\n", &userAcc)

	fmt.Print("请输入用户名：")
	fmt.Scanf("%s\n", &userName)

	for {
		fmt.Print("请输入用户密码：")
		fmt.Scanf("%s\n", &userPwd)

		fmt.Print("请输入确认密码：")
		fmt.Scanf("%s\n", &confirmPwd)

		// Compare user password and confirm password
		if userPwd != confirmPwd {
			log.Println("确认密码不正确！")
			continue
		}
		break
	}

	// Connect the server
	conn, err := net.Dial("tcp", "localhost:8889")
	// conn, err := net.Dial("tcp", "192.168.1.2:8889")
	// conn, err := net.Dial("tcp", "192.168.0.108:8889")
	if err != nil {
		log.Printf("func Register -> Dial err: %v", err)
		return
	}
	defer conn.Close()

	// Initialize register message
	var registerMsg = &message.RegisterMsg{}
	registerMsg.UserAcc = userAcc
	registerMsg.UserName = userName
	registerMsg.UserPwd = userPwd

	var msg message.Message
	msg.Type = message.RegisterMsgType

	data, err := json.Marshal(registerMsg)
	if err != nil {
		log.Printf("func Register -> Marshal RegisterMsgType err: %v\n", err)
		return
	}
	msg.Data = string(data)

	m, err := json.Marshal(msg)
	if err != nil {
		log.Printf("func Register -> Marshal Message err: %v\n", err)
		return
	}

	// Send register message to server.
	err = utils.WritePkg(conn, m)
	if err != nil {
		log.Printf("func Register -> WritePkg err: %v\n", err)
		return
	}

	// Handle reply message from the server
	replyMsg, err := utils.ReadPkg(conn)
	if err != nil {
		log.Printf("func Register -> ReadPkg err: %v\n", err)
		return
	}

	var registerResMsg message.RegisterResMsg
	err = json.Unmarshal([]byte(replyMsg.Data), &registerResMsg)
	if err != nil {
		log.Printf("func Register -> Unmarshall RegisterResMsg fail, err: %v\n", err)
		return
	}

	if registerResMsg.Code == 100 {
		// Register successfully
		fmt.Println("恭喜您注册成功！")
	} else {
		err = fmt.Errorf(registerResMsg.Error)
		log.Println(err)
		return
	}

	return
}

// Login processes login business logic.
func Login(userAcc string, userPwd string) (err error) {

	// connect the server
	conn, err := net.Dial("tcp", "localhost:8889")
	// conn, err := net.Dial("tcp", "192.168.1.2:8889")
	// conn, err := net.Dial("tcp", "192.168.0.108:8889")
	if err != nil {
		log.Printf("Login -> dial err: %v", err)
		return
	}

	// TODO: Need optimization! Where to close the conn?
	defer conn.Close()

	// send messsages to server
	var msg message.Message
	msg.Type = message.LoginMsgType
	// create a LoginMsg message.
	var loginMsg message.LoginMsg
	loginMsg.UserAcc = userAcc
	loginMsg.UserPwd = userPwd

	// serialize loginMsg
	data, err := json.Marshal(loginMsg)
	if err != nil {
		fmt.Printf("Login -> Marshall LoginMsg err: %v\n", err)
		return
	}

	msg.Data = string(data)

	// serialize msg
	m, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("Login -> Marshall Message err: %v\n", err)
		return
	}

	err = utils.WritePkg(conn, m)
	if err != nil {
		log.Printf("Login -> WritePkg err: %v\n", err)
		return
	}

	// Handle reply message from the server
	replyMsg, err := utils.ReadPkg(conn)
	if err != nil {
		log.Printf("Login -> ReadPkg err: %v\n", err)
		return
	}

	var loginResMsg message.LoginResMsg
	err = json.Unmarshal([]byte(replyMsg.Data), &loginResMsg)
	if err != nil {
		log.Printf("Login -> Unmarshall LoginResMsg fail, err: %v\n", err)
		return
	}

	if loginResMsg.Code == 300 { // Login successfully.
		fmt.Printf("----------恭喜%s登录成功----------\n", loginResMsg.UserName)

		// After successfully login, run a goroutine keeping commnucating with server
		go HandleServerMsg(conn)

		// Show the user menu
		ShowUserMenu(loginResMsg.UserID, loginResMsg.UserName, conn)
	} else { // Login fail.
		err = errors.New(loginResMsg.Error)
	}

	return
}

// HandleServerMsg keeping reading messages from server and processes them.
func HandleServerMsg(conn net.Conn) {
	fmt.Println("\nReading message from server...")
	for {
		msg, err := utils.ReadPkg(conn)
		if err != nil {
			log.Printf("ProcessServerMsg -> ReadMsg err: %v\n", err)
			return
		}

		switch msg.Type {
		case message.NoticeMsgType:

			var noticeMsg message.NoticeMsg
			err := json.Unmarshal([]byte(msg.Data), &noticeMsg)
			if err != nil {
				log.Printf("ProcessServerMsg -> Unmarshal NoticeMsg err: %v\n", err)
				return
			}
			fmt.Println(noticeMsg.Content)
		case message.OnlineUserListReqResMsgType:
			ShowOnlineUsers(msg)
		case message.GroupSmsMsgType:
			ShowGrpMsg(msg)
		case message.SmsMsgType:
			ShowSmsMsg(msg)
		default:
			err = fmt.Errorf("unkonwn message type")
			log.Println(err)
		}
	}
}

// ShowOnlineUsers shows all online users.
func ShowOnlineUsers(msg *message.Message) {
	var resMsg message.OnlineUserListReqResMsg

	err := json.Unmarshal([]byte(msg.Data), &resMsg)
	if err != nil {
		log.Printf("ShowOnlineUsers -> Unmarshal OnlineUserListReqResMsg err: %v\n", err)
	}

	fmt.Println("----------在线用户列表----------")
	for _, user := range resMsg.OnlineUsers {
		fmt.Printf("ID: %d, name: %s\n", user.UserID, user.UserName)
	}
}

// ShowUserMenu shows the menu after user has logined in.
func ShowUserMenu(userID int64, userName string, conn net.Conn) {
	// reader用于读取字符串
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("----------1. 获取在线用户列表----------")
		fmt.Println("----------2. 群聊----------")
		fmt.Println("----------3. 私聊----------")
		fmt.Println("----------4. 退出系统----------")
		fmt.Println("请选择(1-4):")
		var key int
		fmt.Scanf("%d\n", &key)

		var userID int64

		switch key {
		case 1:
			err := ReqOnlineUsers(conn)
			if err != nil {
				fmt.Println(err)
			}
		case 2:
			fmt.Println("----------群聊----------")
			fmt.Println("请输入群聊信息：")
			content, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("input group chat content err: %v", err)
				continue
			}
			content = userName + ": " + content

			err = SendGrpSmsMsg(conn, content)
			if err != nil {
				log.Printf("SendGrpSmsMsg err: %s\n", err)
			}
		case 3:
			fmt.Println("----------私聊----------")
			fmt.Println("请输入对方ID：")
			fmt.Scanf("%d\n", &userID)

			fmt.Println("请输入私聊信息：")
			content, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("input private chat content err: %v", err)
				continue
			}
			content = userName + ": " + content

			err = SendSmsMsg(conn, userID, content)
			if err != nil {
				log.Printf("SendSmsMsg err: %s\n", err)
			}
		case 4:
			fmt.Println("----------退出系统----------")
			err := Logout(userID, conn)
			if err != nil {
				log.Println(err)
			}
			os.Exit(0)
		default:
			fmt.Println("你输入的选项不正确，请重新输入！")
		}
	}
}

// ReqOnlineUsers requests online user list from server.
func ReqOnlineUsers(conn net.Conn) error {
	var msg message.Message
	msg.Type = message.OnlineUserListReqMsgType

	m, err := json.Marshal(msg)
	if err != nil {
		log.Printf("ReqOnlineUsers -> marshal Message err: %v\n", err)
		return err
	}

	err = utils.WritePkg(conn, m)
	if err != nil {
		log.Printf("ReqOnlineUsers -> WritePkg err: %v\n", err)
		return err
	}

	return nil
}

// Logout logout the program.
func Logout(userID int64, conn net.Conn) (err error) {
	defer conn.Close()

	var msg message.Message
	msg.Type = message.LogoutMsgType

	m, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Logout -> Marshal Message err: %v\n", err)
		return
	}

	err = utils.WritePkg(conn, m)
	if err != nil {
		err = utils.WritePkg(conn, m)
		log.Printf("Logout -> WritePkg err: %v\n", err)
	}
	return
}
