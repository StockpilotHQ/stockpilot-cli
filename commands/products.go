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

var productsCmd = &cobra.Command{
	Use:   "products",
	Short: "Manage products",
}

var productsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List products",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		client := api.New(cfg.ClientID, cfg.ClientSecret)

		page, _ := cmd.Flags().GetInt("page")
		pageSize, _ := cmd.Flags().GetInt("page-size")

		data, err := client.Get("/products", url.Values{
			"page":      {fmt.Sprint(page)},
			"page_size": {fmt.Sprint(pageSize)},
		})
		if err != nil {
			return err
		}

		products, err := output.Paginated(data)
		if err != nil {
			return err
		}

		if jsonOutput {
			return output.JSON(products)
		}

		rows := make([][]string, 0, len(products))
		for _, p := range products {
			rows = append(rows, []string{
				fmt.Sprint(p["id"]),
				fmt.Sprint(p["sku"]),
				fmt.Sprint(p["name"]),
				fmt.Sprint(p["retail_price"]),
			})
		}
		output.Table([]string{"ID", "SKU", "Name", "Price"}, rows)
		return nil
	},
}

var productsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a product by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		client := api.New(cfg.ClientID, cfg.ClientSecret)
		id, _ := cmd.Flags().GetString("id")
		if id == "" {
			return fmt.Errorf("--id is required")
		}
		data, err := client.Get("/products/get", url.Values{"id": {id}})
		if err != nil {
			return err
		}
		var v any
		json.Unmarshal(data, &v)
		return output.JSON(v)
	},
}

func init() {
	productsListCmd.Flags().Int("page", 1, "Page number")
	productsListCmd.Flags().Int("page-size", 100, "Items per page")
	productsGetCmd.Flags().String("id", "", "Product ID")
	productsCmd.AddCommand(productsListCmd, productsGetCmd)
}
