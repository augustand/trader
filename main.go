package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"

	cli "gopkg.in/urfave/cli.v2"
)

type Config struct {
	listen    string
	geth      string
	gasUpdate time.Duration
}

var globalConfig Config

func main() {
	app := &cli.App{
		Name:    "trader",
		Usage:   "trader for the wallet",
		Version: "2.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "listen",
				Value: ":8888",
				Usage: "listening address:port",
			},
			&cli.StringFlag{
				Name:  "geth",
				Value: "http://127.0.0.1:8545",
				Usage: "geth nodes",
			},
			&cli.DurationFlag{
				Name:  "gas_update",
				Value: 10 * time.Second,
				Usage: "set gas update period",
			},
		},
		Action: func(c *cli.Context) error {
			globalConfig.listen = c.String("listen")
			globalConfig.geth = c.String("geth")
			globalConfig.gasUpdate = c.Duration("gas_update")
			log.Println("listen:", globalConfig.listen)
			log.Println("geth:", globalConfig.geth)
			log.Println("gas_update:", globalConfig.gasUpdate)

			// init
			go update_gas_task()
			// webapi
			router := httprouter.New()
			router.GET("/gasPrice", gasPrice)
			log.Fatal(http.ListenAndServe(globalConfig.listen, router))
			select {}
		},
	}
	app.Run(os.Args)
}

func gasPrice(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp, err := http.Post(globalConfig.geth,
		"application/json",
		bytes.NewBufferString(`{"jsonrpc":"2.0","method":"eth_gasPrice","params":[],"id":73}`))

	if err != nil {
		log.Println("call:", err)
	}

	io.Copy(w, resp.Body)
}
