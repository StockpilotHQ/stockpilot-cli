package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	jsonOutput bool
	version    = "dev"
)

var rootCmd = &cobra.Command{
	Use:   "stockpilot",
	Short: "Stockpilot CLI — manage inventory, orders, and fulfillment from your terminal",
	Long: `Stockpilot CLI gives you full access to Stockpilot from your terminal.
Works standalone or with AI agents (Claude, Cursor, Codex, etc.).

Get started:
  stockpilot login
  stockpilot whoami
  stockpilot status`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output as JSON (for scripts and AI agents)")
	rootCmd.Version = version

	rootCmd.AddCommand(
		loginCmd,
		whoamiCmd,
		statusCmd,
		inventoryCmd,
		ordersCmd,
		productsCmd,
		customersCmd,
		analyticsCmd,
		shippingCmd,
	)
}
