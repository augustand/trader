package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Jeffail/gabs"
	"github.com/julienschmidt/httprouter"
)

func getGasPriceHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ret := fmt.Sprintf(`{"gasPrice":"%v"}`, latestGasPrice.Load())
	w.Write([]byte(ret))
}

func getEstimateGas(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	jsonParsed, _ := gabs.ParseJSONBuffer(r.Body)
	// from, to, data string, gas, gasPrice, value float64)
	from, _ := jsonParsed.Path("from").Data().(string)
	to, _ := jsonParsed.Path("to").Data().(string)
	log.Println(jsonParsed.String())
	if len(from) == 0 || len(to) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"message":"%v", "code":%v}`, "illegal params", http.StatusBadRequest)))
		return
	}

	// optional
	data, _ := jsonParsed.Path("data").Data().(string)
	gas, _ := jsonParsed.Path("gas").Data().(string)
	gasPrice, _ := jsonParsed.Path("gasPrice").Data().(string)
	value, _ := jsonParsed.Path("value").Data().(string)
	ret, err := ethEstimateGas(from, to, data, gas, gasPrice, value)
	log.Println(ret, err)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"message":"%v", "code":%v}`, err, http.StatusBadRequest)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"value":"%v"}`, ret)))
}

func getTransactionCountHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	jsonParsed, _ := gabs.ParseJSONBuffer(r.Body)
	value, ok := jsonParsed.Path("address").Data().(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if resp, err := http.Post(globalConfig.geth,
		"application/json",
		bytes.NewBufferString(fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getTransactionCount","params":["%v","pending"],"id":1}`, value))); err == nil {
		jsonParsed, _ = gabs.ParseJSONBuffer(resp.Body)
		value, ok = jsonParsed.Path("result").Data().(string)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(jsonParsed.Path("error").String()))
			return
		}
		ret := fmt.Sprintf(`{"count":"%v"}`, value)
		w.Write([]byte(ret))
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func sendRawTransactionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	jsonParsed, _ := gabs.ParseJSONBuffer(r.Body)
	value, ok := jsonParsed.Path("data").Data().(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if resp, err := http.Post(globalConfig.geth,
		"application/json",
		bytes.NewBufferString(fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_sendRawTransaction","params":["%v"],"id":1}`, value))); err == nil {
		jsonParsed, _ = gabs.ParseJSONBuffer(resp.Body)
		value, ok = jsonParsed.Path("result").Data().(string)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(jsonParsed.Path("error").String()))
			return
		}
		ret := fmt.Sprintf(`{"txHash":"%v"}`, value)
		w.Write([]byte(ret))
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func getBalanceHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	jsonParsed, _ := gabs.ParseJSONBuffer(r.Body)
	value, ok := jsonParsed.Path("address").Data().(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if resp, err := http.Post(globalConfig.geth,
		"application/json",
		bytes.NewBufferString(fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getBalance","params":["%v","latest"],"id":1}`, value))); err == nil {
		jsonParsed, _ = gabs.ParseJSONBuffer(resp.Body)
		value, ok = jsonParsed.Path("result").Data().(string)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(jsonParsed.Path("error").String()))
			return
		}
		ret := fmt.Sprintf(`{"value":"%v"}`, value)
		w.Write([]byte(ret))
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

type getTransactionStruct struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Value    string `json:"value"`
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`
	Nonce    string `json:"nonce"`
}

func getTransactionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	jsonParsed, _ := gabs.ParseJSONBuffer(r.Body)
	txhash, ok := jsonParsed.Path("txHash").Data().(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if resp, err := http.Post(globalConfig.geth,
		"application/json",
		bytes.NewBufferString(fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getTransactionByHash","params":["%v"], "id":1}`, txhash))); err == nil {
		c, _ := gabs.ParseJSONBuffer(resp.Body)
		if !c.ExistsP("result") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(jsonParsed.Path("error").String()))
			return
		}

		tx := getTransactionStruct{}
		tx.From, _ = c.Path("result.from").Data().(string)
		tx.To, _ = c.Path("result.to").Data().(string)
		tx.Value, _ = c.Path("result.value").Data().(string)
		tx.Gas, _ = c.Path("result.gas").Data().(string)
		tx.GasPrice, _ = c.Path("result.gasPrice").Data().(string)
		tx.Nonce, _ = c.Path("result.nonce").Data().(string)

		enc := json.NewEncoder(w)
		enc.Encode(tx)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func blockNumberHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if resp, err := http.Post(globalConfig.geth,
		"application/json",
		bytes.NewBufferString(`{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`)); err == nil {
		jsonParsed, _ := gabs.ParseJSONBuffer(resp.Body)
		value, ok := jsonParsed.Path("result").Data().(string)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(jsonParsed.Path("error").String()))
			return
		}
		ret := fmt.Sprintf(`{"value":"%v"}`, value)
		w.Write([]byte(ret))
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func blockPerSecondHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bps := defaultBlockTimeEstimator.getBPS()
	ret := fmt.Sprintf(`{"bps":%v}`, bps)
	w.Write([]byte(ret))
}
