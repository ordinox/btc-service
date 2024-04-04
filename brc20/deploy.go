package brc20

import (
	"encoding/json"
	"fmt"

	"github.com/ordinox/btc-service/config"
	"github.com/ordinox/btc-service/inscriptions"
)

type deploy struct {
	P    string `json:"p"`
	Op   string `json:"op"`
	Tick string `json:"tick"`
	Max  string `json:"max"`
	Lim  string `json:"lim"`
}

func InscribeDeploy(ticker string, cap uint, destination string, feeRate uint64, config config.BtcConfig) (*inscriptions.InscriptionResultRaw, error) {
	deploy := deploy{
		P:    "brc-20",
		Op:   "deploy",
		Tick: ticker,
		Max:  fmt.Sprintf("%d", cap),
		Lim:  fmt.Sprintf("%d", cap),
	}

	bz, _ := json.Marshal(deploy)
	inscription := string(bz)
	return inscriptions.Inscribe(inscription, destination, feeRate, config)
}
