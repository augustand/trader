package main

import (
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"

	cli "gopkg.in/urfave/cli.v2"
)

type Config struct {
	listen           string
	geth             string
	gasUpdate        time.Duration
	coinMarketCapURL string
	postgres         string
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
			&cli.StringFlag{
				Name:  "coinmarketcapurl",
				Value: "https://api.coinmarketcap.com/v1/ticker/?convert=CNY",
				Usage: "query market price",
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
			&cli.StringFlag{
				Name:  "postgres",
				Value: "host=localhost port=5432 user=postgres dbname=trader sslmode=disable password=qwer1234",
				Usage: "postgres connection string",
			},
		},
		Action: func(c *cli.Context) error {
			globalConfig.listen = c.String("listen")
			globalConfig.geth = c.String("geth")
			globalConfig.gasUpdate = c.Duration("gas_update")
			globalConfig.coinMarketCapURL = c.String("coinmarketcapurl")
			globalConfig.postgres = c.String("postgres")
			log.Println("listen:", globalConfig.listen)
			log.Println("geth:", globalConfig.geth)
			log.Println("gas_update:", globalConfig.gasUpdate)
			log.Println("coinmarketcapurl:", globalConfig.coinMarketCapURL)
			log.Println("postgres:", globalConfig.postgres)

			// init
			db, err := gorm.Open("postgres", globalConfig.postgres)
			if err != nil {
				panic("failed to connect database")
			}
			go update_gas_task()
			defaultETHTXManager = NewTransactionManger(db)

			// webapi
			router := httprouter.New()
			router.GET("/eth/getGasPrice", getGasPriceHandler)
			router.POST("/eth/getBalance", getBalanceHandler)
			router.POST("/eth/getTransactionCount", getTransactionCountHandler)
			router.POST("/eth/sendRawTransaction", sendRawTransactionHandler)
			router.POST("/eth/tokens/balanceOf", tokenBalanceOfHandler)
			router.POST("/eth/tokens/totalSupply", tokenTotalSupplyHandler)
			router.POST("/eth/tokens/transferABI", transferABIHandler)
			router.POST("/market/priceList", priceListHandler)
			log.Fatal(http.ListenAndServe(globalConfig.listen, router))
			select {}
		},
	}
	app.Run(os.Args)
}
