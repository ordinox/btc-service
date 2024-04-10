package common

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

func LoadPrivateKey(pkHex string) *btcec.PrivateKey {
	pkBytes, _ := hex.DecodeString(pkHex)
	pk, _ := btcec.PrivKeyFromBytes(pkBytes)
	return pk
}

// Verify that the address belongs to the private key, regardless of pubkey compression
func VerifyPrivateKey(privateKey *btcec.PrivateKey, addr btcutil.Address, chainCfg *chaincfg.Params) (btcutil.Address, []byte, error) {
	pubkey := privateKey.PubKey()
	pubkeyData := pubkey.SerializeCompressed()
	derivedAddr, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(pubkeyData), chainCfg)
	if err != nil {
		return nil, nil, err
	}

	if derivedAddr.EncodeAddress() != addr.EncodeAddress() {
		pubkeyData = pubkey.SerializeUncompressed()
		derivedAddr, err = btcutil.NewAddressPubKeyHash(btcutil.Hash160(pubkeyData), chainCfg)
		if err != nil {
			return nil, nil, err
		}
		if derivedAddr.EncodeAddress() != addr.EncodeAddress() {
			return nil, nil, fmt.Errorf("private key does not match the address")
		}
	}
	return derivedAddr, pubkeyData, nil
}
