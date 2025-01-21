package externals

type DataGetterService interface {
	ReadDataFromFiles() error
	Channel() chan Record
	EndChannel() chan bool
}

type CoinGeckoAPI interface {
	GetTokenID(token string) string
	GetPrice(symbol, date string) (float64, error)
}

type Database interface {
}
