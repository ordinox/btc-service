package cmd

import (
	"fmt"
	"os"

	"github.com/ordinox/btc-service/brc20"
	"github.com/ordinox/btc-service/config"
	"github.com/spf13/cobra"
)

func e2eCmd(config config.Config) *cobra.Command {
	e2eCmd := cobra.Command{
		Use:    "e2e TOKEN AMT FROM_ADDRESS TO_ADDRESS INSCRIBER_PRIVATE_KEY SENDER_PRIVATE_KEY",
		Short:  "mint and transfer in one command [ONLY FOR REGTEST]",
		PreRun: preRunForceArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			genBlocks := func() {
				// you don't need block generation as long as txns are in the mempool
				fmt.Println("no blocks generated")
				// if err := generateBlocks(); err != nil {
				// 	fmt.Println("Error generating blocks")
				// 	fmt.Println(err.Error())
				// 	os.Exit(1)
				// }
			}

			feeRate := forceFeeRateFlag(cmd)
			ticker := parseTicker(args[0])
			amt := parseUint64(args[1])
			fromAddr := parseBtcAddress(args[2], config)
			toAddr := parseBtcAddress(args[3], config)
			inscriberPrivateKey := parsePrivateKey(args[4])
			senderPrivateKey := parsePrivateKey(args[5])

			_, err := brc20.InscribeMint(ticker, amt, fromAddr, inscriberPrivateKey, uint64(feeRate), config)
			if err != nil {
				fmt.Println("Error occured while minting")
				fmt.Println(err.Error())
				os.Exit(1)
			}

			genBlocks()
			fmt.Println("inscribing transfer inscription...")

			insc, err := brc20.InscribeTransfer(ticker, amt, fromAddr, inscriberPrivateKey, uint64(feeRate), config)
			if err != nil {
				fmt.Println("Error occured while inscribing transfer")
				fmt.Println(err.Error())
				os.Exit(1)
			}

			genBlocks()

			transferInscription := insc.RevealTx

			fmt.Println("transferring inscription...")

			res, err := brc20.TransferBrc20(fromAddr, toAddr, transferInscription, senderPrivateKey, uint64(feeRate), config)
			if err != nil {
				fmt.Println("Error occured while transferring")
				fmt.Println(err.Error())
				os.Exit(1)
			}

			genBlocks()

			fmt.Println("done")
			fmt.Println("inscription ID transferred: ", transferInscription)
			fmt.Println("commit hash", *res)

			return nil
		},
	}
	_ = e2eCmd.MarkFlagRequired("fee-rate")
	_ = e2eCmd.Flags().StringP("fee-rate", "f", "", "Fee rate for submitting transactions")
	return &e2eCmd
}

func transferCmd(config config.Config) *cobra.Command {
	transferCmd := cobra.Command{
		Use:   "transfer FROM_ADDR TO_ADDR TRANSFER_INSCRIPTION SENDER_PRIVATE_KEY",
		Short: "transfer inscriptions",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			feeRate := forceFeeRateFlag(cmd)
			fromAddr := parseBtcAddress(args[0], config)
			toAddr := parseBtcAddress(args[1], config)
			transferInscription := parseString(args[2])
			privateKey := parsePrivateKey(args[3])
			hashPtr, err := brc20.TransferBrc20(fromAddr, toAddr, transferInscription, privateKey, uint64(feeRate), config)
			if err != nil {
				fmt.Println("Error occured while transferring")
				fmt.Println(err.Error())
				os.Exit(1)
			}
			fmt.Println("Transfer Initiated:", *hashPtr)
			return nil
		},
	}
	_ = transferCmd.MarkFlagRequired("fee-rate")
	_ = transferCmd.Flags().StringP("fee-rate", "f", "", "Fee rate for submitting transactions")
	return &transferCmd
}

func sendBrc20Cmd(config config.Config) *cobra.Command {
	transferCmd := cobra.Command{
		Use:   "send TOKEN AMT FROM_ADDRESS TO_ADDRESS INSCRIBER_PRIVATE_KEY SENDER_PRIVATE_KEY",
		Short: "inscribe transfer + transfer inscription",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			feeRate := forceFeeRateFlag(cmd)
			ticker := parseTicker(args[0])
			amt := parseUint64(args[1])
			fromAddr := parseBtcAddress(args[2], config)
			toAddr := parseBtcAddress(args[3], config)
			inscriberPrivateKey := parsePrivateKey(args[4])
			senderPrivateKey := parsePrivateKey(args[5])

			inscriptionId, hash, err := brc20.SendBrc20(ticker, fromAddr, toAddr, uint64(amt), uint64(feeRate), inscriberPrivateKey, senderPrivateKey, config)
			if err != nil {
				fmt.Println("Error occured while transferring")
				fmt.Println(err.Error())
				os.Exit(1)
			}
			fmt.Println("----")
			fmt.Println("Inscription ID", inscriptionId)
			fmt.Println("Tx Hash", hash)
			fmt.Println("----")
			return nil
		},
	}
	_ = transferCmd.MarkFlagRequired("fee-rate")
	_ = transferCmd.Flags().StringP("fee-rate", "f", "", "Fee rate for submitting transactions")
	return &transferCmd
}
