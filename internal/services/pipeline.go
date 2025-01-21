package services

import (
	"github.com/lat1992/blockchain-data-aggregator/externals"
	"github.com/lat1992/blockchain-data-aggregator/internal"
)

type Pipeline struct {
	coingecko  externals.CoinGeckoAPI
	dataGetter externals.DataGetterService
	clickhosue externals.Database
}

func NewPipeline(cg externals.CoinGeckoAPI, dg externals.DataGetterService, ch externals.Database) *Pipeline {
	cg.InitTokenIDs()
	return &Pipeline{
		coingecko:  cg,
		dataGetter: dg,
		clickhosue: ch,
	}
}

func (p *Pipeline) Run() {
	go p.dataGetter.ReadDataFromFiles()

	var stats []internal.MarketStat
	for {
		select {
		case record := <-p.dataGetter.Channel():
			stats = append(stats, p.GetMarketStats(record))
		case <-p.dataGetter.EndChannel():
			p.clickhosue.InsertMarket(stats)
			return
		}
	}
}

func (p *Pipeline) GetMarketStats(record internal.Record) internal.MarketStat {
	return internal.MarketStat{}
}
