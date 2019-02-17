package model

// User model
type User struct {
	UserID   int64  `json:"userID"`
	UserAcc  string `json:"userAcc"`
	UserPwd  string `json:"userPwd"`
	UserName string `json:"userName"`
}
