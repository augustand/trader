package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Jeffail/gabs"
	"github.com/julienschmidt/httprouter"
)

var tmpl = `{error:%v}`

func getBtcTransactionById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	jsonParsed, err := gabs.ParseJSONBuffer(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(tmpl, err)))
		return
	}

	var id string
	var ok bool
	if id, ok = jsonParsed.Path("txid").Data().(string); !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(tmpl, "get txid err")))
		return
	}

	get(w, fmt.Sprintf("%v/insight-api/tx/%v", globalConfig.insight, id))
}

func getBtcTransactions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var address string
	var from, to float64
	var ok bool
	jsonParsed, err := gabs.ParseJSONBuffer(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(tmpl, err)))
		return
	}

	if address, ok = jsonParsed.Path("address").Data().(string); !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(tmpl, "get address err")))
		return
	}

	if exists := jsonParsed.Exists("from"); exists {
		if from, ok = jsonParsed.Path("from").Data().(float64); !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf(tmpl, "parse from err")))
			return
		}
	}

	if exists := jsonParsed.Exists("to"); exists {
		if to, ok = jsonParsed.Path("to").Data().(float64); !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf(tmpl, "parse to err")))
			return
		}
	}

	if from > to {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(tmpl, fmt.Sprintf("from:%v, to:%v", from, to))))
		return
	}

	get(w, fmt.Sprintf("%v/insight-api/addrs/%v/txs?from=%v&to=%v", globalConfig.insight, address, from, to))
}

// /insight-api/addrs/2NF2baYuJAkCKo5onjUKEPdARQkZ6SYyKd5,2NAre8sX2povnjy4aeiHKeEh97Qhn97tB1f/utxo
func send(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var rawtx string
	var ok bool
	jsonParsed, err := gabs.ParseJSONBuffer(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(tmpl, err)))
		return
	}

	if rawtx, ok = jsonParsed.Path("rawtx").Data().(string); !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(tmpl, "get rawtx err")))
		return
	}

	if resp, err := http.PostForm(fmt.Sprintf("%v/insight-api/tx/send", globalConfig.insight), url.Values{"rawtx": {rawtx}}); err == nil {
		defer resp.Body.Close()
		if bts, err := ioutil.ReadAll(resp.Body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf(tmpl, err)))
			return
		} else {
			w.Write([]byte(bts))
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(tmpl, err)))
		return
	}
}

// /insight-api/tx/send
func getUtxo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var address string
	var ok bool
	jsonParsed, err := gabs.ParseJSONBuffer(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(tmpl, err)))
		return
	}

	if address, ok = jsonParsed.Path("address").Data().(string); !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(tmpl, "get address err")))
		return
	}

	get(w, fmt.Sprintf("%v/insight-api/addrs/%v/utxo", globalConfig.insight, address))
}

func estimatefee(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	get(w, fmt.Sprintf("%v/insight-api/utils/estimatefee", globalConfig.insight))
}

func balance(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var address string
	var ok bool
	jsonParsed, err := gabs.ParseJSONBuffer(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(tmpl, err)))
		return
	}
	if address, ok = jsonParsed.Path("address").Data().(string); !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(tmpl, "get address err")))
		return
	}
	get(w, fmt.Sprintf("%v/insight-api/addr/%v/balance", globalConfig.insight, address))
}

func get(w http.ResponseWriter, url string) {
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		if bts, err := ioutil.ReadAll(resp.Body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf(tmpl, err)))
			return
		} else {
			w.Write([]byte(bts))
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(tmpl, err)))
		return
	}
}
