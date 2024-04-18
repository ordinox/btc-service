package cmd

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/ordinox/btc-service/config"
	"github.com/ordinox/btc-service/runes"
	"github.com/spf13/cobra"
)

func preRunForceArgs(argLength int) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		_ = cmd.MarkFlagRequired("fee-rate")
		if len(args) < argLength {
			_ = cmd.Help()
			os.Exit(1)
		}
	}
}

func forceFeeRateFlag(cmd *cobra.Command) int {
	val, err := cmd.Flags().GetString("fee-rate")
	if err != nil {
		fmt.Println("Error: Fee Rate not set. Use --fee-rate")
		_ = cmd.Help()
		os.Exit(1)
	}
	if val == "" {
		fmt.Println("Error: Fee Rate not set. Use --fee-rate")
		_ = cmd.Help()
		os.Exit(1)
	}
	fee, err := strconv.Atoi(val)
	if err != nil {
		fmt.Println("Error: Fee Rate should be an integer")
		os.Exit(1)
	}
	return fee
}

func parseBtcAddress(addrStr string, config config.Config) btcutil.Address {
	addr, err := btcutil.DecodeAddress(addrStr, config.BtcConfig.GetChainConfigParams())
	if err != nil {
		fmt.Printf("Error: Invalid bitcoin address %s\n", addrStr)
		os.Exit(1)
	}
	return addr
}

func parsePrivateKey(pkStr string) *btcec.PrivateKey {
	privKeyB, err := hex.DecodeString(pkStr)
	if err != nil {
		fmt.Printf("Error: Invalid private key %s\n", pkStr)
		os.Exit(1)
	}
	privKey, _ := btcec.PrivKeyFromBytes(privKeyB)
	return privKey
}

func parseRune(runeStr string) runes.Rune {
	split := strings.Split(runeStr, ":")
	if len(split) != 2 {
		fmt.Printf("Error: Invalid Rune ID %s\n", runeStr)
		os.Exit(1)
	}

	blockNumber, err := strconv.Atoi(split[0])
	if err != nil {
		fmt.Printf("Error: Invalid Rune ID %s\n", runeStr)
		os.Exit(1)
	}

	txIdx, err := strconv.Atoi(split[1])
	if err != nil {
		fmt.Printf("Error: Invalid Rune ID %s\n", runeStr)
		os.Exit(1)
	}
	return runes.Rune{
		BlockNumber: uint64(blockNumber),
		TxIndex:     uint32(txIdx),
	}
}

func parseUint64(str string) uint64 {
	res, err := strconv.Atoi(str)
	if err != nil {
		fmt.Printf("Error: Invalid Number %s\n", str)
		os.Exit(1)
	}
	return uint64(res)
}

func parseTicker(str string) string {
	if len(str) != 4 {
		fmt.Printf("Error: Invalid Ticker %s\n", str)
		os.Exit(1)
	}
	return str
}

func parseString(str string) string {
	if len(str) == 0 {
		fmt.Printf("Error: Invalid String %s\n", str)
		os.Exit(1)
	}
	return str
}

func parseBigInt(str string) *big.Int {
	val, ok := big.NewInt(0).SetString(str, 10)
	if !ok {
		fmt.Printf("Error: Invalid number string%s\n", str)
		os.Exit(1)
	}
	return val
}
