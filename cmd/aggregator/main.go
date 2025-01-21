package main

import (
	"log"

	"github.com/lat1992/blockchain-data-aggregator/config"
	"github.com/lat1992/blockchain-data-aggregator/externals/clickhouse"
	"github.com/lat1992/blockchain-data-aggregator/externals/coingecko"
	"github.com/lat1992/blockchain-data-aggregator/externals/dataGetter"
	"github.com/lat1992/blockchain-data-aggregator/internal/services"
	"github.com/spf13/viper"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Cannot load config:", err)
	}

	dg := dataGetter.New(viper.GetString("DATA_PATH"))
	cg := coingecko.New(viper.GetString("COINGECKO_URL"), viper.GetString("COINGECKO_API_KEY"))
	ch := clickhouse.New("", "")

	pipeline := services.NewPipeline(cg, dg, ch)
	pipeline.Run()
}
