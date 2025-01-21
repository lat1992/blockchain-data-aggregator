package coinGecko

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTokenID(t *testing.T) {
	client := New("https://api.coingecko.com/api/v3", "demo")

	err := client.InitTokenIDs()
	assert.NoError(t, err)
	assert.Equal(t, "sunflower-land", client.GetTokenID("SFL"))
}

func TestGetTokenPrice(t *testing.T) {
	client := New("https://api.coingecko.com/api/v3", "CG-GKvKPioBeTZQzkgGz4AKwgEe")

	err := client.InitTokenIDs()
	assert.NoError(t, err)

	testCases := []struct {
		name   string
		symbol string
		date   string
		price  float64
		err    error
	}{
		{
			name:   "normal case",
			symbol: "SFL",
			date:   "01-01-2025",
			price:  0.046057701457628754,
		},
		{
			name:   "invalid symbol",
			symbol: "invalid",
			date:   "01-01-2025",
			price:  0,
		},
		{
			name:   "invalid date",
			symbol: "SFL",
			date:   "invalid",
			price:  0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			price, err := client.GetPrice(tc.symbol, tc.date)
			if tc.err != nil {
				assert.Error(t, err)
			}
			assert.Equal(t, tc.price, price)
		})
	}
}
