package coingecko

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

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

func (c *Client) buildAndSendRequest(endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("x-cg-demo-api-key", c.apiKey)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get token ids from coingecko: %w", err)
	}

	return res, nil
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

type coinGeckoCoinsListResponse struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
}

func (c *Client) getTokenIDs() ([]coinGeckoCoinsListResponse, error) {
	res, err := c.buildAndSendRequest(c.url + "/coins/list")
	if err != nil {
		return nil, fmt.Errorf("failed to get token ids from coingecko: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			slog.Error("failed to close response body", "err", err)
		}
	}()

	var result []coinGeckoCoinsListResponse
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

func (c *Client) GetPrice(symbol, date string) (float64, error) {
	key := symbol + "-" + date
	price, exist := c.priceCache.Load(key)
	if exist {
		return price.(float64), nil
	}
	result, err := c.getPriceFromSource(symbol, date)
	if err != nil {
		return 0, fmt.Errorf("failed to get price from source: %w", err)
	}
	c.priceCache.Store(key, result)
	return result, nil
}

type coinGeckoCoinsHistoryResponse struct {
	MarketData struct {
		CurrentPrice struct {
			USD float64 `json:"usd"`
		} `json:"current_price"`
	} `json:"market_data"`
}

func (c *Client) getPriceFromSource(symbol, date string) (float64, error) {
	id := c.GetTokenID(symbol)
	if id == "" {
		slog.Error("id not found with token symbol", "symbol", symbol)
		return 0, nil
	}

	_, err := time.Parse("02-01-2006", date)
	if err != nil {
		slog.Error("date format not valid", "date", date)
		return 0, nil
	}

	res, err := c.buildAndSendRequest(c.url + "/coins/" + id + "/history?date=" + date + "&localization=false")
	if err != nil {
		return 0, fmt.Errorf("failed to get token price from coingecko: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			slog.Error("failed to close response body", "err", err)
		}
	}()

	var result coinGeckoCoinsHistoryResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode response body: %w", err)
	}

	return result.MarketData.CurrentPrice.USD, nil
}
