package cmd

import (
	"fmt"
	"os"

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
			fmt.Println((*hash).String())
		},
	}

	return
}
