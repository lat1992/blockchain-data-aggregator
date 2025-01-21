package externals

type DataGetterService interface {
	ReadDataFromFiles(path string) error
	Channel() chan Record
}

type CoinGeckoAPI interface {
	GetTokenID(token string) string
	GetPrice(symbol, date string) (float64, error)
}
