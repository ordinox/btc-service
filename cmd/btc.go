package cmd

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/alexellis/go-execute/v2"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/ordinox/btc-service/btc"
	"github.com/ordinox/btc-service/client"
	"github.com/ordinox/btc-service/common"
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
			fmt.Println("blocks generated")
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
				fmt.Println(utxos.Result)
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
				fmt.Println(u)
			}
			return nil
		},
	}
	return cmd
}
