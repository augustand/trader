package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/julienschmidt/httprouter"
)

type tokenBalanceOfStruct struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func tokenBalanceOfHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	abi := tokenBalanceOfStruct{}
	dec := json.NewDecoder(r.Body)
	dec.Decode(&abi)
	contract, ok := contractAddresses[abi.Name]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var err error
	abi.Address, err = paduint(abi.Address)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println(abi)

	data := fmt.Sprintf("0x%v%v", ERC20Signatures[signBalanceOf], abi.Address)
	if ret, err := eth_call(contract, data); err == nil {
		w.Write([]byte(ret))
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

type tokenTotalSupplyStruct struct {
	Name string `json:"name"`
}

func tokenTotalSupplyHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	abi := tokenTotalSupplyStruct{}
	dec := json.NewDecoder(r.Body)
	dec.Decode(&abi)
	contract, ok := contractAddresses[abi.Name]
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

func paduint(value string) (string, error) {
	if !strings.HasPrefix(value, "0x") {
		return value, errors.New("must start with 0x")
	}

	value = value[2:]
	n := len(value)
	if n%32 == 0 {
		return value, nil
	} else {
		return strings.Repeat("0", 32-n%32) + value, nil
	}
}

type transferABIStruct struct {
	Name  string `json:"name"`
	To    string `json:"to"`
	Value string `json:"value"`
}

func transferABIHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	abi := transferABIStruct{}
	dec := json.NewDecoder(r.Body)
	dec.Decode(&abi)
	contract, ok := contractAddresses[abi.Name]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var err error
	abi.To, err = paduint(abi.To)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	abi.Value, err = paduint(abi.Value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data := fmt.Sprintf("0x%v%v%v", ERC20Signatures[signTransfer], abi.To, abi.Value)
	w.Write([]byte(fmt.Sprintf(`{"contract":"%v", "data":"%v"}`, contract, data)))
}
