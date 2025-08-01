package main

import (
	"fmt"
	"log"

	"github.com/bestk/dmxstart_auto_outbound/pkg/client"
	"github.com/bestk/dmxstart_auto_outbound/pkg/config"
	"github.com/bestk/dmxstart_auto_outbound/pkg/logger"
)

func main() {
	logger.Init()

	logger.Logger.Info("Starting application")

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
	orders, err := dmxClient.GetWaitingPickOrders(1, 20, cfg.CustomerIDs)
	if err != nil {
		log.Fatalf("Failed to get waiting pick orders: %v", err)
	}
	fmt.Printf("Waiting pick orders: %s\n", orders)

	// Create pickup wave
	// result, err := dmxClient.CreatePickupWave(true, 1, true)
	// if err != nil {
	// 	log.Fatalf("Failed to create pickup wave: %v", err)
	// }
	// fmt.Printf("Pickup wave created: %s\n", result)
}
