package services

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/lat1992/blockchain-data-aggregator/externals"
	"github.com/lat1992/blockchain-data-aggregator/internal"
)

type Pipeline struct {
	coingecko        externals.CoinGeckoAPI
	dataGetter       externals.DataGetterService
	clickhosue       externals.Database
	goroutineNum     int
	marketStatsCache *marketStatCache
}

func NewPipeline(cg externals.CoinGeckoAPI, dg externals.DataGetterService, ch externals.Database, gNum int) *Pipeline {
	cg.InitTokenIDs()
	return &Pipeline{
		coingecko:    cg,
		dataGetter:   dg,
		clickhosue:   ch,
		goroutineNum: gNum,
		marketStatsCache: &marketStatCache{
			stats: make(map[string]internal.MarketStat),
		},
	}
}

func (p *Pipeline) Run() {
	slog.Info("pipeline started")
	var wg sync.WaitGroup
	wg.Add(p.goroutineNum + 1)

	go func() {
		defer wg.Done()
		if err := p.dataGetter.ReadDataFromFiles(); err != nil {
			slog.Error("failed to read data from files", "err", err)
		}
	}()

	for i := 0; i < p.goroutineNum; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case record, ok := <-p.dataGetter.Channel():
					if !ok {
						return
					}
					if err := p.GetMarketStats(record); err != nil {
						slog.Error("failed to get market stats", "err", err)
					}
				case <-p.dataGetter.EndChannel():
					return
				}
			}
		}()
	}
	wg.Wait()

	p.clickhosue.InsertMarket(p.marketStatsCache.stats)
	slog.Info("pipeline ended")
}

type marketStatCache struct {
	mutex sync.Mutex
	stats map[string]internal.MarketStat
}

type propsSchema struct {
	CurrencySymbol string `json:"currencySymbol"`
}

type numsSchema struct {
	CurrencyValueDecimal string `json:"currencyValueDecimal"`
}

func (p *Pipeline) GetMarketStats(record internal.Record) error {
	date, err := time.Parse(time.DateTime+".000", record.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to parse timestamp: %w", err)
	}
	y, m, d := date.Date()
	dateString := fmt.Sprintf("%02d-%02d-%d", d, m, y)

	var props propsSchema
	if err := json.Unmarshal([]byte(record.Props), &props); err != nil {
		return fmt.Errorf("failed to unmarshal props: %w", err)
	}
	var nums numsSchema
	if err := json.Unmarshal([]byte(record.Nums), &nums); err != nil {
		return fmt.Errorf("failed to unmarshal props: %w", err)
	}
	amount, err := strconv.ParseFloat(nums.CurrencyValueDecimal, 64)
	if err != nil {
		return fmt.Errorf("failed to parse currency value decimal: %w", err)
	}

	price, err := p.coingecko.GetPrice(props.CurrencySymbol, dateString)
	if err != nil {
		return fmt.Errorf("failed to get price: %w", err)
	}

	key := dateString + "-" + record.ProjectID

	p.marketStatsCache.Update(key, record.ProjectID, date, price, amount)
	return nil
}

func (c *marketStatCache) Update(key, projectID string, date time.Time, price, amount float64) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	ms, exist := c.stats[key]
	if exist {
		c.stats[key] = internal.MarketStat{
			Date:        date,
			ProjectID:   ms.ProjectID,
			NumTx:       ms.NumTx + 1,
			TotalVolume: ms.TotalVolume + (price * amount),
		}
	} else {
		pID, err := strconv.ParseUint(projectID, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse project id: %w", err)
		}
		c.stats[key] = internal.MarketStat{
			Date:        date,
			ProjectID:   pID,
			NumTx:       1,
			TotalVolume: price * amount,
		}
	}
	return nil
}
