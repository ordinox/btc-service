package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/alexellis/go-execute/v2"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/ordinox/btc-service/btc"
	"github.com/ordinox/btc-service/client"
	"github.com/ordinox/btc-service/config"
	"github.com/spf13/cobra"
)

func getKeyPair(config config.BtcConfig) (string, string) {
	pk, _ := btcec.NewPrivateKey()
	pkHex := hex.EncodeToString(pk.Serialize())
	pubKey := pk.PubKey()
	addr, _ := btcutil.NewAddressPubKeyHash(btcutil.Hash160(pubKey.SerializeCompressed()), config.GetChainConfigParams())
	addrStr := addr.EncodeAddress()
	return addrStr, pkHex
}

func getKeyPairCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "keypair",
		Short: "get a btc address & privatekey hex",
		Run: func(cmd *cobra.Command, args []string) {
			addr, privKey := getKeyPair(config.GetDefaultConfig().BtcConfig)
			fmt.Println()
			fmt.Println("Address: ", addr)
			fmt.Println()
			fmt.Println("PrivKeyHex: ", privKey)
			fmt.Println()
		},
	}
	return &cmd
}

func genBlocksCmd(config config.BtcConfig) *cobra.Command {
	cmd := cobra.Command{
		Use:   "genblocks [amt] [address]",
		Short: "generate regtest blocks",
		RunE: func(cmd *cobra.Command, args []string) error {
			amt := 10
			addr := "n3uNm2T4TisRQd3TUmYSrENtbPWEhzqhC2"
			if args[0] != "" {
				amtP, err := strconv.Atoi(args[0])
				if err != nil {
					return err
				}
				amt = amtP
			}
			if args[1] != "" {
				addrP, err := btcutil.DecodeAddress(args[1], config.GetChainConfigParams())
				if err != nil {
					return err
				}
				addr = addrP.EncodeAddress()
			}
			c := execute.ExecTask{
				Command: "bitcoin-cli",
				Args:    []string{"-regtest", "generatetoaddress", fmt.Sprintf("%d", amt), addr},
			}
			_, err := c.Execute(context.Background())
			if err != nil {
				return err
			}
			fmt.Printf("%d blocks generated to %s\n", amt, addr)
			return nil
		},
	}
	return &cmd
}

func generateBlocks() error {
	amt := 10
	addr := "n3uNm2T4TisRQd3TUmYSrENtbPWEhzqhC2"
	c := execute.ExecTask{
		Command: "bitcoin-cli",
		Args:    []string{"-regtest", "generatetoaddress", fmt.Sprintf("%d", amt), addr},
	}
	_, err := c.Execute(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func getUtxosCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "utxos",
		Short: "get utxos for a legacy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if config.GetDefaultConfig().BtcConfig.ChainConfig == "mainnet" {
				// Get utxos using electrum
				utxos, err := btc.GetUtxos(args[0])
				if err != nil {
					return err
				}
				for _, u := range utxos.Result {
					fmt.Println("hash: ", u.TxHash)
					fmt.Println("pos: ", u.Vout)
					fmt.Println("val: ", u.Value)
					fmt.Println("---------------")
				}
				return nil
			}
			client := client.NewBitcoinClient(config.GetDefaultConfig())
			addr, err := btcutil.DecodeAddress(args[0], config.GetDefaultConfig().BtcConfig.GetChainConfigParams())
			if err != nil {
				return err
			}
			utxos, err := client.GetUtxos(addr)
			if err != nil {
				return err
			}
			for _, u := range utxos {
				utxo := u
				fmt.Println("hash: ", utxo.GetTxID())
				fmt.Println("pos: ", utxo.GetVout())
				fmt.Println("val: ", utxo.GetValueInSats())
				fmt.Println("---------------")
			}
			return nil
		},
	}
	return cmd
}

func transferBtcCmd(config config.Config) *cobra.Command {
	transferCmd := cobra.Command{
		Use:   "transfer [fromAddr] [toAddr] [feeRate] [amt] [privateKey]",
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

			feeRate, err := strconv.Atoi(args[2])
			if err != nil {
				return err
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

			err = btc.TransferBtc(
				*privKey,
				fromAddr,
				toAddr,
				nil,
				uint64(amt),
				uint32(feeRate),
			)
			if err != nil {
				return err
			}
			return nil
		},
	}
	return &transferCmd
}
