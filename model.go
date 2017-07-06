package main

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	ETHTransactions []ETHTranaction
	Contacts        []Contact
}

type ETHTranaction struct {
	gorm.Model
	UserId      int
	From        string
	To          string
	Value       string
	Nonce       string
	Gas         string
	GasPrice    string
	BlockNumber string
	Hash        string
}

type Contact struct {
	gorm.Model
	UserId  int
	Address string
}
