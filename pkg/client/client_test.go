package client

import (
	"path/filepath"
	"testing"

	"github.com/bestk/dmxstart_auto_outbound/pkg/config"
	"github.com/bestk/dmxstart_auto_outbound/pkg/logger"
)

func TestLoginWithAutoOCR(t *testing.T) {
	logger.Init()
	// 创建带有账号信息的配置
	configPath := filepath.Join("..", "..", "config.yaml")
	config, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if config.Account == "" || config.Password == "" {
		t.Fatalf("Account or password is empty")
	}

	client := NewClient(config)
	client.SetLogger(logger.Logger)

	// 执行登录测试
	resp, err := client.LoginWithAutoOCR(config.Account, config.Password)
	if err != nil {
		t.Errorf("LoginWithAutoOCR() error = %v", err)
		return
	}

	// 验证响应
	if !resp.Success {
		t.Errorf("登录失败: %s", resp.ErrorMessage)
		return
	}

	// 验证token
	if resp.Data.Token == "" {
		t.Error("登录成功但未获取到token")
		return
	}

	t.Logf("登录成功，token: %s", resp.Data.Token)
}

func TestValidateSession(t *testing.T) {
	logger.Init()
	configPath := filepath.Join("..", "..", "config.yaml")
	config, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	client := NewClient(config)
	client.SetLogger(logger.Logger)

	err = client.ValidateSession()
	if err != nil {
		t.Errorf("ValidateSession() error = %v", err)
		return
	}

	t.Log("ValidateSession() success")
}

func TestGetWaitingPickOrders(t *testing.T) {
	logger.Init()
	configPath := filepath.Join("..", "..", "config.yaml")
	config, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	client := NewClient(config)
	client.SetLogger(logger.Logger)

	resp, err := client.GetWaitingPickOrders(1, 20, config.CustomerIDs)
	if err != nil {
		t.Errorf("GetWaitingPickOrders() error = %v", err)
		return
	}

	t.Logf("GetWaitingPickOrders() success, total: %d, data: %+v", resp.Total, resp.Data)
}

func TestCreatePickupWave(t *testing.T) {
	logger.Init()
	configPath := filepath.Join("..", "..", "config.yaml")
	config, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	client := NewClient(config)
	client.SetLogger(logger.Logger)

	resp, err := client.CreatePickupWave(true, 1, true, config.CustomerIDs, "[BOT]")
	if err != nil {
		t.Errorf("CreatePickupWave() error = %v", err)
		return
	}

	t.Logf("CreatePickupWave() success, id: %d", resp.Data.ID)
}
