package runes_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ordinox/btc-service/client"
	"github.com/ordinox/btc-service/cmd"
	"github.com/ordinox/btc-service/common"
	"github.com/ordinox/btc-service/config"
	"github.com/ordinox/btc-service/runes"
	"github.com/stretchr/testify/require"
)

const TEST_PRIVATE_KEY = "37b4de3cdd51c013addf9c727be371b3a2258d457552798a76778712bacd01f2"
const TEST_RUNE_ID = "310:1"

func testTransfer(t *testing.T, amt int64) {
	config.Init()
	// Init
	var (
		runeBalance *client.RunesBalance

		config      = config.GetDefaultConfig()
		opiClient   = client.NewOpiClient(config.OpiConfig)
		btcClient   = client.NewBitcoinClient(config)
		pk          = common.LoadPrivateKey(TEST_PRIVATE_KEY)
		addr, err   = common.GetP2PKHAddress(pk.PubKey().SerializeCompressed(), &chaincfg.RegressionNetParams)
		destAddr, _ = btcutil.DecodeAddress("msTNoMYXynujefraNmpeKKuUsAsSYMZ8C2", &chaincfg.RegressionNetParams)
	)
	require.Nil(t, err)

	fmt.Printf("Transferring %d runes\n", amt)
	// Check Balance - Store balance
	runesBalances, err := opiClient.GetRunesBalance(addr.EncodeAddress())
	require.Nil(t, err)
	for _, i := range runesBalances {
		if i.RuneID == TEST_RUNE_ID {
			runeBalance = i.Copy()
			break
		}
	}

	require.NotNil(t, runeBalance)
	fmt.Println("balance before transfer", runeBalance.TotalBalance)

	// Do Transfer
	rune, err := runes.ParseRune(runeBalance.RuneID)
	require.Nil(t, err)

	hash, err := runes.TransferRune(*rune, big.NewInt(amt), addr, destAddr, pk, 22, config)
	require.Nil(t, err)

	// Check balance
	err = cmd.GenerateBlocks()
	h, _ := btcClient.GetBlockCount()
	fmt.Println("BlockHeight", h)
	require.Nil(t, err)
	// Wait for the indexer to catch up

	time.Sleep(5 * time.Second)

	runesBalances, err = opiClient.GetRunesBalance(addr.EncodeAddress())
	require.Nil(t, err)
	for _, i := range runesBalances {
		if i.RuneID == TEST_RUNE_ID {
			runeBalance = i.Copy()
			break
		}
	}
	fmt.Println("balance after transfer", runeBalance.TotalBalance)
	fmt.Println((*hash).String())
	fmt.Println("------")
}

func TestRuneTransfer(t *testing.T) {
	testTransfer(t, 1)
	testTransfer(t, 2)
	testTransfer(t, 3)
	testTransfer(t, 4)
	testTransfer(t, 5)
}
