package btc

import (
	"bytes"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ordinox/btc-service/client"
	"github.com/ordinox/btc-service/common"
	"github.com/ordinox/btc-service/config"
)

func TransferBtc(
	senderPrivKey btcec.PrivateKey,
	senderAddr, destinationAddr btcutil.Address,
	utxo common.Utxo,
	amtInSats uint64,
	feeRate uint32,
) error {
	tx := wire.NewMsgTx(wire.TxVersion)
	destinationAddrScript, err := txscript.PayToAddrScript(destinationAddr)
	if err != nil {
		return err
	}
	senderAddrScript, err := txscript.PayToAddrScript(senderAddr)
	if err != nil {
		return err
	}
	utxo0, err := wire.NewOutPointFromString(fmt.Sprintf("%s:%d", utxo.GetTxID(), utxo.GetVout()))
	if err != nil {
		return err
	}
	txin0 := wire.NewTxIn(utxo0, nil, [][]byte{})

	pkData := senderPrivKey.PubKey().SerializeCompressed()

	sigHash0, err := txscript.CalcSignatureHash(senderAddrScript, txscript.SigHashAll, tx, 0)
	if err != nil {
		return err
	}

	tx.AddTxIn(txin0)
	txout0 := wire.NewTxOut(int64(amtInSats), destinationAddrScript)
	tx.AddTxOut(txout0)
	// AmtInSats here is just a placeholder for calculating the size
	txout1 := wire.NewTxOut(int64(amtInSats), senderAddrScript)
	tx.AddTxOut(txout1)

	dummySigScript := bytes.Repeat([]byte{0x00}, 105)
	tx.TxIn[0].SignatureScript = dummySigScript

	var buf bytes.Buffer
	err = tx.Serialize(&buf)
	if err != nil {
		return err
	}

	totalFee := feeRate * uint32(buf.Len())

	// Now that we have the fee, replace the change txout
	// change = total - amt - fee
	change := int64(utxo.GetValueInSats()) - int64(amtInSats) - int64(totalFee)
	if change < 1 {
		fmt.Printf("change [%d] = totalBal [%d] - amt [%d] - fee [%d]", change, utxo.GetValueInSats(), amtInSats, totalFee)
		return fmt.Errorf("low balance")
	}

	tx.TxOut[1].Value = change

	sig0 := ecdsa.Sign(&senderPrivKey, sigHash0).Serialize()
	sig0 = append(sig0, byte(txscript.SigHashAll))
	sigScript0, err := txscript.NewScriptBuilder().AddData(sig0).AddData(pkData).Script()
	if err != nil {
		return err
	}

	tx.TxIn[0].SignatureScript = sigScript0

	// TODO: Send tx
	client := client.NewBitcoinClient(config.GetDefaultConfig())
	h, err := client.SendRawTransaction(tx, true)
	if err != nil {
		fmt.Println("Error broadcasting tx")
		fmt.Println(err)
		return err
	}
	fmt.Println("TxHash", h.String())
	return nil
}
