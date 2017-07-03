package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"

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
	abi.Address, err = paduint(abi.Address, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println(abi)

	data := fmt.Sprintf("0x%v%v", ERC20Signatures[signBalanceOf], abi.Address)
	if ret, err := eth_call(contract, data); err == nil {
		fmt.Fprintf(w, `{"value":"%v"}`, ret)
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
		fmt.Fprintf(w, `{"value":"%v"}`, ret)
	} else {
		w.WriteHeader(http.StatusBadRequest)
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
	abi.To, err = paduint(abi.To, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	abi.Value, err = paduint(abi.Value, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data := fmt.Sprintf("0x%v%v%v", ERC20Signatures[signTransfer], abi.To, abi.Value)
	w.Write([]byte(fmt.Sprintf(`{"contract":"%v", "data":"%v"}`, contract, data)))
}
