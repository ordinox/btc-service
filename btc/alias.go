package btc

import (
	"bytes"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

var (
	NewMsgTx    = wire.NewMsgTx
	NewTxIn     = wire.NewTxIn
	NewTxOut    = wire.NewTxOut
	TxVersion   = wire.TxVersion
	NewOutPoint = wire.NewOutPoint

	PayToAddrScript  = txscript.PayToAddrScript
	NewScriptBuilder = txscript.NewScriptBuilder

	NewHashFromStr = chainhash.NewHashFromStr

	DummySig = bytes.Repeat([]byte{0x00}, 105)

	Sign = ecdsa.Sign
)

type (
	PrivateKey *btcec.PrivateKey
	Address    btcutil.Address
	Hash       *chainhash.Hash
)
