package main

import (
	"fmt"
	"os"

	"github.com/bestk/dmxstart_auto_outbound/pkg/client"
	"github.com/bestk/dmxstart_auto_outbound/pkg/config"
	"github.com/bestk/dmxstart_auto_outbound/pkg/logger"
)

var banner = `
___________________________________________________

	DMXSTART AUTO OUTBOUND v0.0.1
____________________________________________________

`

func main() {
	logger.Init()

	log := logger.Logger
	log.Info("Starting...")

	log.Info(banner)

	// Load configuration
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create DMXSmart client
	dmxClient := client.NewClient(cfg)

	err = dmxClient.ValidateSession()
	if err != nil {
		log.Fatalf("Failed to validate session: %v", err)
	}

	// Get waiting pick orders
	orders, err := dmxClient.GetWaitingPickOrders(1, 100, cfg.CustomerIDs)
	if err != nil {
		log.Fatalf("Failed to get waiting pick orders: %v", err)
	}

	log.Printf("Waiting pick orders: %d\n", orders.Total)

	if orders.Total > 0 {
		// Create pickup wave
		_, err = dmxClient.CreatePickupWave(true, 1, true, cfg.CustomerIDs, "[DIY OUTBOUND BY BOT]")
		if err != nil {
			log.Fatalf("Failed to create pickup wave: %v", err)
		}
		log.Info("Pickup wave created")
	} else {
		log.Error("No waiting pick orders")
	}

	// 等待用户输入
	log.Info("Press Enter to continue...")
	fmt.Scanln()
	// 退出
	os.Exit(0)
}
