package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
)

const (
	coinMarketCapURL = "https://api.coinmarketcap.com/v1/ticker/?convert=CNY"
)

type priceCoinMarketCap struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	Symbol            string `json:"symbol"`
	Rank              int    `json:"rank"`
	PriceUSD          string `json:"price_usd"`
	PriceBTC          string `json:"price_btc"`
	VolumeUSD24H      string `json:"24h_volume_usd"`
	MarketCapUSD      string `json:"market_cap_usd"`
	AvailableSupply   string `json:"available_supply"`
	TotalSupply       string `json:"total_supply"`
	PercentChanage1H  string `json:"percent_change_1h"`
	PercentChanage24H string `json:"percent_change_24h"`
	PercentChanage7D  string `json:"percent_change_7d"`
	LastUpdated       string `json:"last_updated"`
	PriceCNY          string `json:"price_cny"`
	VolumeCNY24H      string `json:"24h_volume_cny"`
	MarketCapCNY      string `json:"market_cap_cny"`
}

var defaultCoinMarketCap *coinMarketCap

type coinMarketCap struct {
	priceList   map[string]priceCoinMarketCap // symbol -> item
	muPriceList sync.Mutex
}

func newCoinMarketCap() *coinMarketCap {
	m := &coinMarketCap{}
	go m.updateTask()
	return m
}

func (m *coinMarketCap) updateTask() {
	for {
		if resp, err := http.Get(coinMarketCapURL); err == nil {
			var list []priceCoinMarketCap
			dec := json.NewDecoder(resp.Body)
			dec.Decode(&list)
			m.muPriceList.Lock()
			m.priceList = make(map[string]priceCoinMarketCap)
			for k := range list {
				m.priceList[list[k].Symbol] = list[k]
			}
			// add a pseudo asset test
			m.priceList["test"] = priceCoinMarketCap{PriceCNY: "7.0", PriceUSD: "1.0"}
			m.muPriceList.Unlock()
		} else {
			log.Println(err)
		}
		<-time.After(5 * time.Second)
	}
}

func (m *coinMarketCap) getJSON() ([]byte, error) {
	m.muPriceList.Lock()
	defer m.muPriceList.Unlock()
	return json.Marshal(m.priceList)
}

func init() {
	defaultCoinMarketCap = newCoinMarketCap()
}

func priceListHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if bts, err := defaultCoinMarketCap.getJSON(); err == nil {
		w.Write(bts)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
