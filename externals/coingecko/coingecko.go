package coinGecko

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

type coinGeckoAPI interface {
}

type Client struct {
	httpClient *http.Client
	url        string
	apiKey     string
	tokenIDs   map[string]string
	priceCache sync.Map
}

func New(url, apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: time.Second * 3,
		},
		url:      url,
		tokenIDs: make(map[string]string),
		apiKey:   apiKey,
	}
}

func (c *Client) InitTokenIDs() error {
	response, err := c.getTokenIDs()
	if err != nil {
		return fmt.Errorf("failed to get token ids: %w", err)
	}
	for _, coin := range response {
		c.tokenIDs[coin.Symbol] = coin.ID
	}
	return nil
}

type CoinGeckoCoinsListResponse struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
}

func (c *Client) getTokenIDs() ([]CoinGeckoCoinsListResponse, error) {
	req, _ := http.NewRequest("GET", c.url+"/coins/list", nil)

	req.Header.Add("accept", "application/json")
	// req.Header.Add("x-cg-demo-api-key", c.apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get token ids from coingecko: %w", err)
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			slog.Error("failed to close response body", "err", err)
		}
	}()
	var result []CoinGeckoCoinsListResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return result, nil
}

func (c *Client) GetTokenID(symbol string) string {
	id, exist := c.tokenIDs[strings.ToLower(symbol)]
	if exist {
		return id
	}
	return ""
}

func (c *Client) GetPrice(coin string) (float64, error) {
	price, exist := c.priceCache.Load(coin)
	if exist {
		return price.(float64), nil
	}

	return 0, nil
}
