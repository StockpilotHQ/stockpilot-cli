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

var inventoryCmd = &cobra.Command{
	Use:   "inventory",
	Short: "Manage inventory",
}

var inventoryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all inventory items",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		client := api.New(cfg.ClientID, cfg.ClientSecret)

		page, _ := cmd.Flags().GetInt("page")
		pageSize, _ := cmd.Flags().GetInt("page-size")

		data, err := client.Get("/inventory", url.Values{
			"page":      {fmt.Sprint(page)},
			"page_size": {fmt.Sprint(pageSize)},
		})
		if err != nil {
			return err
		}

		if jsonOutput {
			var v any
			json.Unmarshal(data, &v)
			return output.JSON(v)
		}

		var items []map[string]any
		if err := json.Unmarshal(data, &items); err != nil {
			return err
		}
		rows := make([][]string, 0, len(items))
		for _, item := range items {
			rows = append(rows, []string{
				fmt.Sprint(item["id"]),
				fmt.Sprint(item["sku"]),
				fmt.Sprint(item["name"]),
				fmt.Sprint(item["quantity"]),
				fmt.Sprint(item["location"]),
			})
		}
		output.Table([]string{"ID", "SKU", "Name", "Qty", "Location"}, rows)
		return nil
	},
}

var inventoryGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a single inventory item by SKU, ID, or barcode",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		client := api.New(cfg.ClientID, cfg.ClientSecret)

		params := url.Values{}
		if v, _ := cmd.Flags().GetString("sku"); v != "" {
			params.Set("sku", v)
		} else if v, _ := cmd.Flags().GetString("id"); v != "" {
			params.Set("id", v)
		} else if v, _ := cmd.Flags().GetString("barcode"); v != "" {
			params.Set("barcode", v)
		} else {
			return fmt.Errorf("one of --sku, --id, or --barcode is required")
		}

		data, err := client.Get("/inventory/get", params)
		if err != nil {
			return err
		}

		var v any
		json.Unmarshal(data, &v)
		return output.JSON(v)
	},
}

var inventoryUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update inventory item quantity, location, or threshold",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		client := api.New(cfg.ClientID, cfg.ClientSecret)

		body := map[string]any{}
		if v, _ := cmd.Flags().GetString("sku"); v != "" {
			body["sku"] = v
		} else if v, _ := cmd.Flags().GetString("id"); v != "" {
			body["item_id"] = v
		} else if v, _ := cmd.Flags().GetString("barcode"); v != "" {
			body["barcode"] = v
		} else {
			return fmt.Errorf("one of --sku, --id, or --barcode is required")
		}
		if v, _ := cmd.Flags().GetInt("quantity"); cmd.Flags().Changed("quantity") {
			body["quantity"] = v
		}
		if v, _ := cmd.Flags().GetString("location"); v != "" {
			body["location"] = v
		}
		if v, _ := cmd.Flags().GetString("threshold"); v != "" {
			body["threshold"] = v
		}

		data, err := client.Post("/inventory/update", body)
		if err != nil {
			return err
		}

		var v any
		json.Unmarshal(data, &v)
		return output.JSON(v)
	},
}

func init() {
	inventoryListCmd.Flags().Int("page", 1, "Page number")
	inventoryListCmd.Flags().Int("page-size", 100, "Items per page")

	inventoryGetCmd.Flags().String("sku", "", "SKU")
	inventoryGetCmd.Flags().String("id", "", "Item ID")
	inventoryGetCmd.Flags().String("barcode", "", "Barcode / EAN")

	inventoryUpdateCmd.Flags().String("sku", "", "SKU")
	inventoryUpdateCmd.Flags().String("id", "", "Item ID")
	inventoryUpdateCmd.Flags().String("barcode", "", "Barcode / EAN")
	inventoryUpdateCmd.Flags().Int("quantity", 0, "New quantity")
	inventoryUpdateCmd.Flags().String("location", "", "Bin location (e.g. A1-001-01)")
	inventoryUpdateCmd.Flags().String("threshold", "", "Stock threshold (e.g. 5u or 33w)")

	inventoryCmd.AddCommand(inventoryListCmd, inventoryGetCmd, inventoryUpdateCmd)
}
