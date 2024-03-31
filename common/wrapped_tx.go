package common

import (
	"bytes"
	"fmt"

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

func (x *WrappedTx) EstimateGas(feeRate uint64) (uint64, error) {
	rawTx := x.Copy()
	tx := NewWrappedTx(rawTx, x.SenderPkScript)

	dummySigScript := bytes.Repeat([]byte{0x00}, 105)

	for i := range tx.TxIn {
		tx.TxIn[i].SignatureScript = dummySigScript
	}
	changeTxOut := wire.NewTxOut(5000, tx.SenderPkScript)
	tx.AddTxOut(changeTxOut)

	var buf bytes.Buffer
	err := tx.Serialize(&buf)
	if err != nil {
		fmt.Println("error serializing for size est")
		return 0, err
	}

	totalFee := feeRate * uint64(buf.Len())
	return totalFee, nil
}

func NewWrappedTx(raw *wire.MsgTx, senderPkScript []byte) WrappedTx {
	return WrappedTx{
		raw, senderPkScript,
	}
}
