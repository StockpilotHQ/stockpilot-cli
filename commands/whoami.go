package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/StockpilotHQ/stockpilot-cli/internal/api"
	"github.com/StockpilotHQ/stockpilot-cli/internal/config"
	"github.com/StockpilotHQ/stockpilot-cli/internal/output"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Verify credentials and show your organization",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		client := api.New(cfg.ClientID, cfg.ClientSecret)
		data, err := client.Get("/auth/who-is", nil)
		if err != nil {
			return err
		}

		if jsonOutput {
			var v any
			json.Unmarshal(data, &v)
			return output.JSON(v)
		}

		var resp struct {
			ID               int    `json:"id"`
			OrganizationName string `json:"organization_name"`
			UniqueID         string `json:"unique_id"`
		}
		if err := json.Unmarshal(data, &resp); err != nil {
			return err
		}
		output.Table(
			[]string{"ID", "Organization", "Unique ID"},
			[][]string{{
				fmt.Sprint(resp.ID),
				resp.OrganizationName,
				resp.UniqueID,
			}},
		)
		return nil
	},
}
