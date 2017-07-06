package main

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type transactions struct {
	gorm.Model
	From      string
	To        string
	Value     string
	Timestamp time.Time
}

type transactionManager struct {
	db *gorm.DB
}

func (txm *transactionManager) init(conn string) {
	db, err := gorm.Open("postgres", conn)
	if err != nil {
		panic("failed to connect database")
	}
	txm.db = db
}

var defaultTransactionManager transactionManager
