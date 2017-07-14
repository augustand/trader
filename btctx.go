package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Jeffail/gabs"
	"github.com/julienschmidt/httprouter"
)

func getBtcTransactions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var address string
	var from, to int
	var ok bool
	jsonParsed, err := gabs.ParseJSONBuffer(r.Body)
	tmpl := `{error:%v}`

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
		if from, ok = jsonParsed.Path("from").Data().(int); !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf(tmpl, "get from err")))
			return
		}
	}

	if exists := jsonParsed.Exists("to"); exists {
		if from, ok = jsonParsed.Path("to").Data().(int); !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf(tmpl, "get to err")))
			return
		}
	}

	if from > to {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(tmpl, "from gt to")))
		return
	}

	if resp, err := http.Get(fmt.Sprintf("%v/insight-api/addrs/%v/txs?from=%v&to=%v", globalConfig.insight, address, from, to)); err == nil {
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
