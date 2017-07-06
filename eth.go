package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"encoding/hex"
	"log"

	"github.com/Jeffail/gabs"

	"github.com/xtaci/trader/sha3"
)

const (
	signBalanceOf   = "balanceOf(address)"
	signTotalSupply = "totalSupply()"
	signTransfer    = "transfer(address,uint256)"
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

	contractAddresses["test"] = "0x0125bb97c11a9d6c62814f0645471972804214ed"
	contractAddresses["GNT"] = "0xa74476443119A942dE498590Fe1f2454d7D4aC0d"
	contractAddresses["GNO"] = "0x6810e776880c02933d47db1b9fc05908e5386b96"
	contractAddresses["ICN"] = "0x888666CA69E0f178DED6D75b5726Cee99A87D698"
	contractAddresses["REP"] = "0xC66eA802717bFb9833400264Dd12c2bCeAa34a6d"
	contractAddresses["SNT"] = "0x744d70FDBE2Ba4CF95131626614a1763DF805B9E"
	contractAddresses["DGD"] = "0xE0B7927c4aF23765Cb51314A0E0521A9645F0E2A"
	contractAddresses["BNT"] = "0x1F573D6Fb3F13d689FF844B4cE37794d79a7FF1C"
	contractAddresses["BAT"] = "0x0D8775F648430679A709E98d2b0Cb6250d2887EF"
	contractAddresses["1ST"] = "0xAf30D2a7E90d7DC361c8C4585e9BB7D2F6f15bc7"
	contractAddresses["SNGLS"] = "0xaeC2E87E0A235266D9C5ADc9DEb4b2E29b54D009"
	contractAddresses["ANT"] = "0x960b236A07cf122663c4303350609A66A7B288C0"
	contractAddresses["EDG"] = "0x08711D3B02C8758F2FB3ab4e80228418a7F8e39c"
	contractAddresses["RLC"] = "0x607F4C5BB672230e8672085532f7e901544a7375"
	contractAddresses["MLN"] = "0xBEB9eF514a379B997e0798FDcC901Ee474B6D9A1"
	contractAddresses["SJCX"] = "0xB64ef51C888972c908CFacf59B47C1AfBC0Ab8aC"
	contractAddresses["VSL"] = "0x5c543e7AE0A1104f78406C340E9C64FD9fCE5170"
	contractAddresses["WINGS"] = "0x667088b212ce3d06a1b553a7221E1fD19000d9aF"
	contractAddresses["TKN"] = "0xaAAf91D9b90dF800Df4F55c205fd6989c977E73a"
	contractAddresses["TRST"] = "0xCb94be6f13A1182E4A4B6140cb7bf2025d28e41B"
	contractAddresses["TAAS"] = "0xe7775a6e9bcf904eb39da2b68c5efb4f9360e08c"
	contractAddresses["BCAP"] = "0xFf3519eeeEA3e76F1F699CCcE5E23ee0bdDa41aC"
	contractAddresses["SWT"] = "0xB9e7F8568e08d5659f5D29C4997173d84CdF2607"
	contractAddresses["GUP"] = "0xf7b098298f7c69fc14610bf71d5e02c60792894c"
	contractAddresses["TIME"] = "0x6531f133e6DeeBe7F2dcE5A0441aA7ef330B4e53"
	contractAddresses["PLU"] = "0xD8912C10681D8B21Fd3742244f44658dBA12264E"
	contractAddresses["LUN"] = "0xfa05A73FfE78ef8f1a739473e462c54bae6567D9"
}

func eth_call(to, data string) (string, error) {
	if resp, err := http.Post(globalConfig.geth,
		"application/json",
		bytes.NewBufferString(fmt.Sprintf(`{"jsonrpc":"2.0","method": "eth_call", "params": [{"to": "%v", "data": "%v"}, "latest"], "id": 0}`, to, data))); err == nil {
		jsonParsed, _ := gabs.ParseJSONBuffer(resp.Body)
		value, ok := jsonParsed.Path("result").Data().(string)
		if !ok {
			return jsonParsed.Path("error").String(), nil
		}
		return value, nil
	} else {
		return "", err
	}
}

func paduint(value string, size int) (string, error) {
	if !strings.HasPrefix(value, "0x") {
		return value, errors.New("must start with 0x")
	}

	value = value[2:]
	n := len(value)
	if n%size == 0 {
		return value, nil
	} else {
		return strings.Repeat("0", size-n%size) + value, nil
	}
}
