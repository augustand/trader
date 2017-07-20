package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Jeffail/gabs"
	"github.com/julienschmidt/httprouter"
)

var tmpl = `{"code":%v, "message":"%v"}`

var (
	btc = 100000000.0
)

func getBtcTransactionById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	jsonParsed, err := gabs.ParseJSONBuffer(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(tmpl, http.StatusBadRequest, err)))
		return
	}

	var id string
	var ok bool
	if id, ok = jsonParsed.Path("txid").Data().(string); !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(tmpl, http.StatusBadRequest, "get txid err")))
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
		w.Write([]byte(fmt.Sprintf(tmpl, http.StatusBadRequest, err)))
		return
	}

	if address, ok = jsonParsed.Path("address").Data().(string); !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(tmpl, http.StatusBadRequest, "get address err")))
		return
	}

	if exists := jsonParsed.Exists("from"); exists {
		if from, ok = jsonParsed.Path("from").Data().(float64); !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf(tmpl, http.StatusBadRequest, "parse from err")))
			return
		}
	}

	if exists := jsonParsed.Exists("to"); exists {
		if to, ok = jsonParsed.Path("to").Data().(float64); !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf(tmpl, http.StatusBadRequest, "parse to err")))
			return
		}
	}

	if from > to {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(tmpl, http.StatusBadRequest, fmt.Sprintf("from:%v, to:%v", from, to))))
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
		w.Write([]byte(fmt.Sprintf(tmpl, http.StatusBadRequest, err)))
		return
	}

	if rawtx, ok = jsonParsed.Path("rawtx").Data().(string); !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(tmpl, http.StatusBadRequest, "get rawtx err")))
		return
	}

	resp, err := http.Post(fmt.Sprintf("%v/insight-api/tx/send", globalConfig.insight), "application/json", bytes.NewBufferString(fmt.Sprintf(`{"rawtx": "%v"}`, rawtx)))
	if err == nil {
		retval, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusOK {
			w.Write(retval)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf(tmpl, http.StatusBadRequest, string(retval))))
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(tmpl, http.StatusBadRequest, err)))
	}

	return
}

// /insight-api/tx/send
func getUtxo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var address string
	var ok bool
	jsonParsed, err := gabs.ParseJSONBuffer(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(tmpl, http.StatusBadRequest, err)))
		return
	}

	if address, ok = jsonParsed.Path("address").Data().(string); !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(tmpl, http.StatusBadRequest, "get address err")))
		return
	}

	get(w, fmt.Sprintf("%v/insight-api/addrs/%v/utxo", globalConfig.insight, address))
}

func estimatefee(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var nbBlocks float64
	var ok bool
	jsonParsed, err := gabs.ParseJSONBuffer(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(tmpl, http.StatusBadRequest, err)))
		return
	}

	if nbBlocks, ok = jsonParsed.Path("nbBlocks").Data().(float64); !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(tmpl, http.StatusBadRequest, "get nbBlocks err")))
		return
	}
	//					  /insight-api/utils/estimatefee[?nbBlocks=2]
	url := fmt.Sprintf("%v/insight-api/utils/estimatefee?nbBlocks=%v", globalConfig.insight, nbBlocks)
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		jsonParsed, err := gabs.ParseJSONBuffer(resp.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf(tmpl, http.StatusBadRequest, err)))
			return
		}

		val, ok := jsonParsed.Path(`2`).Data().(float64)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf(tmpl, http.StatusBadRequest, "paser ret err")))
			return
		}
		w.Write([]byte(fmt.Sprintf(`{"staoshi":"0x%x"}`, int64(val*btc))))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(tmpl, http.StatusBadRequest, err)))
		return
	}
}

func getAddress(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var address string
	var ok bool
	jsonParsed, err := gabs.ParseJSONBuffer(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(tmpl, http.StatusBadRequest, err)))
		return
	}
	if address, ok = jsonParsed.Path("address").Data().(string); !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(tmpl, http.StatusBadRequest, "get address err")))
		return
	}
	get(w, fmt.Sprintf("%v/insight-api/addr/%v?noTxList=1", globalConfig.insight, address))
}

func get(w http.ResponseWriter, url string) {
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		bts, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(resp.StatusCode)
			w.Write([]byte(fmt.Sprintf(tmpl, resp.StatusCode, err)))
			return
		}

		if resp.StatusCode != http.StatusOK {
			w.WriteHeader(resp.StatusCode)
			w.Write([]byte(fmt.Sprintf(tmpl, resp.StatusCode, string(bts))))
			return
		}
		w.Write([]byte(bts))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(tmpl, http.StatusBadRequest, err)))
		return
	}
}
