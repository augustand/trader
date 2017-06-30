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
	signBalanceOf = "balanceOf(address)"
)

var signatures []string = []string{signBalanceOf}

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

	code := fmt.Sprintf("0x%v%v", ERC20Signatures[signBalanceOf], addrb)
	log.Println("code:", code, "#")
	if resp, err := http.Post(globalConfig.geth,
		"application/json",
		bytes.NewBufferString(fmt.Sprintf(`{"jsonrpc":"2.0","method": "eth_call", "params": [{"from": "%v", "to": "%v", "data": "%v"}, "latest"], "id": 0}`, globalConfig.account, contract, code))); err == nil {
		jsonParsed, _ = gabs.ParseJSONBuffer(resp.Body)
		count, ok := jsonParsed.Path("result").Data().(string)
		log.Println(jsonParsed)
		if !ok {
			w.Write([]byte(jsonParsed.Path("error").String()))
			return
		}
		ret := fmt.Sprintf(`{"count":"%v"}`, count)
		w.Write([]byte(ret))
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
