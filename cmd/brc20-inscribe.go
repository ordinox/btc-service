package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/ordinox/btc-service/brc20"
	"github.com/ordinox/btc-service/client"
	"github.com/ordinox/btc-service/config"
	"github.com/ordinox/btc-service/inscriptions"
	"github.com/spf13/cobra"
)

func brc20Cmd(config config.Config) *cobra.Command {
	brc20Cmd := cobra.Command{
		Use:   "brc20",
		Short: "interact with brc20 tokens on bitcoin",
	}
	brc20Cmd.AddCommand(
		getBalance(config.OpiConfig),
		inscribeDeployCmd(config.BtcConfig),
		inscribeMintCmd(config.BtcConfig),
		inscribeTransferCmd(config),
		transferCmd(config),
		e2eCmd(config),
	)
	return &brc20Cmd
}

func getBalance(config config.OpiConfig) *cobra.Command {
	getBalanceCmd := cobra.Command{
		Use:   "balance  [ticker] [address]",
		Short: "get brc20 balance for a token",
		Long: strings.TrimSpace(`
Example:
brc20 balance <TICKER> <ADDRESS>
		`),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] == "" {
				return fmt.Errorf("ticker name cannot be empty")
			}
			if args[1] == "" {
				return fmt.Errorf("address cannot be empty")
			}
			balance, err := client.NewOpiClient(config).GetBalance(args[1], args[0])
			if err != nil {
				return err
			}
			fmt.Println("Overall Balance: ", balance.OverallBalance)
			fmt.Println("Available Balance: ", balance.AvailableBalance)
			fmt.Println("Block Height", balance.BlockHeight)
			return nil
		},
	}
	return &getBalanceCmd
}

func inscribeDeployCmd(config config.BtcConfig) *cobra.Command {
	deployCmd := cobra.Command{
		Use:   "inscribe-deploy [ticker] [supply] [destination]",
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
			insc, err := brc20.InscribeDeploy(args[0], uint(amt), args[2], config)
			if err != nil {
				return err
			}
			printInscriptionRes(insc)
			return nil
		},
	}
	return &deployCmd
}

func inscribeMintCmd(config config.BtcConfig) *cobra.Command {
	mintCmd := cobra.Command{
		Use:   "inscribe-mint [ticker] [amount] [destination]",
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
			insc, err := brc20.InscribeMint(args[0], uint(amt), args[2], config)
			if err != nil {
				return err
			}
			printInscriptionRes(insc)
			return nil
		},
	}
	return &mintCmd
}

func inscribeTransferCmd(config config.Config) *cobra.Command {
	transferCmd := cobra.Command{
		Use:        "inscribe-transfer [fromAddr] [ticker] [amt]",
		Args:       cobra.MatchAll(cobra.ExactArgs(3)),
		ArgAliases: []string{"fromAddr", "ticker", "amt"},
		ValidArgs:  []string{"fromAddr", "ticker", "amt"},
		Short:      "transfer brc20 tokens from the given address to another",
		RunE: func(cmd *cobra.Command, args []string) error {
			fromAddr, err := btcutil.DecodeAddress(args[0], config.BtcConfig.GetChainConfigParams())
			if err != nil {
				return err
			}
			ticker := args[1]
			if strings.TrimSpace(ticker) == "" {
				return fmt.Errorf("ticker cannot be empty")
			}

			amt, err := strconv.Atoi(args[2])
			if err != nil {
				return err
			}

			inscriptionsRes, err := brc20.InscribeTransfer(ticker, fromAddr, uint(amt), config)
			if err != nil {
				return err
			}
			printInscriptionRes(inscriptionsRes)

			return nil
		},
	}
	return &transferCmd
}

func printInscriptionRes(res *inscriptions.InscriptionResultRaw) {
	fmt.Println("commit: ", res.Commit)
	fmt.Println("inscriptionId: ", res.Inscriptions[0].Id)
}
