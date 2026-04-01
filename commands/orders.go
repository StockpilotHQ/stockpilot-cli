package commands

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/StockpilotHQ/stockpilot-cli/internal/api"
	"github.com/StockpilotHQ/stockpilot-cli/internal/config"
	"github.com/StockpilotHQ/stockpilot-cli/internal/output"
)

var ordersCmd = &cobra.Command{
	Use:   "orders",
	Short: "Manage orders",
}

var ordersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List orders",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		client := api.New(cfg.ClientID, cfg.ClientSecret)

		page, _ := cmd.Flags().GetInt("page")
		pageSize, _ := cmd.Flags().GetInt("page-size")
		params := url.Values{
			"page":      {fmt.Sprint(page)},
			"page_size": {fmt.Sprint(pageSize)},
		}
		if status, _ := cmd.Flags().GetString("status"); status != "" {
			params.Set("status", status)
		}

		data, err := client.Get("/orders", params)
		if err != nil {
			return err
		}

		orders, err := output.Paginated(data)
		if err != nil {
			return err
		}

		if jsonOutput {
			return output.JSON(orders)
		}

		rows := make([][]string, 0, len(orders))
		for _, o := range orders {
			rows = append(rows, []string{
				fmt.Sprint(o["id"]),
				fmt.Sprint(o["status"]),
				fmt.Sprint(o["customer_name"]),
				fmt.Sprint(o["created_at"]),
			})
		}
		output.Table([]string{"ID", "Status", "Customer", "Created"}, rows)
		return nil
	},
}

var ordersGetCmd = &cobra.Command{
	Use:   "get ORDER_ID",
	Short: "Get a single order",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		client := api.New(cfg.ClientID, cfg.ClientSecret)
		data, err := client.Get("/orders/get-single", url.Values{"order_id": {args[0]}})
		if err != nil {
			return err
		}
		var v any
		json.Unmarshal(data, &v)
		return output.JSON(v)
	},
}

var ordersFulfilCmd = &cobra.Command{
	Use:   "fulfil ORDER_ID",
	Short: "Mark an order as fulfilled",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		client := api.New(cfg.ClientID, cfg.ClientSecret)
		data, err := client.Post("/orders/fulfil", map[string]any{
			"order_id": args[0],
			"items":    []any{},
		})
		if err != nil {
			return err
		}
		var v any
		json.Unmarshal(data, &v)
		return output.JSON(v)
	},
}

var ordersCancelCmd = &cobra.Command{
	Use:   "cancel ORDER_ID",
	Short: "Cancel an order",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		client := api.New(cfg.ClientID, cfg.ClientSecret)
		data, err := client.Put("/orders/cancel-order", map[string]any{
			"order_pk": args[0],
		})
		if err != nil {
			return err
		}
		var v any
		json.Unmarshal(data, &v)
		return output.JSON(v)
	},
}

func init() {
	ordersListCmd.Flags().Int("page", 1, "Page number")
	ordersListCmd.Flags().Int("page-size", 100, "Items per page")
	ordersListCmd.Flags().String("status", "", "Filter by status (e.g. pending, fulfilled, cancelled)")
	ordersCmd.AddCommand(ordersListCmd, ordersGetCmd, ordersFulfilCmd, ordersCancelCmd)
}
