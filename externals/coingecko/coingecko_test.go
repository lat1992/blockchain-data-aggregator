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
