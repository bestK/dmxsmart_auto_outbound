package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/bestk/dmxstart_auto_outbound/pkg/config"
	"github.com/bestk/dmxstart_auto_outbound/pkg/ocr"
	"github.com/go-resty/resty/v2"
)

const (
	BaseURL = "https://wms.dmxsmart.com"
)

// Client represents a DMXSmart API client
type Client struct {
	httpClient *resty.Client
	baseURL    string
	config     *config.ConfigStruct
}

// Config holds the configuration for the DMXSmart client
type Config struct {
	AccessToken string
	WarehouseID string
	CustomerIDs []string
	Username    string // 用户名
	Password    string // 密码
}

// NewClient creates a new DMXSmart client
func NewClient(config *config.ConfigStruct) *Client {
	client := &Client{
		httpClient: resty.New().SetDebug(config.Debug),
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

func (c *Client) SetLogger(logger resty.Logger) {
	c.httpClient.SetLogger(logger)
}

// validateSession validates the session
func (c *Client) ValidateSession() error {
	url := fmt.Sprintf("%s/api/user/getUserInfo", c.baseURL)

	resp, err := c.httpClient.R().
		Get(url)

	if err != nil {
		return fmt.Errorf("failed to validate session: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	var respData Response
	err = json.Unmarshal(resp.Body(), &respData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !respData.Success {
		return fmt.Errorf("validate session failed: %s", respData.ErrorMessage)
	}

	return nil
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

// GetCaptcha retrieves a captcha image
func (c *Client) GetCaptcha() (*CaptchaResponse, error) {
	url := fmt.Sprintf("%s/api/login/captcha", c.baseURL)

	authHeader := c.httpClient.Header.Get("Authorization")
	if authHeader != "" {
		c.httpClient.Header.Del("Authorization")
	}

	resp, err := c.httpClient.R().
		SetQueryParam("lang", "zh-CN").
		SetHeader("Accept", "application/json, text/plain, */*").
		SetHeader("Referer", fmt.Sprintf("%s/user/login", c.baseURL)).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to get captcha: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	var captchaResp CaptchaResponse
	if err := json.Unmarshal(resp.Body(), &captchaResp); err != nil {
		return nil, fmt.Errorf("failed to parse captcha response: %w", err)
	}

	return &captchaResp, nil
}

// Login 执行登录操作
func (c *Client) Login(username, password, captcha, uuid string) (*LoginResponse, error) {
	// 加密密码
	encryptedPassword, err := encryptPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt password: %w", err)
	}

	// 构建登录请求
	loginReq := LoginRequest{
		Username:    username,
		Password:    encryptedPassword,
		Captcha:     captcha,
		UUID:        uuid,
		LoginType:   "USERNAME",
		DeviceToken: nil,
		Lang:        "zh-CN",
	}

	url := fmt.Sprintf("%s/api/login/authenticate", c.baseURL)

	// 发送登录请求
	resp, err := c.httpClient.R().
		SetHeader("Accept", "application/json, text/plain, */*").
		SetHeader("Referer", fmt.Sprintf("%s/user/login", c.baseURL)).
		SetHeader("Cookie", "locale=zh-CN").
		SetBody(loginReq).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("failed to send login request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	// 解析响应
	var loginResp LoginResponse
	if err := json.Unmarshal(resp.Body(), &loginResp); err != nil {
		return nil, fmt.Errorf("failed to parse login response: %w", err)
	}

	// 如果登录成功，更新客户端的 token
	if loginResp.Success {
		c.httpClient.SetHeader("Authorization", "Bearer "+loginResp.Data.Token)
		c.config.AccessToken = loginResp.Data.Token
		config.SaveConfig()
	}

	return &loginResp, nil
}

// LoginWithAutoOCR 自动识别验证码并登录，失败时最多重试3次
func (c *Client) LoginWithAutoOCR(username, password string) (*LoginResponse, error) {
	maxRetries := 3
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		// 获取验证码
		captchaResp, err := c.GetCaptcha()
		if err != nil {
			lastErr = fmt.Errorf("attempt %d: failed to get captcha: %w", attempt+1, err)
			continue
		}

		// 从base64图片中提取图片数据
		imgData := captchaResp.Data.Img
		imgData = strings.TrimPrefix(imgData, "data:image/png;base64,")
		imgData = strings.TrimPrefix(imgData, "data:image/jpeg;base64,")

		// 识别验证码
		captchaText, err := ocr.RecognizeBase64Image(imgData)
		if err != nil {
			lastErr = fmt.Errorf("attempt %d: failed to recognize captcha: %w", attempt+1, err)
			continue
		}

		// 使用识别出的验证码进行登录
		resp, err := c.Login(username, password, captchaText, captchaResp.Data.UUID)
		if err != nil {
			lastErr = fmt.Errorf("attempt %d: failed to login: %w", attempt+1, err)
			continue
		}

		// 如果登录成功，直接返回结果
		if resp.Success {
			return resp, nil
		}

		// 如果登录失败但没有报错，记录失败原因
		lastErr = fmt.Errorf("attempt %d: login failed: %s", attempt+1, resp.ErrorMessage)
	}

	return nil, fmt.Errorf("login failed after %d attempts, last error: %w", maxRetries, lastErr)
}
