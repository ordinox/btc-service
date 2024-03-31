package common

import (
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

type WrappedTx struct {
	*wire.MsgTx
	SenderPkScript []byte
}

func (tx *WrappedTx) SigHash(idx int) ([]byte, error) {
	return txscript.CalcSignatureHash(tx.SenderPkScript, txscript.SigHashAll, tx.MsgTx, idx)
}

func NewWrappedTx(raw *wire.MsgTx, senderPkScript []byte) WrappedTx {
	return WrappedTx{
		raw, senderPkScript,
	}
}
