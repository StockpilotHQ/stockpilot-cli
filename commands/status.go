package commands

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/StockpilotHQ/stockpilot-cli/internal/api"
	"github.com/StockpilotHQ/stockpilot-cli/internal/config"
	"github.com/StockpilotHQ/stockpilot-cli/internal/output"
)

// status is a smart compound command: shows open orders + inventory count.
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Open orders and inventory overview",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		client := api.New(cfg.ClientID, cfg.ClientSecret)

		// Fetch open orders — use count from wrapper, not len(results)
		ordersData, err := client.Get("/orders", url.Values{
			"page_size": {"1"},
			"status":    {"open"},
		})
		if err != nil {
			return err
		}
		_, openOrderCount, err := output.PaginatedCount(ordersData)
		if err != nil {
			return err
		}

		// Fetch inventory total count
		inventoryData, err := client.Get("/inventory", url.Values{
			"page_size": {"1"},
		})
		if err != nil {
			return err
		}
		_, inventoryCount, err := output.PaginatedCount(inventoryData)
		if err != nil {
			return err
		}

		if jsonOutput {
			return output.JSON(map[string]any{
				"open_orders":     openOrderCount,
				"inventory_items": inventoryCount,
			})
		}

		fmt.Printf("Open orders:     %d\n", openOrderCount)
		fmt.Printf("Inventory items: %d\n", inventoryCount)
		return nil
	},
}
