package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Jeffail/gabs"
)

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
