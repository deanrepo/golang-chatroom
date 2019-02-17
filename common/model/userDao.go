package model

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Db used to manipulate mysql database.
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

// CreateUser creates a user to into the database.
func CreateUser(user *User) (userID int64, err error) {
	// Check whether user exists in the database.
	u, _ := GetUserByAcc(user.UserAcc)
	if u != nil {
		err = ErrUserExists
		return
	}

	stmt, err := Db.Prepare(`INSERT INTO user (user_acc,user_name,user_pwd) values (?,?,?)`)
	if err != nil {
		// log.Printf("CreateUser -> Prepare err: %v\n", err)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(user.UserAcc, user.UserName, user.UserPwd)
	if err != nil {
		// log.Printf("CreateUser -> Exec err: %v\n", err)
		return
	}
	userID, err = res.LastInsertId()
	if err != nil {
		// log.Printf("CreateUser -> LastInsertId err: %v\n", err)
		return
	}

	return
}

// GetUserByID gets user by user ID.
func GetUserByID(userID int64) (*User, error) {
	user := &User{}
	err := Db.QueryRow("SELECT user_id, user_acc, user_name, user_pwd FROM user WHERE user_id = ?",
		userID).Scan(&user.UserID, &user.UserAcc, &user.UserName, &user.UserPwd)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByAcc gets user by user account.
func GetUserByAcc(userAcc string) (*User, error) {
	user := User{}
	err := Db.QueryRow("SELECT user_id, user_acc, user_name, user_pwd FROM user WHERE user_acc = ?",
		userAcc).Scan(&user.UserID, &user.UserAcc, &user.UserName, &user.UserPwd)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
