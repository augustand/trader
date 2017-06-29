package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Jeffail/gabs"
	"github.com/julienschmidt/httprouter"
)

func getGasPriceHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ret := fmt.Sprintf(`{"jsonrpc":"2.0","id":1,"result":"%v"}`, latestGasPrice.Load())
	w.Write([]byte(ret))
}

func getTransactionCountHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	jsonParsed, _ := gabs.ParseJSONBuffer(r.Body)
	value, ok := jsonParsed.Path("address").Data().(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := http.Post(globalConfig.geth,
		"application/json",
		bytes.NewBufferString(fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getTransactionCount","params":["%v","latest"],"id":1}`, value)))
	if err != nil {
		log.Println("call:", err)
	}

	io.Copy(w, resp.Body)
}
