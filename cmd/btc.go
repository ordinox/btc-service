package cmd

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
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
