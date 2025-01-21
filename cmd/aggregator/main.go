package main

import (
	"log/slog"

	"github.com/lat1992/blockchain-data-aggregator/config"
	"github.com/lat1992/blockchain-data-aggregator/externals/clickhouse"
	"github.com/lat1992/blockchain-data-aggregator/externals/coingecko"
	"github.com/lat1992/blockchain-data-aggregator/externals/dataGetter"
	"github.com/lat1992/blockchain-data-aggregator/internal/services"
	"github.com/spf13/viper"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		slog.Error("Cannot load config", "error", err)
		return
	}

	dg := dataGetter.New(viper.GetString("DATA_PATH"))
	cg := coingecko.New(viper.GetString("COINGECKO_URL"), viper.GetString("COINGECKO_API_KEY"))

	ch, err := clickhouse.New(viper.GetString("CLICKHOUSE_HOSTNAME"), viper.GetString("CLICKHOUSE_DATABASE"), viper.GetString("CLICKHOUSE_USERNAME"), viper.GetString("CLICKHOUSE_PASSWORD"))
	if err != nil {
		slog.Error("Cannot connect to clickhouse", "error", err)
		return
	}

	pipeline := services.NewPipeline(cg, dg, ch)
	pipeline.Run()
}
