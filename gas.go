package main

import (
	"bytes"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/Jeffail/gabs"
	log "github.com/Sirupsen/logrus"
)

var latestGasPrice atomic.Value

func update_gas_task() {
	for {
		resp, err := http.Post(globalConfig.geth,
			"application/json",
			bytes.NewBufferString(`{"jsonrpc":"2.0","method":"eth_gasPrice","params":[],"id":0}`))

		if err != nil {
			log.Println("call:", err)
		}
		jsonParsed, _ := gabs.ParseJSONBuffer(resp.Body)
		value, ok := jsonParsed.Path("result").Data().(string)
		if !ok {
			log.Println("cannot get gasPrice", jsonParsed)
		}
		latestGasPrice.Store(value)
		<-time.After(globalConfig.gasUpdate)
	}
}
