package main

import (
	"bytes"
	"errors"
	"math/big"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/Jeffail/gabs"
	log "github.com/Sirupsen/logrus"
)

type BlockTimeEstimator struct {
	blkPerSecond atomic.Value
}

func (est *BlockTimeEstimator) init() {
	est.blkPerSecond.Store(float64(0.1))
	go est.estimateTask()
}

func (est *BlockTimeEstimator) getBPS() float64 {
	return est.blkPerSecond.Load().(float64)
}

func (est *BlockTimeEstimator) estimateTask() {
	ticker := time.NewTicker(time.Second * 10)
	var lastMesaure time.Time
	var lastBlkId int64

	for {
		if blkId, err := est.getBlock(); err == nil {
			lastBlkId = blkId
			lastMesaure = time.Now()
			break
		}
	}

	for _ = range ticker.C {
		if blkId, err := est.getBlock(); err == nil {
			if blkId > lastBlkId {
				bps := float64(blkId-lastBlkId) / time.Now().Sub(lastMesaure).Seconds()
				est.blkPerSecond.Store(bps)
				lastBlkId = blkId
				lastMesaure = time.Now()
				log.Println("blk/sec:", est.blkPerSecond.Load())
			}
		} else {
			log.Println(err)
		}
	}
}

func (est *BlockTimeEstimator) getBlock() (int64, error) {
	if resp, err := http.Post(globalConfig.geth,
		"application/json",
		bytes.NewBufferString(`{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`)); err == nil {
		jsonParsed, _ := gabs.ParseJSONBuffer(resp.Body)
		value, ok := jsonParsed.Path("result").Data().(string)
		if !ok {
			return 0, errors.New("getBlock failed")
		}

		i := new(big.Int)
		if _, ok := i.SetString(value, 0); ok {
			return i.Int64(), nil
		}
	}
	return 0, errors.New("getBlock failed")
}

var defaultBlockTimeEstimator BlockTimeEstimator
