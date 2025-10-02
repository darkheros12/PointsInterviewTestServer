package client

import (
	"PointsInterviewTestServer/internal/models"
	"PointsInterviewTestServer/internal/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	timeoutDuration = 200
)

type TaxRateProvider interface {
	GetTaxBrackets(ctx context.Context, year int) ([]models.TaxBracket, error)
}

type HTTPTaxClient struct {
	baseURL    string
	httpClient *http.Client
	cache      *utils.MemoryCache
	maxRetries int
}

func NewHTTPTaxClient(baseURL string, cache *utils.MemoryCache) *HTTPTaxClient {
	return &HTTPTaxClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 1 * time.Minute,
		},
		cache:      cache,
		maxRetries: 2,
	}
}

func (c *HTTPTaxClient) GetTaxBrackets(ctx context.Context, year int) ([]models.TaxBracket, error) {
	// Check if the year is being cached in-memory
	if cacheResult, ok := c.cache.Get(year); ok {
		return cacheResult, nil
	}

	// process to request the mock api for year/bands
	var err error
	var resp *http.Response
	url := fmt.Sprintf("%s/tax-calculator/tax-year/%d", c.baseURL, year)

	for attempt := 0; attempt < c.maxRetries; attempt++ {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		resp, err = c.httpClient.Do(req)
		if err != nil {
			time.Sleep(time.Duration(timeoutDuration*attempt) * time.Millisecond)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("mock api status: %d", resp.StatusCode)
			if resp.StatusCode >= 500 {
				time.Sleep(time.Duration(timeoutDuration*attempt) * time.Millisecond)
				continue
			}
			return nil, err
		}

		var mockApiResp models.MockApiTaxResponse
		if err = json.NewDecoder(resp.Body).Decode(&mockApiResp); err != nil {
			log.Printf("Error with Unmarshal JSON response: %v", err)
			return nil, err
		}
		if len(mockApiResp.TaxBrackets) == 0 {
			log.Printf("empty tax brackets")
			return nil, errors.New("empty tax brackets")
		}
		c.cache.Set(year, mockApiResp.TaxBrackets)
		return mockApiResp.TaxBrackets, nil
	}

	return nil, err
}
