package main

import (
	"fmt"
	"log"

	"github.com/bestk/dmxstart_auto_outbound/pkg/client"
	"github.com/bestk/dmxstart_auto_outbound/pkg/config"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create DMXSmart client
	dmxClient := client.NewClient(&client.Config{
		AccessToken: cfg.AccessToken,
		WarehouseID: cfg.WarehouseID,
		CustomerIDs: cfg.CustomerIDs,
	})

	// Get waiting pick orders
	orders, err := dmxClient.GetWaitingPickOrders(1, 20)
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
