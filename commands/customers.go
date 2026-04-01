package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/StockpilotHQ/stockpilot-cli/internal/api"
	"github.com/StockpilotHQ/stockpilot-cli/internal/config"
	"github.com/StockpilotHQ/stockpilot-cli/internal/output"
)

var customersCmd = &cobra.Command{
	Use:   "customers",
	Short: "List customers",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		client := api.New(cfg.ClientID, cfg.ClientSecret)
		data, err := client.Get("/customers", nil)
		if err != nil {
			return err
		}

		if jsonOutput {
			var v any
			json.Unmarshal(data, &v)
			return output.JSON(v)
		}

		customers, err := output.Paginated(data)
		if err != nil {
			return err
		}

		rows := make([][]string, 0, len(customers))
		for _, c := range customers {
			rows = append(rows, []string{
				fmt.Sprint(c["id"]),
				fmt.Sprint(c["customer_code"]),
				fmt.Sprint(c["business_name"]),
				fmt.Sprint(c["invoice_email"]),
			})
		}
		output.Table([]string{"ID", "Code", "Business", "Email"}, rows)
		return nil
	},
}
