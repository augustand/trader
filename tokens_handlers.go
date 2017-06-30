package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/julienschmidt/httprouter"
	"github.com/xtaci/trader/sha3"
)

const (
	signBalanceOf   = "balanceOf(address)"
	signTotalSupply = "totalSupply()"
	signTransfer    = "transfer(address,uint)"
)

var signatures []string = []string{signBalanceOf, signTotalSupply, signTransfer}

var ERC20Signatures = make(map[string]string)

func init() {
	for _, sign := range signatures {
		d := sha3.NewKeccak256()
		d.Write([]byte(sign))
		h := d.Sum(nil)
		ERC20Signatures[sign] = hex.EncodeToString(h[0:4])
		log.Println(sign, ERC20Signatures[sign])
	}
}

func tokenBalanceOfHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	jsonParsed, _ := gabs.ParseJSONBuffer(r.Body)
	address, ok := jsonParsed.Path("address").Data().(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contract, ok := jsonParsed.Path("contract").Data().(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(address, "0x") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	addrb := strings.Repeat("0", 24) + address[2:]

	data := fmt.Sprintf("0x%v%v", ERC20Signatures[signBalanceOf], addrb)
	if ret, err := eth_call(contract, data); err == nil {
		w.Write([]byte(ret))
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func tokenTotalSupplyHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	jsonParsed, _ := gabs.ParseJSONBuffer(r.Body)
	contract, ok := jsonParsed.Path("contract").Data().(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data := fmt.Sprintf("0x%v", ERC20Signatures[signTotalSupply])
	if ret, err := eth_call(contract, data); err == nil {
		w.Write([]byte(ret))
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func eth_call(to, data string) (string, error) {
	if resp, err := http.Post(globalConfig.geth,
		"application/json",
		bytes.NewBufferString(fmt.Sprintf(`{"jsonrpc":"2.0","method": "eth_call", "params": [{"to": "%v", "data": "%v"}, "latest"], "id": 0}`, to, data))); err == nil {
		jsonParsed, _ := gabs.ParseJSONBuffer(resp.Body)
		value, ok := jsonParsed.Path("result").Data().(string)
		if !ok {
			return jsonParsed.Path("error").String(), nil
		}
		return fmt.Sprintf(`{"value":"%v"}`, value), nil
	} else {
		return "", err
	}
}
