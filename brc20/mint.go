package brc20

import (
	"encoding/json"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/ordinox/btc-service/config"
	"github.com/ordinox/btc-service/inscriptions"
	"github.com/ordinox/btc-service/taproot"
)

type mint struct {
	P    string `json:"p"`
	Op   string `json:"op"`
	Tick string `json:"tick"`
	Amt  string `json:"amt"`
}

func InscribeMint(ticker string, amt *big.Float, destination btcutil.Address, privateKey *btcec.PrivateKey, feeRate uint64, config config.Config) (*inscriptions.SingleInscriptionResult, error) {
	mint := mint{
		P:    "brc-20",
		Op:   "mint",
		Tick: ticker,
		Amt:  amt.Text('f', 4),
	}

	bz, _ := json.Marshal(mint)
	inscription := taproot.NewInscriptionData(string(bz), taproot.ContentTypeText)
	return inscriptions.InscribeNative(destination, privateKey, inscription, feeRate, config)
}
