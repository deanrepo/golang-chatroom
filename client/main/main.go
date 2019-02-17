package main

import (
	"fmt"
	"log"
	"projects/chatroom1/client/process"
)

// User ID
var userID int64

// User account
var userAcc string

// User password
var userPwd string

// User confirm password
var confirmPwd string

// User name
var userName string

func main() {

	// User's choice
	var key int

	// 显示主界面
loop:
	for {
		fmt.Println("-------------------欢迎登录多人聊天系统-------------------")
		fmt.Println("                   1. 登录聊天室")
		fmt.Println("                   2. 注册用户")
		fmt.Println("                   3. 退出系统")
		fmt.Println("=>请选择(1-3):")

		fmt.Scanf("%d\n", &key)

		switch key {
		case 1:
			fmt.Println("-------------------登录聊天室-------------------")
			// TODO: LOGIN BUSINESS LOGIC.
			fmt.Print("请输入用户账号：")
			fmt.Scanf("%s\n", &userAcc)

			fmt.Print("请输入用户密码：")
			fmt.Scanf("%s\n", &userPwd)

			// TODO: TO PROCESS THE ERROR
			err := process.Login(userAcc, userPwd)
			if err != nil {
				log.Printf("Login err: %v.\n", err)
				continue
			}

		case 2:
			fmt.Println("-------------------注册用户-------------------")
			err := process.Register(userAcc, userName, userPwd, confirmPwd)
			if err != nil {
				log.Printf("Register err: %v, please register again.\n", err)
				continue
			}
		case 3:
			fmt.Println("-------------------退出系统-------------------")
			break loop
		default:
			fmt.Println("您输入有误， 请重新输入！")
		}
	}
}
