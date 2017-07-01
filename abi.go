package main

import (
	"encoding/hex"
	"log"

	"github.com/xtaci/trader/sha3"
)

const (
	signBalanceOf   = "balanceOf(address)"
	signTotalSupply = "totalSupply()"
	signTransfer    = "transfer(address,uint)"
)

var (
	contractAddresses = make(map[string]string)
	signatures        = []string{signBalanceOf, signTotalSupply, signTransfer}
	ERC20Signatures   = make(map[string]string)
)

func init() {
	for _, sign := range signatures {
		d := sha3.NewKeccak256()
		d.Write([]byte(sign))
		h := d.Sum(nil)
		ERC20Signatures[sign] = hex.EncodeToString(h[0:4])
		log.Println(sign, ERC20Signatures[sign])
	}

	contractAddresses["test"] = "0x11e1268fc0c49ada36151d4fa2cd0d945e062c1a"
}
