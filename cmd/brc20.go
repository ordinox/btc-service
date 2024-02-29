package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ordinox/btc-service/brc20"
	"github.com/ordinox/btc-service/config"
	"github.com/spf13/cobra"
)

func brc20Cmd(config config.BtcConfig) *cobra.Command {
	brc20Cmd := cobra.Command{
		Use:   "brc20",
		Short: "interact with brc20 tokens on bitcoin",
	}
	brc20Cmd.AddCommand(deployCmd(config), mintCmd(config))
	return &brc20Cmd
}

func getBalance(config config.BtcConfig) *cobra.Command {
	return nil
}

func deployCmd(config config.BtcConfig) *cobra.Command {
	deployCmd := cobra.Command{
		Use:   "deploy",
		Short: "deploy a brc20 token",
		Long: strings.TrimSpace(`
Example: 
brc20 deploy <TICKER> <SUPPLY> <DESTINATION_ADDR>
		`),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] == "" {
				return fmt.Errorf("ticker name cannot be empty")
			}
			if args[1] == "" {
				return fmt.Errorf("token supply/limit cannot be empty")
			}
			if args[2] == "" {
				return fmt.Errorf("destination cannot be empty")
			}
			amt, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			insc, err := brc20.DeployToken(args[0], uint(amt), args[2], config)
			if err != nil {
				return err
			}
			res, _ := json.MarshalIndent(insc, "", "  ")
			fmt.Print(string(res))
			return nil
		},
	}
	return &deployCmd
}

func mintCmd(config config.BtcConfig) *cobra.Command {
	mintCmd := cobra.Command{
		Use:   "mint",
		Short: "mint a brc20 token",
		Long: strings.TrimSpace(`
Example: 
brc20 mint <TICKER> <AMOUNT> <DESTINATION_ADDR>
		`),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] == "" {
				return fmt.Errorf("ticker name cannot be empty")
			}
			if args[1] == "" {
				return fmt.Errorf("mint amount cannot be empty")
			}
			if args[2] == "" {
				return fmt.Errorf("destination cannot be empty")
			}
			amt, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			insc, err := brc20.MintToken(args[0], uint(amt), args[2], config)
			if err != nil {
				return err
			}
			res, _ := json.MarshalIndent(insc, "", "  ")
			fmt.Print(string(res))
			return nil
		},
	}
	return &mintCmd
}
