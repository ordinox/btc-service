package runes

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/wire"
	"github.com/ordinox/btc-service/btc"
	"github.com/ordinox/btc-service/client"
	"github.com/ordinox/btc-service/common"
	"github.com/ordinox/btc-service/config"
)

// Given a list of rune utxos, get the most appropirate utxo for a transfer
func SelectRunesUnspentOutput(rune *Rune, amt *big.Int, utxos []client.RunesUnspentOutput) (*client.RunesUnspentOutput, error) {
	runeStr := rune.String()

	for i, utxo := range utxos {
		for j, rune := range utxo.RuneIds {
			// Check if the rune
			if rune == runeStr && utxos[i].Balances[j].Cmp(amt) >= 0 {
				return &utxo, nil
			}
		}
	}
	return nil, fmt.Errorf("Runes UTXO not found")
}

// Transfer one specific rune from one address to another
// TODO: Port this into a batch transfer function to save on fees
func TransferRune(rune Rune, amount *big.Int, addr btc.Address, toAddr btc.Address, privateKey btc.PrivateKey, feeRate uint64, config config.Config) (btc.Hash, error) {
	var (
		// Flag to identify if the transaction has 2 inputs (FEE UTXO + RUNE UTXO) or 1 (RUNE UTXO)
		includesFeeUtxo = false
	)

	opiClient := client.NewOpiClient(config.OpiConfig)
	addr, pubkeyData, err := common.VerifyPrivateKey(privateKey, addr, config.BtcConfig.GetChainConfigParams())
	if err != nil {
		return nil, err
	}

	// P2PKH Script of the sender
	senderScript, err := btc.PayToAddrScript(addr)
	if err != nil {
		return nil, err
	}

	// P2PKH Script of the receiver
	receiverScript, err := btc.PayToAddrScript(toAddr)
	if err != nil {
		return nil, err
	}

	// Get a utxo for fee payment
	utxo, err := common.SelectOneUtxo(addr.EncodeAddress(), 1000, config.BtcConfig)
	if err != nil {
		return nil, err
	}

	rawTx := btc.NewMsgTx(int32(btc.TxVersion))
	tx := common.NewWrappedTx(rawTx, senderScript)

	outTxId, _ := btc.NewHashFromStr(utxo.TxID)

	outPoint := btc.NewOutPoint(outTxId, utxo.Vout)
	txIn0 := btc.NewTxIn(outPoint, btc.DummySig, nil) // Dummy Sig used for gas estimation

	// This is the fee utxo but, can be also the rune utxo
	tx.AddTxIn(txIn0)

	// Get the rune inputs
	runeUtxos, err := opiClient.GetRunesUnspentOutpoints(addr.EncodeAddress())
	if err != nil {
		return nil, err
	}

	// Get the Rune UTXO for the given Rune
	runeOutpoint, err := SelectRunesUnspentOutput(&rune, amount, runeUtxos)
	if err != nil {
		return nil, err
	}

	runeOutTxIdStr := strings.Split(runeOutpoint.Outpoint, ":")[0]
	vout, err := strconv.Atoi(strings.Split(runeOutpoint.Outpoint, ":")[1])
	if err != nil {
		return nil, err
	}

	// If the Rune UTXO and the cardinal UTXOs are different, then the tx needs to have 2 txins
	if runeOutTxIdStr != utxo.TxID {
		runeOutTxId, err := btc.NewHashFromStr(runeOutTxIdStr)
		if err != nil {
			return nil, err
		}

		runeOutPoint := btc.NewOutPoint(runeOutTxId, uint32(vout))
		txIn1 := btc.NewTxIn(runeOutPoint, btc.DummySig, nil)

		tx.AddTxIn(txIn1)

		includesFeeUtxo = true // Set the fee flag
	}

	// Runestone script
	transferScript, err := createTransferScript(rune, amount, 1, true)
	if err != nil {
		return nil, err
	}

	transferTxOut := btc.NewTxOut(0, transferScript)
	destTxOut := btc.NewTxOut(546, receiverScript) // Min transfer

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

	// If there is a seperate rune txin, sign the second txin as well
	if includesFeeUtxo {
		if err := tx.SignP2PKH(privateKey, pubkeyData, 1); err != nil {
			return nil, err
		}
	}

	client := client.NewBitcoinClient(config)
	h, err := client.SendRawTransaction(tx.MsgTx, true)
	return h, err
}

// Create txout script which contains the runestone
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
