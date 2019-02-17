package model

import (
	"fmt"
	"testing"
)

func TestCreateUser(t *testing.T) {
	user := &User{
		UserAcc:  "test@qq.com",
		UserName: "test",
		UserPwd:  "123",
	}
	userID, err := CreateUser(user)
	if err != nil {
		t.Fatalf("create user fail, err: %v\n", err)
	}

	fmt.Printf("create user success, userID: %d\n", userID)
}

func TestGetUserByAcc(t *testing.T) {
	newUser := &User{
		UserAcc:  "dean@qq.com",
		UserName: "dean",
		UserPwd:  "123",
	}

	user, err := GetUserByAcc("dean@qq.com")
	if err != nil {
		t.Fatalf("GetUserByAcc fail, err: %v\n", err)
	}

	if newUser.UserName != user.UserName || newUser.UserPwd != user.UserPwd {
		t.Errorf("user != newUser, want %v, got %v\n", newUser, user)
	}

	user1, err := GetUserByAcc("noone@qq.com")
	if user1 != nil {
		t.Errorf("GetUserByAcc fail, want a nil user, got a user")
	}
}

func TestGetUserByID(t *testing.T) {
	// Test already exists data.
	user := &User{
		UserID:   2,
		UserAcc:  "test@qq.com",
		UserName: "test",
		UserPwd:  "123",
	}

	newUser, err := GetUserByID(2)
	if err != nil {
		t.Errorf("GetUserByID fail, err: %v\n", err)
	}
	if *newUser != *user {
		t.Errorf("GetUserByID fail, want %+v, got %+v\n", user, newUser)
	}

	// Test not exists data.
	nilUser, _ := GetUserByID(200)
	if nilUser != nil {
		t.Errorf("GetUserByID fail, want nil user, got %+v\n", nilUser)
	}
}
