package cmd

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/ordinox/btc-service/brc20"
	"github.com/ordinox/btc-service/config"
	"github.com/spf13/cobra"
)

func e2eCmd(config config.Config) *cobra.Command {
	e2eCmd := cobra.Command{
		Use:   "e2e [tokenName] [amt] [fromAddr] [toAddr] [privateKey] [fee-rate]",
		Short: "mint and transfer in one command [ONLY FOR REGTEST]",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			fromAddr, err := btcutil.DecodeAddress(args[2], config.BtcConfig.GetChainConfigParams())
			if err != nil {
				return err
			}
			toAddr, err := btcutil.DecodeAddress(args[3], config.BtcConfig.GetChainConfigParams())
			if err != nil {
				return err
			}
			amt, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			privKeyB, err := hex.DecodeString(args[4])
			if err != nil {
				return err
			}

			privKey, _ := btcec.PrivKeyFromBytes(privKeyB)

			feeRate, err := strconv.Atoi(args[5])
			if err != nil {
				return err
			}

			tokenName := args[0]

			fmt.Println("inscribing mint inscription...")

			_, err = brc20.InscribeMint(tokenName, uint(amt), args[2], uint(feeRate), config.BtcConfig)
			if err != nil {
				return err
			}

			if err := generateBlocks(); err != nil {
				return err
			}

			fmt.Println("done")
			time.Sleep(1 * time.Second)
			fmt.Println("inscribing transfer inscription...")

			insc, err := brc20.InscribeTransfer(tokenName, fromAddr, uint(amt), uint(feeRate), config)
			if err != nil {
				return err
			}

			if err := generateBlocks(); err != nil {
				return err
			}
			fmt.Println("done")
			time.Sleep(1 * time.Second)

			transferInscription := insc.Inscriptions[0].Id

			fmt.Println("transferring inscription...")

			res, err := brc20.TransferBrc20(fromAddr, toAddr, transferInscription, uint(amt), *privKey, uint64(feeRate), config)
			if err != nil {
				return err
			}

			if err := generateBlocks(); err != nil {
				return err
			}

			fmt.Println("done")
			fmt.Println("inscription ID transferred: ", transferInscription)
			fmt.Println("commit hash", *res)

			return nil
		},
	}
	return &e2eCmd
}

func transferCmd(config config.Config) *cobra.Command {
	transferCmd := cobra.Command{
		Use:   "transfer [fromAddr] [toAddr] [transferInscriptionId] [amt] [privateKey] [fee-rate]",
		Short: "transfer inscriptions",
		Args:  cobra.ExactArgs(6),
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

			feeRate, err := strconv.Atoi(args[5])
			if err != nil {
				return err
			}
			privKey, _ := btcec.PrivKeyFromBytes(privKeyB)
			res, err := brc20.TransferBrc20(fromAddr, toAddr, transferInscription, uint(amt), *privKey, uint64(feeRate), config)
			if err != nil {
				return err
			}
			fmt.Println(*res)
			return nil
		},
	}
	return &transferCmd
}
