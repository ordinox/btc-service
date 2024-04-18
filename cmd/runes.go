package cmd

import (
	"fmt"
	"os"

	"github.com/markkurossi/tabulate"
	"github.com/ordinox/btc-service/client"
	"github.com/ordinox/btc-service/config"
	"github.com/ordinox/btc-service/runes"
	"github.com/spf13/cobra"
)

func runesCmd(config config.Config) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "runes",
		Short: "interact with brc20 tokens on bitcoin",
	}
	cmd.AddCommand(
		mintRunesCmd(config),
		transferRuneCmd(config),
		runesBalanceCmd(config),
	)
	return
}

func mintRunesCmd(config config.Config) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:    "mint RUNE_ID FROM_ADDR PRIV_KEY_HEX",
		PreRun: preRunForceArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			feeRate := forceFeeRateFlag(cmd)
			rune := parseRune(args[0])
			addr := parseBtcAddress(args[1], config)
			privKey := parsePrivateKey(args[2])
			hash, err := runes.MintRunes(rune, addr, privKey, uint64(feeRate), config)
			if err != nil {
				fmt.Println("error executing mint")
				fmt.Println(err.Error())
				os.Exit(1)
			}
			fmt.Println("runes minted successfully")
			fmt.Println("commit", (*hash).String())
		},
	}

	_ = cmd.MarkFlagRequired("fee-rate")
	_ = cmd.Flags().StringP("fee-rate", "f", "", "Fee rate for submitting transactions")
	return
}

func transferRuneCmd(config config.Config) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:    "transfer RUNE_ID AMT FROM_ADDR TO_ADDR PRIV_KEY_HEX",
		PreRun: preRunForceArgs(5),
		Run: func(cmd *cobra.Command, args []string) {
			feeRate := forceFeeRateFlag(cmd)
			rune := parseRune(args[0])
			amt := parseBigInt(args[1])
			addr := parseBtcAddress(args[2], config)
			toAddr := parseBtcAddress(args[3], config)
			privKey := parsePrivateKey(args[4])
			hash, err := runes.TransferRune(rune, amt, addr, toAddr, privKey, uint64(feeRate), config)
			if err != nil {
				fmt.Println("error executing mint")
				fmt.Println(err.Error())
				os.Exit(1)
			}
			fmt.Println("runes tranferred successfully")
			fmt.Println("commit", (*hash).String())
		},
	}

	_ = cmd.MarkFlagRequired("fee-rate")
	_ = cmd.Flags().StringP("fee-rate", "f", "", "Fee rate for submitting transactions")
	return
}

func runesBalanceCmd(c config.Config) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:    "balance ADDRESS",
		PreRun: preRunForceArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			addr := parseBtcAddress(args[0], c)
			opiClient := client.NewOpiClient(c.OpiConfig)
			bal, err := opiClient.GetRunesBalance(addr.EncodeAddress())
			if err != nil {
				fmt.Println("error connecting to OPI Runes")
				fmt.Println(err)
				os.Exit(1)
			}
			tab := tabulate.New(tabulate.ASCII)
			_ = tabulate.Reflect(tab, 0, nil, bal)
			tab.Print(os.Stdout)
		},
	}
	return
}
