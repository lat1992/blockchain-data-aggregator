package internal

import "time"

type Record struct {
	Timestamp string
	Event     string
	ProjectID string
	Props     string
	Nums      string
}

type MarketStat struct {
	Date        time.Time
	ProjectID   uint64
	NumTx       uint64
	TotalVolume float64
}
