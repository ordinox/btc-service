package runes_test

import (
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ordinox/btc-service/btc"
	"github.com/ordinox/btc-service/client"
	"github.com/ordinox/btc-service/cmd"
	"github.com/ordinox/btc-service/common"
	"github.com/ordinox/btc-service/config"
	"github.com/ordinox/btc-service/runes"
	"github.com/stretchr/testify/require"
)

const TEST_PRIVATE_KEY = "37b4de3cdd51c013addf9c727be371b3a2258d457552798a76778712bacd01f2"
const TEST_RUNE_ID = "310:1"

// Transfer some runes into an arbitrary wallet and check balances
func testTransfer(t *testing.T, amt int64) {
	config.Init()
	// Init
	var (
		senderRuneBalance   *client.RunesBalance
		receiverRuneBalance *client.RunesBalance

		config      = config.GetDefaultConfig()
		opiClient   = client.NewOpiClient(config.OpiConfig)
		pk          = common.LoadPrivateKey(TEST_PRIVATE_KEY)
		addr, err   = common.GetP2PKHAddress(pk.PubKey().SerializeCompressed(), &chaincfg.RegressionNetParams)
		destAddr, _ = btcutil.DecodeAddress("msTNoMYXynujefraNmpeKKuUsAsSYMZ8C2", &chaincfg.RegressionNetParams)
	)
	require.Nil(t, err)

	// Check Balance - Store balance
	balances, err := opiClient.GetRunesBalance(addr.EncodeAddress())
	require.Nil(t, err)
	for _, i := range balances {
		if i.RuneID == TEST_RUNE_ID {
			senderRuneBalance = i.Copy()
			break
		}
	}
	require.NotNil(t, senderRuneBalance)
	senderPrevBalanceStr := senderRuneBalance.TotalBalance

	// Receiver prev address
	balances, err = opiClient.GetRunesBalance(destAddr.EncodeAddress())
	require.Nil(t, err)
	for _, i := range balances {
		if i.RuneID == TEST_RUNE_ID {
			receiverRuneBalance = i.Copy()
			break
		}
	}
	require.NotNil(t, receiverRuneBalance)
	receiverPrevBalanceStr := receiverRuneBalance.TotalBalance

	// Do Transfer
	rune, err := runes.ParseRune(senderRuneBalance.RuneID)
	require.Nil(t, err)

	_, err = runes.TransferRune(*rune, big.NewInt(amt), addr, destAddr, pk, 22, config)
	require.Nil(t, err)

	// Check balance
	err = cmd.GenerateBlocks()
	require.Nil(t, err)
	// Wait for the indexer to catch up
	_ = time.Tick

	// Sleep so the indexer can catch up on the new block
	// Mainnet has 15min block times, so indexer always catches up
	time.Sleep(3 * time.Second)

	// Get new sender balances
	balances, err = opiClient.GetRunesBalance(addr.EncodeAddress())
	require.Nil(t, err)
	for _, i := range balances {
		if i.RuneID == TEST_RUNE_ID {
			senderRuneBalance = i.Copy()
			break
		}
	}
	senderNewBalanceStr := senderRuneBalance.TotalBalance

	// Get new receiver balances
	balances, err = opiClient.GetRunesBalance(destAddr.EncodeAddress())
	require.Nil(t, err)
	for _, i := range balances {
		if i.RuneID == TEST_RUNE_ID {
			receiverRuneBalance = i.Copy()
			break
		}
	}
	receiverNewBalanceStr := receiverRuneBalance.TotalBalance

	senderPrevBalance, ok := big.NewInt(0).SetString(senderPrevBalanceStr, 10)
	require.True(t, ok)

	receiverPrevBalance, ok := big.NewInt(0).SetString(receiverPrevBalanceStr, 10)
	require.True(t, ok)

	senderNewBalance, ok := big.NewInt(0).SetString(senderNewBalanceStr, 10)
	require.True(t, ok)

	receiverNewBalance, ok := big.NewInt(0).SetString(receiverNewBalanceStr, 10)
	require.True(t, ok)

	// Check if runes are deducted from the sender
	require.Equal(t, big.NewInt(0).Sub(senderPrevBalance, senderNewBalance).String(), big.NewInt(amt).String())

	// Check if runes are added to the receiver
	require.Equal(t, big.NewInt(0).Sub(receiverNewBalance, receiverPrevBalance).String(), big.NewInt(amt).String())
}

// Outputs given by OPI should be valid outpoints
func TestRunesOutpointFetch(t *testing.T) {
	config.Init()
	config := config.GetDefaultConfig()
	pk := common.LoadPrivateKey(TEST_PRIVATE_KEY)
	addr, _ := common.GetP2PKHAddress(pk.PubKey().SerializeCompressed(), &chaincfg.RegressionNetParams)
	opiClient := client.NewOpiClient(config.OpiConfig)
	btcClient := client.NewBitcoinClient(config)

	rune, _ := runes.ParseRune("310:1")

	utxos, err := opiClient.GetRunesUnspentOutpoints(addr.EncodeAddress())
	require.Nil(t, err)

	outpoint, err := runes.SelectRunesUnspentOutput(rune, big.NewInt(100), utxos)
	require.Nil(t, err)

	txId := strings.Split(outpoint.Outpoint, ":")[0]
	_, err = strconv.Atoi(strings.Split(outpoint.Outpoint, ":")[1])
	require.Nil(t, err)

	hash, err := btc.NewHashFromStr(txId)
	require.Nil(t, err)

	_, err = btcClient.GetRawTransactionVerbose(hash)
	require.Nil(t, err)
}

func TestRunesTransfer(t *testing.T) {
	t.Run("transfer 1 rune", func(t *testing.T) {
		testTransfer(t, 1)
	})
	t.Run("transfer 2 rune", func(t *testing.T) {
		testTransfer(t, 2)
	})
	t.Run("transfer 3 rune", func(t *testing.T) {
		testTransfer(t, 3)
	})
	t.Run("transfer 4 rune", func(t *testing.T) {
		testTransfer(t, 4)
	})
	t.Run("transfer 5 rune", func(t *testing.T) {
		testTransfer(t, 5)
	})
}
