package brc20

import (
	"encoding/json"
	"fmt"

	"github.com/ordinox/btc-service/config"
	"github.com/ordinox/btc-service/inscriptions"
)

type mint struct {
	P    string `json:"p"`
	Op   string `json:"op"`
	Tick string `json:"tick"`
	Amt  string `json:"amt"`
}

func InscribeMint(ticker string, cap uint, destination string, feeRate uint64, config config.BtcConfig) (*inscriptions.InscriptionResultRaw, error) {
	mint := mint{
		P:    "brc-20",
		Op:   "mint",
		Tick: ticker,
		Amt:  fmt.Sprintf("%d", cap),
	}

	bz, _ := json.Marshal(mint)
	inscription := string(bz)
	return inscriptions.Inscribe(inscription, destination, feeRate, config)
}
