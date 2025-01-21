package externals

import "github.com/lat1992/blockchain-data-aggregator/internal"

type DataGetterService interface {
	ReadDataFromFiles() error
	Channel() chan internal.Record
	EndChannel() chan bool
}

type CoinGeckoAPI interface {
	InitTokenIDs() error
	GetTokenID(token string) string
	GetPrice(symbol, date string) (float64, error)
}

type Database interface {
	InsertMarket(stats map[string]internal.MarketStat) error
}
