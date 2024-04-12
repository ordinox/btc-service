package inscriptions

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/ordinox/btc-service/common"
	"github.com/ordinox/btc-service/config"
	"github.com/ordinox/btc-service/taproot"
)

const PK_HEX = "f80c5f4802ac331207fe47f5a36cb4a0c17a0dbd4fe57f4d1243d56e4000e79a"

type deploy struct {
	P    string `json:"p"`
	Op   string `json:"op"`
	Tick string `json:"tick"`
	Max  string `json:"max"`
	Lim  string `json:"lim"`
}

type mint struct {
	P    string `json:"p"`
	Op   string `json:"op"`
	Tick string `json:"tick"`
	Amt  string `json:"amt"`
}

var _ = fmt.Println

func TestInscriptionCommit(t *testing.T) {
	mint := mint{
		P:    "brc-20",
		Op:   "mint",
		Tick: "opiz",
		Amt:  "100",
	}
	str, _ := json.Marshal(&mint)
	config.Init()
	config := config.GetDefaultConfig()
	pk := common.LoadPrivateKey(PK_HEX)
	addr2Str := "moMEj6u1Eb4jXPDVKBfUVwzir2RqhbUcSc"
	addr2, _ := btcutil.DecodeAddress(addr2Str, config.BtcConfig.GetChainConfigParams())
	fmt.Println("Dest =", addr2.EncodeAddress())
	_, err := InscribeNative(
		addr2,
		pk,
		taproot.NewInscriptionData(string(str), taproot.ContentTypeText),
		20,
		config,
	)
	if err != nil {
		panic(err)
	}
}
