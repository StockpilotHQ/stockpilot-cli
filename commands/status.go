package commands

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/stockpilot/stockpilot-cli/internal/api"
	"github.com/stockpilot/stockpilot-cli/internal/config"
	"github.com/stockpilot/stockpilot-cli/internal/output"
)

// status is a smart compound command: shows pending orders + low stock items.
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Pending orders and low stock overview",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		client := api.New(cfg.ClientID, cfg.ClientSecret)

		// Fetch pending orders
		ordersData, err := client.Get("/orders", url.Values{
			"page_size": {"100"},
		})
		if err != nil {
			return err
		}
		var orders []map[string]any
		json.Unmarshal(ordersData, &orders)

		// Fetch inventory
		inventoryData, err := client.Get("/inventory", url.Values{
			"page_size": {"1000"},
		})
		if err != nil {
			return err
		}
		var inventory []map[string]any
		json.Unmarshal(inventoryData, &inventory)

		if jsonOutput {
			return output.JSON(map[string]any{
				"orders":    orders,
				"inventory": inventory,
			})
		}

		fmt.Printf("Orders: %d\n\n", len(orders))
		fmt.Printf("Inventory items: %d\n", len(inventory))
		return nil
	},
}
