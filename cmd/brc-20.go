package cmd

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/ordinox/btc-service/brc20"
	"github.com/ordinox/btc-service/config"
	"github.com/spf13/cobra"
)

func transferCmd(config config.Config) *cobra.Command {
	transferCmd := cobra.Command{
		Use:   "transfer [fromAddr] [toAddr] [transferInscriptionId] [amt] [privateKey]",
		Short: "transfer inscriptions",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			fromAddr, err := btcutil.DecodeAddress(args[0], config.BtcConfig.GetChainConfigParams())
			if err != nil {
				return err
			}

			toAddr, err := btcutil.DecodeAddress(args[1], config.BtcConfig.GetChainConfigParams())
			if err != nil {
				return err
			}

			transferInscription := args[2]
			if strings.TrimSpace(transferInscription) == "" {
				return fmt.Errorf("ticker cannot be empty")
			}

			amt, err := strconv.Atoi(args[3])
			if err != nil {
				return err
			}

			privKeyB, err := hex.DecodeString(args[4])
			if err != nil {
				return err
			}
			privKey, _ := btcec.PrivKeyFromBytes(privKeyB)
			res, err := brc20.TransferBrc20(fromAddr, toAddr, transferInscription, uint(amt), *privKey, config)
			if err != nil {
				return err
			}
			fmt.Println(*res)
			return nil
		},
	}
	return &transferCmd
}
