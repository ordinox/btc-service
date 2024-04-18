package runes

import (
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/wire"
	"github.com/ordinox/btc-service/btc"
	"github.com/ordinox/btc-service/client"
	"github.com/ordinox/btc-service/common"
	"github.com/ordinox/btc-service/config"
)

func TransferRune(rune Rune, amount *big.Int, addr btc.Address, toAddr btc.Address, privateKey btc.PrivateKey, feeRate uint64, config config.Config) (btc.Hash, error) {
	addr, pubkeyData, err := common.VerifyPrivateKey(privateKey, addr, config.BtcConfig.GetChainConfigParams())
	if err != nil {
		return nil, err
	}
	senderScript, err := btc.PayToAddrScript(addr)
	if err != nil {
		return nil, err
	}

	receiverScript, err := btc.PayToAddrScript(toAddr)
	if err != nil {
		return nil, err
	}

	utxo, err := common.SelectOneUtxo(addr.EncodeAddress(), 1000, config.BtcConfig)
	if err != nil {
		return nil, err
	}

	rawTx := btc.NewMsgTx(int32(btc.TxVersion))
	tx := common.NewWrappedTx(rawTx, senderScript)

	outTxId, _ := btc.NewHashFromStr(utxo.TxID)

	outPoint := btc.NewOutPoint(outTxId, utxo.Vout)
	txIn0 := btc.NewTxIn(outPoint, btc.DummySig, nil) // Dummy Sig used for gas estimation

	tx.AddTxIn(txIn0)

	transferScript, err := createTransferScript(rune, amount, 1, true)
	if err != nil {
		return nil, err
	}

	transferTxOut := btc.NewTxOut(0, transferScript)
	destTxOut := btc.NewTxOut(546, receiverScript)

	tx.AddTxOut(transferTxOut)
	tx.AddTxOut(destTxOut)

	fee, err := tx.EstimateGas(feeRate)
	if err != nil {
		return nil, err
	}

	change := int64(utxo.Amount) - int64(546) - int64(fee)
	if change < 1 {
		return nil, fmt.Errorf("UTXO Amount is lower than what's needed: change=%d fee=%d amt=%d totalUtxoAmt=%f", change, fee, 546, utxo.Amount)
	}

	changeTxOut := btc.NewTxOut(change, senderScript)

	// Reset TxOut
	tx.TxOut = make([]*wire.TxOut, 0)
	tx.AddTxOut(changeTxOut)
	tx.AddTxOut(destTxOut)
	tx.AddTxOut(transferTxOut)

	if err := tx.SignP2PKH(privateKey, pubkeyData, 0); err != nil {
		return nil, err
	}

	client := client.NewBitcoinClient(config)
	h, err := client.SendRawTransaction(tx.MsgTx, true)
	return h, err
}

// TODO: Handle multiple edicts being sent in the same txn
func createTransferScript(rune Rune, amount *big.Int, output uint64, shouldInit bool) ([]byte, error) {
	scriptBuilder := btc.NewScriptBuilder()
	if shouldInit {
		scriptBuilder.AddOp(OP_RETURN)
		scriptBuilder.AddOp(OP_MAGIC)
	}

	data := ToVarInt(BODY)
	data = append(data, NewEdict(rune, amount, output)...)
	return scriptBuilder.AddData(data).Script()
}
