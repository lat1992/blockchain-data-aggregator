package services

import "github.com/lat1992/blockchain-data-aggregator/externals"

type Pipeline struct {
	coingecko  externals.CoinGeckoAPI
	dataGetter externals.DataGetterService
	clickhosue externals.Database
}

func NewPipeline(cg externals.CoinGeckoAPI, dg externals.DataGetterService, ch externals.Database) *Pipeline {
	return &Pipeline{
		coingecko:  cg,
		dataGetter: dg,
		clickhosue: ch,
	}
}

func (p *Pipeline) Run() {
}
