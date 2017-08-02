package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"encoding/hex"

	"github.com/Jeffail/gabs"
	log "github.com/Sirupsen/logrus"
	"github.com/xtaci/trader/sha3"
)

const (
	signBalanceOf   = "balanceOf(address)"
	signTotalSupply = "totalSupply()"
	signTransfer    = "transfer(address,uint256)"
	eventTransfer   = "Transfer(address,address,uint256)"
	initWallet      = "initWallet(address[],uint256,uint256)"
)

var (
	signatures      = []string{signBalanceOf, signTotalSupply, signTransfer, initWallet}
	events          = []string{eventTransfer}
	ERC20Signatures = make(map[string]string)
)

func init() {
	log.Info("Computing signatures")
	for _, sign := range signatures {
		d := sha3.NewKeccak256()
		d.Write([]byte(sign))
		h := d.Sum(nil)
		ERC20Signatures[sign] = hex.EncodeToString(h[0:4])
		log.Println(sign, ERC20Signatures[sign])
	}

	for _, sign := range events {
		d := sha3.NewKeccak256()
		d.Write([]byte(sign))
		h := d.Sum(nil)
		ERC20Signatures[sign] = hex.EncodeToString(h)
		log.Println(sign, ERC20Signatures[sign])
	}
}

func ethEstimateGas(from, to, data, gas, gasPrice, value string) (string, error) {
	if resp, err := http.Post(globalConfig.geth,
		"application/json",
		bytes.NewBufferString(fmt.Sprintf(`{"jsonrpc":"2.0","method": "eth_estimateGas","params": [{"to": "%v", "from": "%v" ,"gas": "%v" , "gasPrice":"%v", "value": "%v", "data": "%v"}, "latest"], "id": 0}`, to, from, gas, gasPrice, value, data))); err == nil {

		log.Println(fmt.Sprintf(`{"jsonrpc":"2.0","method": "eth_estimateGas","params": [{"to": "%v", "from": "%v" ,"gas": "%v" , "gasPrice":"%v", "value": "%v", data": "%v"}, "latest"], "id": 0}`, to, from, gas, gasPrice, value, data))
		jsonParsed, _ := gabs.ParseJSONBuffer(resp.Body)
		value, ok := jsonParsed.Path("result").Data().(string)
		if !ok {
			return jsonParsed.Path("error").String(), nil
		}
		return value, nil
	} else {
		return "", err
	}
}

func ethCall(to, data string) (string, error) {
	if resp, err := http.Post(globalConfig.geth,
		"application/json",
		bytes.NewBufferString(fmt.Sprintf(`{"jsonrpc":"2.0","method": "eth_call", "params": [{"to": "%v", "data": "%v"}, "latest"], "id": 0}`, to, data))); err == nil {
		jsonParsed, _ := gabs.ParseJSONBuffer(resp.Body)
		value, ok := jsonParsed.Path("result").Data().(string)
		if !ok {
			return jsonParsed.Path("error").String(), nil
		}
		return value, nil
	} else {
		return "", err
	}
}

func paduint(value string, size int) (string, error) {
	if !strings.HasPrefix(value, "0x") {
		return value, errors.New("must start with 0x")
	}

	value = value[2:]
	n := len(value)
	if n%size == 0 {
		return value, nil
	}
	return strings.Repeat("0", size-n%size) + value, nil
}
