package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/StockpilotHQ/stockpilot-cli/internal/api"
	"github.com/StockpilotHQ/stockpilot-cli/internal/config"
	"github.com/StockpilotHQ/stockpilot-cli/internal/output"
)

var shippingCmd = &cobra.Command{
	Use:   "shipping",
	Short: "Shipping and label management",
}

var shippingLabelCmd = &cobra.Command{
	Use:   "label ORDER_ID",
	Short: "Request a shipping label for an order",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		client := api.New(cfg.ClientID, cfg.ClientSecret)

		integrationID, _ := cmd.Flags().GetString("integration")
		if integrationID == "" {
			return fmt.Errorf("--integration is required")
		}

		data, err := client.Post("/shipping/request-label", map[string]any{
			"order_id":       args[0],
			"integration_id": integrationID,
		})
		if err != nil {
			return err
		}
		var v any
		json.Unmarshal(data, &v)
		return output.JSON(v)
	},
}

var shippingIntegrationsCmd = &cobra.Command{
	Use:   "integrations",
	Short: "List available shipping integrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		client := api.New(cfg.ClientID, cfg.ClientSecret)
		data, err := client.Get("/shipping/integrations", nil)
		if err != nil {
			return err
		}
		var v any
		json.Unmarshal(data, &v)
		return output.JSON(v)
	},
}

func init() {
	shippingLabelCmd.Flags().String("integration", "", "Shipping integration ID (required)")
	shippingCmd.AddCommand(shippingLabelCmd, shippingIntegrationsCmd)
}
