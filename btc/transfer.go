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
	rawTx := wire.NewMsgTx(wire.TxVersion)
	destinationAddrScript, err := txscript.PayToAddrScript(destinationAddr)
	if err != nil {
		return err
	}
	senderAddrScript, err := txscript.PayToAddrScript(senderAddr)
	if err != nil {
		return err
	}

	tx := common.NewWrappedTx(rawTx, senderAddrScript)

	utxo0, err := wire.NewOutPointFromString(fmt.Sprintf("%s:%d", utxo.GetTxID(), utxo.GetVout()))
	if err != nil {
		return err
	}
	fmt.Println("selected utxo: ", utxo.GetTxID())

	dummySigScript := bytes.Repeat([]byte{0x00}, 105)
	txin0 := wire.NewTxIn(utxo0, dummySigScript, [][]byte{})

	tx.AddTxIn(txin0)
	txout0 := wire.NewTxOut(int64(amtInSats), destinationAddrScript)
	tx.AddTxOut(txout0)

	pkData := senderPrivKey.PubKey().SerializeCompressed()

	totalFee, err := tx.EstimateGas(uint64(feeRate))
	if err != nil {
		return err
	}

	change := int64(utxo.GetValueInSats()) - int64(amtInSats) - int64(totalFee)
	if change < 1 {
		fmt.Printf("change [%d] = totalBal [%d] - amt [%d] - fee [%d]", change, utxo.GetValueInSats(), amtInSats, totalFee)
		return fmt.Errorf("low balance")
	}

	changeTxOut := wire.NewTxOut(int64(change), senderAddrScript)
	tx.AddTxOut(changeTxOut)

	sigHash0, err := tx.SigHash(0)
	if err != nil {
		return err
	}
	sig0 := ecdsa.Sign(&senderPrivKey, sigHash0).Serialize()
	sig0 = append(sig0, byte(txscript.SigHashAll))
	sigScript0, err := txscript.NewScriptBuilder().AddData(sig0).AddData(pkData).Script()
	if err != nil {
		return err
	}

	tx.TxIn[0].SignatureScript = sigScript0

	// TODO: Send tx
	client := client.NewBitcoinClient(config.GetDefaultConfig())
	h, err := client.SendRawTransaction(tx.MsgTx, true)

	if err != nil {
		fmt.Println("Error broadcasting tx")
		fmt.Println(err)
		return err
	}

	fmt.Println("TxHash", h.String())
	fmt.Println("Fee Paid", totalFee)
	fmt.Println("------")
	return nil
}
