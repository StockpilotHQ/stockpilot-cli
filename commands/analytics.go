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

var analyticsCmd = &cobra.Command{
	Use:   "analytics",
	Short: "Sales analytics",
}

var analyticsSalesCmd = &cobra.Command{
	Use:   "sales",
	Short: "Sales analytics for an inventory item",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		client := api.New(cfg.ClientID, cfg.ClientSecret)

		itemID, _ := cmd.Flags().GetString("id")
		if itemID == "" {
			return fmt.Errorf("--id is required")
		}

		params := url.Values{"item_id": {itemID}}
		if v, _ := cmd.Flags().GetString("from"); v != "" {
			params.Set("date_from", v)
		}
		if v, _ := cmd.Flags().GetString("to"); v != "" {
			params.Set("date_to", v)
		}
		if v, _ := cmd.Flags().GetString("interval"); v != "" {
			params.Set("interval", v)
		}

		data, err := client.Get("/analytics/items/sales", params)
		if err != nil {
			return err
		}
		var v any
		json.Unmarshal(data, &v)
		return output.JSON(v)
	},
}

var analyticsSummaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Overall sales summary",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		client := api.New(cfg.ClientID, cfg.ClientSecret)

		params := url.Values{}
		if v, _ := cmd.Flags().GetString("from"); v != "" {
			params.Set("date_from", v)
		}
		if v, _ := cmd.Flags().GetString("to"); v != "" {
			params.Set("date_to", v)
		}

		data, err := client.Get("/analytics/sales-summary", params)
		if err != nil {
			return err
		}
		var v any
		json.Unmarshal(data, &v)
		return output.JSON(v)
	},
}

func init() {
	analyticsSalesCmd.Flags().String("id", "", "Inventory item ID (required)")
	analyticsSalesCmd.Flags().String("from", "", "Start date (YYYY-MM-DD)")
	analyticsSalesCmd.Flags().String("to", "", "End date (YYYY-MM-DD)")
	analyticsSalesCmd.Flags().String("interval", "daily", "Interval: daily, weekly, monthly")

	analyticsSummaryCmd.Flags().String("from", "", "Start date (YYYY-MM-DD)")
	analyticsSummaryCmd.Flags().String("to", "", "End date (YYYY-MM-DD)")

	analyticsCmd.AddCommand(analyticsSalesCmd, analyticsSummaryCmd)
}
