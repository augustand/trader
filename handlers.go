package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/Jeffail/gabs"
	"github.com/julienschmidt/httprouter"
)

func getGasPriceHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ret := fmt.Sprintf(`{"gasPrice":"%v"}`, latestGasPrice.Load())
	w.Write([]byte(ret))
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
		bytes.NewBufferString(fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getTransactionCount","params":["%v","latest"],"id":1}`, value))); err == nil {
		jsonParsed, _ = gabs.ParseJSONBuffer(resp.Body)
		value, ok = jsonParsed.Path("result").Data().(string)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ret := fmt.Sprintf(`{"count":"%v"}`, value)
		w.Write([]byte(ret))
	}
}
