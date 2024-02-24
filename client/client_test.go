package client

import (
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ordinox/btc-service/config"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Println

// Should not error out getting a wallet's UTXOs
func TestUtxo(t *testing.T) {
	addr, err := btcutil.DecodeAddress("mtgNH1MPY5QeGobaZXidMvY3dhRjow2AS8", &chaincfg.RegressionNetParams)
	assert.Nil(t, err)
	client := CreateBitcoinClient(config.GetDefaultConfig())
	utxos, err := client.GetUtxos(addr)
	assert.Nil(t, err)
	fmt.Println(len(utxos))
}

func TestOpiClient(t *testing.T) {
	opiClient := NewOpiClient(config.GetDefaultConfig().OpiConfig)
	evts, err := opiClient.GetEventsByInscriptionId("5ccae5f073c372698108f4fc65b56d8f82d0857828a8987d1e6f8b10052f05ffi0")
	if err != nil {
		panic(err)
	}
	for _, evt := range evts {
		fmt.Println(evt.Tick, evt.EventType, evt.SourceWallet, evt.SpentWallet)
	}
}
