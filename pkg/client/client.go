package client

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

const (
	BaseURL = "https://wms.dmxsmart.com"
)

// Client represents a DMXSmart API client
type Client struct {
	httpClient *resty.Client
	baseURL    string
	config     *Config
}

// Config holds the configuration for the DMXSmart client
type Config struct {
	AccessToken string
	WarehouseID string
	CustomerIDs []string
}

// NewClient creates a new DMXSmart client
func NewClient(config *Config) *Client {
	client := &Client{
		httpClient: resty.New(),
		baseURL:    BaseURL,
		config:     config,
	}

	// Set default headers
	client.httpClient.SetHeaders(map[string]string{
		"Accept":             "*/*",
		"Accept-Language":    "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7",
		"Authorization":      "Bearer " + config.AccessToken,
		"Cache-Control":      "no-cache",
		"Connection":         "keep-alive",
		"Content-Type":       "application/json",
		"Origin":             BaseURL,
		"Pragma":             "no-cache",
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "same-origin",
		"User-Agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36",
		"sec-ch-ua":          `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": `"macOS"`,
	})

	return client
}

// GetWaitingPickOrders retrieves the list of waiting pick orders
func (c *Client) GetWaitingPickOrders(page, pageSize int) (string, error) {
	url := fmt.Sprintf("%s/api/tenant/outbound/pickupwave/listWaitingPickOrder", c.baseURL)

	resp, err := c.httpClient.R().
		SetQueryParams(map[string]string{
			"current":     fmt.Sprintf("%d", page),
			"pageSize":    fmt.Sprintf("%d", pageSize),
			"warehouseId": c.config.WarehouseID,
			"keywordType": "referenceId",
			"timeType":    "createTime",
			"keyword":     "",
		}).
		SetBody(map[string]string{
			"keyword": "",
		}).
		Post(url)

	if err != nil {
		return "", fmt.Errorf("failed to get waiting pick orders: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return resp.String(), nil
}

// CreatePickupWave creates a new pickup wave
func (c *Client) CreatePickupWave(isAll bool, pickupType int, isOutbound bool) (string, error) {
	url := fmt.Sprintf("%s/api/tenant/outbound/pickupwave/createPickupWave", c.baseURL)

	resp, err := c.httpClient.R().
		SetQueryParams(map[string]string{
			"isAll":       fmt.Sprintf("%t", isAll),
			"pickupType":  fmt.Sprintf("%d", pickupType),
			"isOutbound":  fmt.Sprintf("%t", isOutbound),
			"warehouseId": c.config.WarehouseID,
		}).
		Post(url)

	if err != nil {
		return "", fmt.Errorf("failed to create pickup wave: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return resp.String(), nil
}
