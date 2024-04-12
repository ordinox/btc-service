package cmd

import (
	"fmt"
	"os"

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
		brc20InscribeCmd(config),
		transferCmd(config),
		e2eCmd(config),
		sendBrc20Cmd(config),
	)
	return &brc20Cmd
}

func brc20InscribeCmd(config config.Config) *cobra.Command {
	brc20InscribeCmd := cobra.Command{
		Use:   "inscribe",
		Short: "inscribe brc20 inscriptions",
	}
	_ = brc20InscribeCmd.MarkFlagRequired("fee-rate")
	_ = brc20InscribeCmd.Flags().StringP("fee-rate", "f", "", "Fee rate for submitting transactions")
	_ = brc20InscribeCmd.PersistentFlags().StringP("fee-rate", "f", "", "Fee rate for submitting transactions")
	brc20InscribeCmd.AddCommand(
		inscribeDeployCmd(config),
		inscribeMintCmd(config),
		inscribeTransferCmd(config),
	)
	return &brc20InscribeCmd
}

func getBalance(config config.OpiConfig) *cobra.Command {
	getBalanceCmd := cobra.Command{
		Use:   "balance  [ticker] [address]",
		Short: "get brc20 balance for a token",
		Args:  cobra.ExactArgs(2),
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

func inscribeDeployCmd(config config.Config) *cobra.Command {
	deployCmd := cobra.Command{
		Use:    "deploy TICKER SUPPLY DESTINATION_ADDR SENDER_PRIVATE_KEY",
		Short:  "deploy a brc20 token",
		PreRun: preRunForceArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			feeRate := forceFeeRateFlag(cmd)
			ticker := parseTicker(args[0])
			supply := parseUint64(args[1])
			addr := parseBtcAddress(args[2], config)
			privateKey := parsePrivateKey(args[3])

			insc, err := brc20.InscribeDeploy(ticker, uint(supply), privateKey, addr, uint64(feeRate), config)
			if err != nil {
				fmt.Println("Error occured while deploying")
				fmt.Println(err.Error())
				os.Exit(1)
			}
			fmt.Println("CommitTx:", insc.CommitTx)
			fmt.Println("RevelTx:", insc.RevealTx)
			fmt.Println("Fee:", insc.TotalFeePaid)
			return nil
		},
	}
	return &deployCmd
}

func inscribeMintCmd(config config.Config) *cobra.Command {
	mintCmd := cobra.Command{
		Use:    "mint TICKER AMT DESTINATION_ADDR SENDER_PRIVATE_KEY",
		Short:  "mint a brc20 token",
		PreRun: preRunForceArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			feeRate := forceFeeRateFlag(cmd)
			ticker := parseTicker(args[0])
			amt := parseUint64(args[1])
			addr := parseBtcAddress(args[2], config)
			privateKey := parsePrivateKey(args[3])
			insc, err := brc20.InscribeMint(ticker, amt, addr, privateKey, uint64(feeRate), config)
			if err != nil {
				fmt.Println("Error occured while minting")
				fmt.Println(err.Error())
				os.Exit(1)
			}
			fmt.Println("CommitTx:", insc.CommitTx)
			fmt.Println("RevelTx:", insc.RevealTx)
			fmt.Println("Fee:", insc.TotalFeePaid)
			return nil
		},
	}
	return &mintCmd
}

func inscribeTransferCmd(config config.Config) *cobra.Command {
	transferCmd := cobra.Command{
		Use:    "transfer TICKER AMT DESTINATION_ADDR SENDER_PRIVATE_KEY",
		Args:   cobra.MatchAll(cobra.ExactArgs(4)),
		Short:  "transfer brc20 tokens from the given address to another",
		PreRun: preRunForceArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			feeRate := forceFeeRateFlag(cmd)
			ticker := parseTicker(args[0])
			amt := parseUint64(args[1])
			addr := parseBtcAddress(args[2], config)
			privateKey := parsePrivateKey(args[3])

			insc, err := brc20.InscribeTransfer(ticker, amt, addr, privateKey, uint64(feeRate), config)
			if err != nil {
				fmt.Println("Error occured while inscribing transfer")
				fmt.Println(err.Error())
				os.Exit(1)
			}
			fmt.Println("CommitTx:", insc.CommitTx)
			fmt.Println("RevelTx:", insc.RevealTx)
			fmt.Println("Fee:", insc.TotalFeePaid)
			return nil
		},
	}
	return &transferCmd
}

func printInscriptionRes(res *inscriptions.InscriptionResultRaw) {
	fmt.Println("commit: ", res.Commit)
	fmt.Println("inscriptionId: ", res.Inscriptions[0].Id)
}
