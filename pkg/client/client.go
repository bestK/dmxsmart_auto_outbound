package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

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
	config     *config.ConfigStruct
}

// NewClient creates a new DMXSmart client
func NewClient(config *config.ConfigStruct) *Client {
	client := &Client{
		httpClient: resty.New().
			SetDebug(config.Debug).
			SetBaseURL(BaseURL).
			// 设置超时
			SetTimeout(time.Duration(config.Timeout) * time.Second).
			// 设置重试
			SetRetryCount(3).
			SetRetryWaitTime(5 * time.Second).
			SetRetryMaxWaitTime(20 * time.Second).
			// 设置TLS配置
			SetTLSClientConfig(&tls.Config{
				MinVersion: tls.VersionTLS12,
			}),
		config: config,
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
	url := "/api/user/getUserInfo"

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
func (c *Client) GetWaitingPickOrders(page, pageSize int, customerIds []int) (WaitingPickOrderResponse, error) {
	urlStr := "/api/tenant/outbound/pickupwave/listWaitingPickOrder"

	params := url.Values{}
	params.Set("current", fmt.Sprintf("%d", page))
	params.Set("pageSize", fmt.Sprintf("%d", pageSize))
	params.Set("warehouseId", c.config.WarehouseID)
	params.Set("keywordType", "referenceId")
	params.Set("timeType", "createTime")
	params.Set("keyword", "")

	for _, customerId := range customerIds {
		params.Add("customerIds[]", strconv.Itoa(customerId))
	}

	fullURL := fmt.Sprintf("%s?%s", urlStr, params.Encode())

	var result WaitingPickOrderResponse

	resp, err := c.httpClient.R().
		SetBody(map[string]string{
			"keyword": "",
		}).
		SetResult(&result).
		Post(fullURL)

	if err != nil {
		return WaitingPickOrderResponse{}, fmt.Errorf("failed to get waiting pick orders: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return WaitingPickOrderResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	if !result.Success {
		return WaitingPickOrderResponse{}, fmt.Errorf("get waiting pick orders failed: %s", result.ErrorMessage)
	}

	return result, nil
}

// CreatePickupWave creates a new pickup wave
func (c *Client) CreatePickupWave(isAll bool, pickupType int, isOutbound bool, customerIds []int, remark string) (CreatePickupWaveResponse, error) {
	urlStr := "/api/tenant/outbound/pickupwave/createPickupWave"

	params := url.Values{}
	params.Set("isAll", fmt.Sprintf("%t", isAll))
	params.Set("pickupType", fmt.Sprintf("%d", pickupType))
	params.Set("isOutbound", fmt.Sprintf("%t", isOutbound))
	params.Set("warehouseId", c.config.WarehouseID)
	params.Set("remark", remark)

	for _, customerId := range customerIds {
		params.Add("customerIds[]", strconv.Itoa(customerId))
	}

	fullURL := fmt.Sprintf("%s?%s", urlStr, params.Encode())

	var result CreatePickupWaveResponse

	resp, err := c.httpClient.R().
		SetResult(&result).
		Post(fullURL)

	if err != nil {
		return CreatePickupWaveResponse{}, fmt.Errorf("failed to create pickup wave: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return CreatePickupWaveResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return result, nil
}

// GetCaptcha retrieves a captcha image
func (c *Client) GetCaptcha() (*CaptchaResponse, error) {
	url := "/api/login/captcha"

	authHeader := c.httpClient.Header.Get("Authorization")
	if authHeader != "" {
		c.httpClient.Header.Del("Authorization")
	}

	var result CaptchaResponse

	resp, err := c.httpClient.R().
		SetQueryParam("lang", "zh-CN").
		SetHeader("Accept", "application/json, text/plain, */*").
		SetHeader("Referer", fmt.Sprintf("%s/user/login", BaseURL)).
		SetResult(&result).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to get captcha: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return &result, nil
}

// Login 执行登录操作
func (c *Client) Login(username, password, captcha, uuid string) (*LoginResponse, error) {
	url := "/api/login/authenticate"

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

	var result LoginResponse
	// 发送登录请求
	resp, err := c.httpClient.R().
		SetHeader("Accept", "application/json, text/plain, */*").
		SetHeader("Referer", fmt.Sprintf("%s/user/login", BaseURL)).
		SetHeader("Cookie", "locale=zh-CN").
		SetBody(loginReq).
		SetResult(&result).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("failed to send login request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	// 如果登录成功，更新客户端的 token
	if result.Success {
		c.httpClient.SetHeader("Authorization", "Bearer "+result.Data.Token)
		c.config.AccessToken = result.Data.Token
		config.SaveConfig()
	}

	return &result, nil
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
