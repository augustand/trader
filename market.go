package main

import (
	"encoding/json"
	"fmt"
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
			// add a pseudo asset test
			m.priceList["test"] = priceCoinMarketCap{PriceCNY: "7.0", PriceUSD: "1.0"}
			m.muPriceList.Unlock()
		} else {
			log.Println(err)
		}
		<-time.After(10 * time.Second)
	}
}

type assetItem struct {
	Count string `json:"count"`
	Price string `json:"price"`
}

func (m *coinMarketCap) getPrices(symbols []string, unit string) (prices []string) {
	m.muPriceList.Lock()
	defer m.muPriceList.Unlock()
	for k := range symbols {
		if item, ok := m.priceList[symbols[k]]; ok {
			switch unit {
			case "CNY":
				prices = append(prices, item.PriceCNY)
			case "USD":
				prices = append(prices, item.PriceUSD)
			}

		}
	}
	return
}

func (m *coinMarketCap) getAssets(symbols []string, address string) (prices []string) {
	for k := range symbols {
		contract := contractAddresses[symbols[k]]
		address, err := paduint(address, 64)
		if err != nil {
			continue
		}

		data := fmt.Sprintf("0x%v%v", ERC20Signatures[signBalanceOf], address)
		if ret, err := eth_call(contract, data); err == nil {
			prices = append(prices, ret)
		}
	}

	return
}

func init() {
	defaultCoinMarketCap = NewCoinMarketCap()
}

type priceListStruct struct {
	List    []string `json:"list"`
	Unit    string   `json:"unit"`
	Address string   `json:"address"`
}

func priceListHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	l := priceListStruct{}
	dec := json.NewDecoder(r.Body)
	dec.Decode(&l)

	prices := defaultCoinMarketCap.getPrices(l.List, l.Unit)
	assets := defaultCoinMarketCap.getAssets(l.List, l.Address)

	if len(prices) != len(assets) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	items := make([]assetItem, 0, len(prices))
	for k := range prices {
		items = append(items, assetItem{assets[k], prices[k]})
	}
}
