package common

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec/v2"
)

func LoadPrivateKey(pkHex string) *btcec.PrivateKey {
	pkBytes, _ := hex.DecodeString(pkHex)
	pk, _ := btcec.PrivKeyFromBytes(pkBytes)
	return pk
}
