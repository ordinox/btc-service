package brc20

import (
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/ordinox/btc-service/config"
	"github.com/ordinox/btc-service/inscriptions"
	"github.com/ordinox/btc-service/taproot"
)

type deploy struct {
	P    string `json:"p"`
	Op   string `json:"op"`
	Tick string `json:"tick"`
	Max  string `json:"max"`
	Lim  string `json:"lim"`
}

func InscribeDeploy(ticker string, cap uint, privateKey *btcec.PrivateKey, receiver btcutil.Address, feeRate uint64, config config.Config) (*inscriptions.SingleInscriptionResult, error) {
	deploy := deploy{
		P:    "brc-20",
		Op:   "deploy",
		Tick: ticker,
		Max:  fmt.Sprintf("%d", cap),
		Lim:  fmt.Sprintf("%d", cap),
	}

	bz, _ := json.Marshal(deploy)
	inscriptionData := taproot.NewInscriptionData(string(bz), taproot.ContentTypeText)
	return inscriptions.InscribeNative(receiver, privateKey, inscriptionData, feeRate, config)
}
