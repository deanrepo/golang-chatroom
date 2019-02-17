package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	UserID   int64  `json:"userID"`
	UserAcc  string `json:"userAcc"`
	UserPwd  string `json:"userPwd"`
	UserName string `json:"userName"`
}

var Db *sql.DB

func init() {
	var err error
	driverName := "mysql"
	dsn := "dean:Dean#168168@tcp(192.168.1.12:3306)/chatroom?collation=utf8_general_ci"
	Db, err = sql.Open(driverName, dsn)
	if err != nil {
		log.Printf("open database err: %v\n", err)
		panic(err)
	}

	Db.SetMaxOpenConns(200)
	Db.SetMaxIdleConns(100)
	Db.Ping()
}

func main() {

	// db, err := sql.Open("mysql", "dean:Dean#168168@tcp(192.168.1.12:3306)/chatroom?collation=utf8_general_ci")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	user := User{}
	err := Db.QueryRow("SELECT user_id, user_acc, user_name, user_pwd FROM user WHERE user_acc = ?",
		"dean@qq.com").Scan(&user.UserID, &user.UserAcc, &user.UserName, &user.UserPwd)

	if err != nil {
		log.Printf("GetUserByID -> QueryRow err: %v\n", err)
		return
	}

	fmt.Println(user)

}
