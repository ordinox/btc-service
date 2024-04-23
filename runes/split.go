package runes

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/ordinox/btc-service/btc"
	"github.com/ordinox/btc-service/client"
	"github.com/ordinox/btc-service/common"
	"github.com/ordinox/btc-service/config"
)

func Split(addr btc.Address, privateKey *btcec.PrivateKey, utxoIn common.Utxo, outCount, outValue uint64, feeRate uint64, config config.Config) (*chainhash.Hash, error) {
	addr, pubkeyData, err := common.VerifyPrivateKey(privateKey, addr, config.BtcConfig.GetChainConfigParams())
	if err != nil {
		return nil, err
	}
	senderScript, err := btc.PayToAddrScript(addr)
	if err != nil {
		return nil, err
	}
	totalSatsReq := outValue * outCount
	utxo, err := common.SelectOneUtxo(addr.EncodeAddress(), totalSatsReq, config.BtcConfig)
	if err != nil {
		fmt.Printf("Err: Not enough sats: Total Sats Required: %d\n", totalSatsReq)
		return nil, err
	}

	rawTx := btc.NewMsgTx(int32(btc.TxVersion))
	tx := common.NewWrappedTx(rawTx, senderScript)

	outTxId, err := btc.NewHashFromStr(utxo.TxID)
	outPoint := btc.NewOutPoint(outTxId, utxo.Vout)
	txIn0 := btc.NewTxIn(outPoint, btc.DummySig, nil) // Dummy Sig used for gas estimation
	tx.AddTxIn(txIn0)

	count := 0
	bal := utxo.Amount
	for bal > 0 {
		txOut := btc.NewTxOut(int64(outValue), senderScript)
		tx.AddTxOut(txOut)
		bal = bal - float64(outValue)
		count++
	}

	fee, err := tx.EstimateGas(feeRate)
	if err != nil {
		return nil, err
	}

	pruneCount := int64(fee/outValue) + 1

	change := int64(outValue*uint64(pruneCount)) - int64(fee)

	if change < 1 {
		return nil, fmt.Errorf("UTXO Amount is lower than what's needed: change=%d fee=%d totalUtxoAmt=%f", change, fee, utxo.Amount)
	}

	if len(tx.TxOut) >= int(pruneCount) {
		tx.TxOut = tx.TxOut[:len(tx.TxOut)-int(pruneCount)]
	} else {
		panic("UTXOs less than prune count")
	}
	changeTxOut := btc.NewTxOut(int64(change), senderScript)
	tx.AddTxOut(changeTxOut)

	if err := tx.SignP2PKH(privateKey, pubkeyData, 0); err != nil {
		return nil, err
	}

	client := client.NewBitcoinClient(config)
	h, err := client.SendRawTransaction(tx.MsgTx, true)

	if err != nil {
		return nil, err
	}

	fmt.Println("Count:", count)
	fmt.Println("Change:", change)
	fmt.Println("Fee", fee)

	return h, nil
}
