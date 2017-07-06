package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Jeffail/gabs"
	"github.com/jinzhu/gorm"
)

type TransactionManager struct {
	db *gorm.DB
}

func NewTransactionManger(db *gorm.DB) *TransactionManager {
	db.AutoMigrate(&ETHTranaction{})
	return &TransactionManager{db}
}

func (db *TransactionManager) record(txhash string) {
	if resp, err := http.Post(globalConfig.geth,
		"application/json",
		bytes.NewBufferString(fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getTransactionByHash","params":["%v"], "id":1}`, txhash))); err == nil {
		bts, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return
		}

		if c, err := gabs.ParseJSON(bts); err == nil {
			tx := ETHTranaction{}
			tx.From, _ = c.Path("result.from").Data().(string)
			tx.To, _ = c.Path("result.to").Data().(string)
			tx.Value, _ = c.Path("result.value").Data().(string)
			tx.Gas, _ = c.Path("result.gas").Data().(string)
			tx.GasPrice, _ = c.Path("result.gasPrice").Data().(string)
			tx.Nonce, _ = c.Path("result.nonce").Data().(string)
			tx.Hash, _ = c.Path("result.hash").Data().(string)
			db.db.Create(&tx)
		}
	} else {
		log.Println(err)
	}
}

var defaultETHTXManager *TransactionManager
