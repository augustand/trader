package main

import (
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
			&cli.StringFlag{
				Name:  "account",
				Value: "0x6bd25eb2e60f5cc47c86abf6ba1b3d03fc74ee27",
				Usage: "address for executing constant queries for tokens",
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
			router.GET("/eth/getGasPrice", getGasPriceHandler)
			router.POST("/eth/getBalance", getBalanceHandler)
			router.POST("/eth/getTransactionCount", getTransactionCountHandler)
			router.POST("/eth/sendRawTransaction", sendRawTransactionHandler)
			router.POST("/eth/tokens/balanceOf", tokenBalanceOfHandler)
			router.POST("/eth/tokens/totalSupply", tokenTotalSupplyHandler)
			router.POST("/eth/tokens/transferABI", transferABIHandler)
			log.Fatal(http.ListenAndServe(globalConfig.listen, router))
			select {}
		},
	}
	app.Run(os.Args)
}
