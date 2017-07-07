package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type tokenBalanceOfStruct struct {
	Contract string `json:"contract"`
	Address  string `json:"address"`
}

func tokenBalanceOfHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	abi := tokenBalanceOfStruct{}
	dec := json.NewDecoder(r.Body)
	dec.Decode(&abi)
	var err error
	abi.Address, err = paduint(abi.Address, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data := fmt.Sprintf("0x%v%v", ERC20Signatures[signBalanceOf], abi.Address)
	if ret, err := ethCall(abi.Contract, data); err == nil {
		fmt.Fprintf(w, `{"value":"%v"}`, ret)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

type tokenTotalSupplyStruct struct {
	Contract string `json:"contract"`
}

func tokenTotalSupplyHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	abi := tokenTotalSupplyStruct{}
	dec := json.NewDecoder(r.Body)
	dec.Decode(&abi)
	data := fmt.Sprintf("0x%v", ERC20Signatures[signTotalSupply])
	if ret, err := ethCall(abi.Contract, data); err == nil {
		fmt.Fprintf(w, `{"value":"%v"}`, ret)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

type transferABIStruct struct {
	Contract string `json:"contract"`
	To       string `json:"to"`
	Value    string `json:"value"`
}

func transferABIHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	abi := transferABIStruct{}
	dec := json.NewDecoder(r.Body)
	dec.Decode(&abi)
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
	w.Write([]byte(fmt.Sprintf(`{"contract":"%v", "data":"%v"}`, abi.Contract, data)))
}
