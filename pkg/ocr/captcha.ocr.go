package ocr

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/bestk/dmxstart_auto_outbound/pkg/config"
	"github.com/bestk/dmxstart_auto_outbound/pkg/logger"

	"github.com/go-resty/resty/v2"
)

// RecognizeBase64Image 识别Base64编码的验证码图片
func RecognizeBase64Image(base64Img string) (string, error) {
	// 清理base64前缀（如果存在）
	base64Img = strings.TrimPrefix(base64Img, "data:image/png;base64,")
	base64Img = strings.TrimPrefix(base64Img, "data:image/jpeg;base64,")

	body := map[string]string{
		"base64_image": base64Img,
	}

	result := struct {
		Result string `json:"result"`
	}{}

	resp, err := resty.New().R().
		SetLogger(logger.Logger).
		SetDebug(config.Config.Debug).
		SetBody(body).
		SetResult(&result).
		Post(config.Config.OcrEndpoint)

	if err != nil {
		return "", fmt.Errorf("failed to post request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return result.Result, nil
}

// RecognizeImage 识别图片字节数据
func RecognizeImage(imgData []byte) (string, error) {
	// 转换为base64
	base64Img := base64.StdEncoding.EncodeToString(imgData)

	return RecognizeBase64Image(base64Img)
}

// RecognizeImageFromBuffer 从buffer中识别图片
func RecognizeImageFromBuffer(buffer *bytes.Buffer) (string, error) {
	if buffer == nil {
		return "", fmt.Errorf("buffer is nil")
	}
	return RecognizeImage(buffer.Bytes())
}
