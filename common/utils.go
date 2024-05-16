package common

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"strconv"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

func LoadPrivateKey(pkHex string) *btcec.PrivateKey {
	pkBytes, _ := hex.DecodeString(pkHex)
	pk, _ := btcec.PrivKeyFromBytes(pkBytes)
	return pk
}

func GetP2PKHAddress(pubkeyData []byte, chaincfg *chaincfg.Params) (*btcutil.AddressPubKeyHash, error) {
	derivedAddr, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(pubkeyData), chaincfg)
	if err != nil {
		return nil, err
	}
	return derivedAddr, nil
}

// Verify that the address belongs to the private key, regardless of pubkey compression
// P2PKH (Compressed or Uncompressed)
func VerifyPrivateKey(privateKey *btcec.PrivateKey, p2pkhAddr btcutil.Address, chainCfg *chaincfg.Params) (btcutil.Address, []byte, error) {
	pubkey := privateKey.PubKey()
	pubkeyData := pubkey.SerializeCompressed()
	derivedAddr, err := GetP2PKHAddress(pubkeyData, chainCfg)
	if err != nil {
		return nil, nil, err
	}
	if derivedAddr.EncodeAddress() != p2pkhAddr.EncodeAddress() {
		pubkeyData = pubkey.SerializeUncompressed()
		derivedAddr, err = btcutil.NewAddressPubKeyHash(btcutil.Hash160(pubkeyData), chainCfg)
		if err != nil {
			return nil, nil, err
		}
		if derivedAddr.EncodeAddress() != p2pkhAddr.EncodeAddress() {
			return nil, nil, fmt.Errorf("private key does not match the address")
		}
	}
	return derivedAddr, pubkeyData, nil
}

// Given a token value in big int, parse use the decimals provided to make it into a float64
func ParseStringFloat64(amtStr string, decInt int) (float64, error) {
	amt, ok := new(big.Int).SetString(amtStr, 10)
	if !ok {
		return 0, fmt.Errorf("invalid amount string")
	}
	amtF := new(big.Float).SetInt(amt)
	dec := new(big.Float).SetFloat64(math.Pow10(decInt))
	ans := new(big.Float).Quo(amtF, dec)
	f, err := strconv.ParseFloat(ans.Text('f', 4), 64)
	if err != nil {
		return 0, err
	}
	return f, nil
}
