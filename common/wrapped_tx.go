package common

import (
	"bytes"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
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

// Signs using the given private key and sets the signature in the txin
func (x *WrappedTx) SignP2PKH(privKey *btcec.PrivateKey, pkData []byte, index int) error {
	sigHash, err := x.SigHash(index)
	if err != nil {
		return err
	}
	signature := ecdsa.Sign(privKey, sigHash).Serialize()
	signature = append(signature, byte(txscript.SigHashAll))
	signatureScript, err := txscript.NewScriptBuilder().AddData(signature).AddData(pkData).Script()

	x.TxIn[index].SignatureScript = signatureScript

	return nil
}

func NewWrappedTx(raw *wire.MsgTx, senderPkScript []byte) WrappedTx {
	return WrappedTx{
		raw, senderPkScript,
	}
}
