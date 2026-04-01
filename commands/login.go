package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/StockpilotHQ/stockpilot-cli/internal/api"
	"github.com/StockpilotHQ/stockpilot-cli/internal/config"
	"golang.org/x/term"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Save your Stockpilot API credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Client ID: ")
		clientID, _ := reader.ReadString('\n')
		clientID = strings.TrimSpace(clientID)

		fmt.Print("Client Secret: ")
		secretBytes, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			return err
		}
		clientSecret := strings.TrimSpace(string(secretBytes))

		// Verify credentials before saving
		client := api.New(clientID, clientSecret)
		if _, err := client.Get("/auth/who-is", nil); err != nil {
			return fmt.Errorf("invalid credentials: %w", err)
		}

		cfg := &config.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
		}
		if err := config.Save(cfg); err != nil {
			return err
		}

		fmt.Println("Logged in successfully.")
		return nil
	},
}
