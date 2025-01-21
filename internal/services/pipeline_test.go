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
		},
		{
			name: "invalid timestamp",
			record: internal.Record{
				Timestamp: "invalid",
				ProjectID: "1234",
				Props:     `{"currencySymbol":"BTC"}`,
				Nums:      `{"currencyValueDecimal":"1.5"}`,
			},
			setup: func(mockCG *mocks.CoinGeckoAPI) {},
			err:   assert.AnError,
		},
		{
			name: "invalid project ID",
			record: internal.Record{
				Timestamp: "2024-01-01 12:00:00.000",
				ProjectID: "invalid",
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
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockCG := new(mocks.CoinGeckoAPI)
			mockDG := new(mocks.DataGetterService)
			mockDB := new(mocks.Database)

			mockCG.On("InitTokenIDs").Return(nil)
			tt.setup(mockCG)

			pipeline := NewPipeline(mockCG, mockDG, mockDB, 1)
			err := pipeline.GetMarketStats(tt.record)

			if tt.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				date, _ := time.Parse(time.DateTime+".000", tt.record.Timestamp)
				y, m, d := date.Date()
				dateString := fmt.Sprintf("%02d-%02d-%d", d, m, y)

				stat, exists := pipeline.marketStatsCache.stats[dateString+"+"+tt.record.ProjectID]
				assert.True(t, exists)
				assert.Equal(t, date, stat.Date)
				assert.Equal(t, uint64(1234), stat.ProjectID)
				assert.Equal(t, uint64(1), stat.NumTx)
				assert.Equal(t, 75000.0, stat.TotalVolume)
			}

			mockCG.AssertExpectations(t)
		})
	}
}
