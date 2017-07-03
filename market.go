package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
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

func NewCoinMarketCap() *coinMarketCap {
	m := &coinMarketCap{}
	go m.update_task()
	return m
}

func (m *coinMarketCap) update_task() {
	for {
		if resp, err := http.Get(coinMarketCapURL); err == nil {
			var list []priceCoinMarketCap
			dec := json.NewDecoder(resp.Body)
			dec.Decode(&list)
			m.muPriceList.Lock()
			m.priceList = make(map[string]priceCoinMarketCap)
			for k := range list {
				fmt.Println(list[k].Symbol, list[k].PriceCNY)
				m.priceList[list[k].Symbol] = list[k]
			}
			m.muPriceList.Unlock()
		} else {
			log.Println(err)
		}
		<-time.After(10 * time.Second)
	}
}

func init() {
	defaultCoinMarketCap = NewCoinMarketCap()
}
