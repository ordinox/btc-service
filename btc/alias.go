package btc

import (
	"bytes"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/mempool"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

var (
	NewMsgTx           = wire.NewMsgTx
	NewTxIn            = wire.NewTxIn
	NewTxOut           = wire.NewTxOut
	TxVersion          = wire.TxVersion
	NewOutPoint        = wire.NewOutPoint
	MaxTxInSequenceNum = wire.MaxTxInSequenceNum

	PayToAddrScript     = txscript.PayToAddrScript
	NewScriptBuilder    = txscript.NewScriptBuilder
	TweakTaprootPrivKey = txscript.TweakTaprootPrivKey
	SigHashDefault      = txscript.SigHashDefault

	NewHashFromStr = chainhash.NewHashFromStr

	DummySig = bytes.Repeat([]byte{0x00}, 105)

	Sign = ecdsa.Sign

	NewPrivateKey     = btcec.NewPrivateKey
	NewAddressTaproot = btcutil.NewAddressTaproot
	NewWIF            = btcutil.NewWIF
	NewTx             = btcutil.NewTx
	GetTxVSize        = mempool.GetTxVirtualSize
)

var (
	OP_CHECKSIG = txscript.OP_CHECKSIG
	OP_FALSE    = txscript.OP_FALSE
	OP_IF       = txscript.OP_IF
	OP_DATA_1   = txscript.OP_DATA_1
	OP_0        = txscript.OP_0
	OP_ENDIF    = txscript.OP_ENDIF
)

type (
	PrivateKey *btcec.PrivateKey
	Address    btcutil.Address
	Hash       *chainhash.Hash
	Amount     btcutil.Amount
	TxData     *btcjson.TxRawResult
)
