package internal

type Record struct {
	Timestamp string
	Event     string
	ProjectId string
	Props     string
	Nums      string
}

type MarketStat struct {
	Date        string
	ProjectId   uint64
	NumTx       uint64
	TotalVolume float64
}
