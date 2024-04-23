package runes

import (
	"fmt"

	"github.com/ordinox/btc-service/btc"
	"github.com/ordinox/btc-service/client"
	"github.com/ordinox/btc-service/common"
	"github.com/ordinox/btc-service/config"
)

// Mint runes into a given wallet
// Mostly used only for CLI purposes
func MintRunes(rune Rune, addr btc.Address, privateKey btc.PrivateKey, feeRate uint64, config config.Config) (btc.Hash, error) {
	addr, pubkeyData, err := common.VerifyPrivateKey(privateKey, addr, config.BtcConfig.GetChainConfigParams())
	if err != nil {
		return nil, err
	}
	senderScript, err := btc.PayToAddrScript(addr)
	if err != nil {
		return nil, err
	}
	utxo, err := common.SelectOneUtxo(addr.EncodeAddress(), 1000, config.BtcConfig)
	if err != nil {
		return nil, err
	}

	rawTx := btc.NewMsgTx(int32(btc.TxVersion))
	tx := common.NewWrappedTx(rawTx, senderScript)

	outTxId, err := btc.NewHashFromStr(utxo.TxID)

	outPoint := btc.NewOutPoint(outTxId, utxo.Vout)
	txIn0 := btc.NewTxIn(outPoint, btc.DummySig, nil) // Dummy Sig used for gas estimation

	tx.AddTxIn(txIn0)

	mintScript, err := createMintScript(rune)
	if err != nil {
		return nil, err
	}
	mintTxOut := btc.NewTxOut(0, mintScript)

	tx.AddTxOut(mintTxOut)

	fee, err := tx.EstimateGas(feeRate)
	if err != nil {
		return nil, err
	}

	change := int64(utxo.Amount) - int64(fee)
	changeTxOut := btc.NewTxOut(change, senderScript)
	tx.AddTxOut(changeTxOut)

	if change < 1 {
		return nil, fmt.Errorf("UTXO Amount is lower than what's needed: change=%d fee=%d totalUtxoAmt=%f", change, fee, utxo.Amount)
	}

	if err := tx.SignP2PKH(privateKey, pubkeyData, 0); err != nil {
		return nil, err
	}

	client := client.NewBitcoinClient(config)
	h, err := client.SendRawTransaction(tx.MsgTx, true)

	if err != nil {
		return nil, err
	}

	return h, nil
}

// Create txout to mint the given rune
func createMintScript(rune Rune) ([]byte, error) {
	builder := btc.NewScriptBuilder()
	builder.AddOp(OP_RETURN)
	builder.AddOp(OP_MAGIC)

	data := TagToVarInt(MINT, rune.BlockNumber, uint64(rune.TxIndex))
	builder.AddData(data)
	return builder.Script()
}
