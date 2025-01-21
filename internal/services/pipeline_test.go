package services

import (
	"fmt"
	"testing"
	"time"

	"github.com/lat1992/blockchain-data-aggregator/internal"
	"github.com/lat1992/blockchain-data-aggregator/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/mock"
)

func TestNewPipeline(t *testing.T) {
	mockCG := new(mocks.CoinGeckoAPI)
	mockDG := new(mocks.DataGetterService)
	mockDB := new(mocks.Database)

	mockCG.On("InitTokenIDs").Return(nil)

	pipeline := NewPipeline(mockCG, mockDG, mockDB, 1)

	assert.NotNil(t, pipeline)
	assert.Equal(t, mockCG, pipeline.coingecko)
	assert.Equal(t, mockDG, pipeline.dataGetter)
	assert.Equal(t, mockDB, pipeline.clickhosue)
	assert.Equal(t, 1, pipeline.goroutineNum)

	mockCG.AssertExpectations(t)
}

func TestPipeline_Run(t *testing.T) {
	mockCG := new(mocks.CoinGeckoAPI)
	mockDG := new(mocks.DataGetterService)
	mockDB := new(mocks.Database)

	recordChan := make(chan internal.Record, 1)
	endChan := make(chan bool, 1)

	mockCG.On("InitTokenIDs").Return(nil)
	mockDG.On("ReadDataFromFiles").Return(nil)
	mockDG.On("Channel").Return(recordChan)
	mockDG.On("EndChannel").Return(endChan)
	mockDB.On("InsertMarket", mock.Anything).Return(nil)

	pipeline := NewPipeline(mockCG, mockDG, mockDB, 1)

	go func() {
		recordChan <- internal.Record{
			Timestamp: "2024-01-01 12:00:00.000",
			ProjectID: "1234",
			Props:     `{"currencySymbol":"BTC"}`,
			Nums:      `{"currencyValueDecimal":"1.5"}`,
		}
		endChan <- true
		close(recordChan)
	}()

	mockCG.On("GetPrice", "BTC", "01-01-2024").Return(50000.0, nil)

	pipeline.Run()

	mockCG.AssertExpectations(t)
	mockDG.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

func TestPipeline_GetMarketStats(t *testing.T) {
	testCases := []struct {
		name   string
		record internal.Record
		setup  func(*mocks.CoinGeckoAPI)
		stats  map[string]internal.MarketStat
		err    error
	}{
		{
			name: "successful case",
			record: internal.Record{
				Timestamp: "2024-01-01 12:00:00.000",
				ProjectID: "1234",
				Props:     `{"currencySymbol":"BTC"}`,
				Nums:      `{"currencyValueDecimal":"1.5"}`,
			},
			setup: func(mockCG *mocks.CoinGeckoAPI) {
				mockCG.On("GetPrice", "BTC", "01-01-2024").Return(50000.0, nil)
			},
			stats: map[string]internal.MarketStat{
				"01-01-2024-1234": {
					Date:        time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
					ProjectID:   1234,
					NumTx:       1,
					TotalVolume: 75000.0, // 1.5 * 50000.0
				},
			},
		},
		{
			name: "invalid timestamp",
			record: internal.Record{
				Timestamp: "invalid-timestamp",
				ProjectID: "1234",
				Props:     `{"currencySymbol":"BTC"}`,
				Nums:      `{"currencyValueDecimal":"1.5"}`,
			},
			setup: func(mockCG *mocks.CoinGeckoAPI) {},
			err:   assert.AnError,
		},
		{
			name: "invalid props JSON",
			record: internal.Record{
				Timestamp: "2024-01-01 12:00:00.000",
				ProjectID: "1234",
				Props:     `invalid json`,
				Nums:      `{"currencyValueDecimal":"1.5"}`,
			},
			setup: func(mockCG *mocks.CoinGeckoAPI) {},
			err:   assert.AnError,
		},
		{
			name: "invalid nums JSON",
			record: internal.Record{
				Timestamp: "2024-01-01 12:00:00.000",
				ProjectID: "1234",
				Props:     `{"currencySymbol":"BTC"}`,
				Nums:      `invalid json`,
			},
			setup: func(mockCG *mocks.CoinGeckoAPI) {},
			err:   assert.AnError,
		},
		{
			name: "invalid currency value",
			record: internal.Record{
				Timestamp: "2024-01-01 12:00:00.000",
				ProjectID: "1234",
				Props:     `{"currencySymbol":"BTC"}`,
				Nums:      `{"currencyValueDecimal":"invalid"}`,
			},
			setup: func(mockCG *mocks.CoinGeckoAPI) {},
			err:   assert.AnError,
		},
		{
			name: "coingecko error",
			record: internal.Record{
				Timestamp: "2024-01-01 12:00:00.000",
				ProjectID: "1234",
				Props:     `{"currencySymbol":"BTC"}`,
				Nums:      `{"currencyValueDecimal":"1.5"}`,
			},
			setup: func(mockCG *mocks.CoinGeckoAPI) {
				mockCG.On("GetPrice", "BTC", "01-01-2024").Return(0.0, fmt.Errorf("coingecko error"))
			},
			err: assert.AnError,
		},
		{
			name: "multiple transactions for same project and date",
			record: internal.Record{
				Timestamp: "2024-01-01 12:00:00.000",
				ProjectID: "1234",
				Props:     `{"currencySymbol":"BTC"}`,
				Nums:      `{"currencyValueDecimal":"1.5"}`,
			},
			setup: func(mockCG *mocks.CoinGeckoAPI) {
				mockCG.On("GetPrice", "BTC", "01-01-2024").Return(50000.0, nil)
			},
			stats: map[string]internal.MarketStat{
				"01-01-2024-1234": {
					Date:        time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
					ProjectID:   1234,
					NumTx:       2,
					TotalVolume: 150000.0, // (1.5 * 50000.0) * 2
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCG := new(mocks.CoinGeckoAPI)
			mockDG := new(mocks.DataGetterService)
			mockDB := new(mocks.Database)

			mockCG.On("InitTokenIDs").Return(nil)
			tc.setup(mockCG)

			pipeline := NewPipeline(mockCG, mockDG, mockDB, 1)

			if tc.name == "multiple transactions for same project and date" {
				err := pipeline.GetMarketStats(tc.record)
				assert.NoError(t, err)
			}

			err := pipeline.GetMarketStats(tc.record)

			if tc.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				for key, expectedStat := range tc.stats {
					actualStat, exists := pipeline.marketStatsCache.stats[key]
					assert.True(t, exists)
					assert.Equal(t, expectedStat.Date.Unix(), actualStat.Date.Unix())
					assert.Equal(t, expectedStat.ProjectID, actualStat.ProjectID)
					assert.Equal(t, expectedStat.NumTx, actualStat.NumTx)
					assert.Equal(t, expectedStat.TotalVolume, actualStat.TotalVolume)
				}
			}

			// Verify mock expectations
			mockCG.AssertExpectations(t)
		})
	}
}
