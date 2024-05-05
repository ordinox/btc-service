package inscriptions

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ordinox/btc-service/btc"
	"github.com/ordinox/btc-service/client"
	"github.com/ordinox/btc-service/common"
	"github.com/ordinox/btc-service/config"
	"github.com/ordinox/btc-service/taproot"
)

var defaultSequenceNum = btc.MaxTxInSequenceNum - 10

func InscribeNative(
	receiver btcutil.Address,
	privateKey *btcec.PrivateKey,
	inscriptionData taproot.InscriptionData,
	feeRate uint64,
	config config.Config,
) (*SingleInscriptionResult, error) {
	commitTx := btc.NewMsgTx(int32(btc.TxVersion))
	client := client.NewBitcoinClient(config)
	fromAddr, err := btc.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootKeyNoScript(privateKey.PubKey())), config.BtcConfig.GetChainConfigParams())
	if err != nil {
		return nil, fmt.Errorf("error deriving taproot addresss, %s", err.Error())
	}
	utxo, err := common.SelectOneUtxo(fromAddr.EncodeAddress(), 1000, config.BtcConfig)
	if err != nil {
		return nil, fmt.Errorf("error selecting utxo, %s", err.Error())
	}

	txData, err := GetTxData(client, utxo.TxID)
	if err != nil {
		return nil, fmt.Errorf("error getting txdata, %s", err.Error())
	}
	// Note: Never use prevVout.Value
	prevVout := txData.Vout[utxo.Vout]
	prevVoutScriptPubkey, err := hex.DecodeString(prevVout.ScriptPubKey.Hex)
	if err != nil {
		return nil, fmt.Errorf("error decoding prev-vout-script-pubkey, %s", err.Error())
	}

	prevTxOut := btc.NewTxOut(int64(utxo.Amount), prevVoutScriptPubkey)
	preVoutTxHash, err := chainhash.NewHashFromStr(utxo.TxID)
	if err != nil {
		return nil, fmt.Errorf("error createing hash from prev-vout-tx-id, %s", err.Error())
	}

	outPoint := btc.NewOutPoint(preVoutTxHash, utxo.Vout)

	in := btc.NewTxIn(outPoint, nil, nil)
	in.Sequence = defaultSequenceNum

	commitTx.AddTxIn(in)

	inscriptionMetaData, err := taproot.CreateP2TRInscriptionMetaData(inscriptionData, privateKey, config)
	if err != nil {
		return nil, fmt.Errorf("error creating inscription meta data, %s", err.Error())
	}

	// Inscription commit txout
	commitTx.AddTxOut(btc.NewTxOut(546, inscriptionMetaData.PkScript))

	totalSenderAmt := btc.Amount(utxo.Amount)
	// Change txout
	commitTx.AddTxOut(btc.NewTxOut(0, prevVoutScriptPubkey))

	fee := btc.Amount(btc.GetTxVSize(btc.NewTx(commitTx)) * int64(feeRate))
	payForward := fee + 546
	changeAmount := totalSenderAmt - fee - payForward - 100 // 100 is the buffer amt to get the commit tx before the reveal tx

	if changeAmount < 1 {
		return nil, fmt.Errorf("selected UTXO has insufficient balance, %d", changeAmount)
	}

	// Update the inscription txout value to the payforward value
	commitTx.TxOut[0].Value = int64(payForward)
	// Update the change value
	commitTx.TxOut[len(commitTx.TxOut)-1].Value = int64(changeAmount)

	// Signing
	prevOutputFetcher := txscript.NewMultiPrevOutFetcher(nil)
	prevOutputFetcher.AddPrevOut(*outPoint, prevTxOut)

	witness, err := txscript.TaprootWitnessSignature(
		commitTx,
		txscript.NewTxSigHashes(commitTx, prevOutputFetcher),
		0,
		int64(utxo.Amount),
		prevVoutScriptPubkey,
		txscript.SigHashDefault,
		privateKey,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating taproot-witness-signature, %s", err.Error())
	}
	commitTx.TxIn[0].Witness = witness

	commitTxOut := btc.NewTxOut(
		int64(payForward),
		inscriptionMetaData.PkScript,
	)

	revealTx := btc.NewMsgTx(int32(btc.TxVersion))

	commitTxHash := commitTx.TxHash()
	commitOutpoint := btc.NewOutPoint(&commitTxHash, 0)

	revealIn := btc.NewTxIn(commitOutpoint, nil, nil)
	revealTx.AddTxIn(revealIn)

	recieverPkScript, _ := btc.PayToAddrScript(receiver)
	revealTx.AddTxOut(btc.NewTxOut(546, recieverPkScript))

	revealTxOutputFetcher := txscript.NewMultiPrevOutFetcher(nil)
	revealTxOutputFetcher.AddPrevOut(*commitOutpoint, commitTxOut)

	witnessArray, err := txscript.CalcTapscriptSignaturehash(
		txscript.NewTxSigHashes(revealTx, revealTxOutputFetcher),
		txscript.SigHashDefault,
		revealTx,
		0,
		revealTxOutputFetcher,
		txscript.NewBaseTapLeaf(inscriptionMetaData.LockScript),
	)
	if err != nil {
		return nil, fmt.Errorf("error constructing witness array, %s", err.Error())
	}

	signature, err := schnorr.Sign(privateKey, witnessArray)
	if err != nil {
		return nil, fmt.Errorf("error signing witness array, %s", err.Error())
	}

	revealTx.TxIn[0].Witness = wire.TxWitness{
		signature.Serialize(),
		inscriptionMetaData.LockScript,
		inscriptionMetaData.ControlBlockWitness,
	}

	h1, err := client.SendRawTransaction(commitTx, true)
	if err != nil {
		return nil, fmt.Errorf("error sending commit tx, %s", err.Error())
	}
	fmt.Println("Commit Tx:", (*h1).String())

	h2, err := client.SendRawTransaction(revealTx, true)
	if err != nil {
		return nil, fmt.Errorf("error sending reveal tx, %s", err.Error())
	}
	result := &SingleInscriptionResult{
		TotalFeePaid: int64(fee) + int64(fee) - 546,
		CommitTx:     (*h1).String(),
		RevealTx:     (*h2).String(),
	}
	return result, nil
}

func GetTxData(client *client.BtcRpcClient, hash string) (btc.TxData, error) {
	cHash, err := btc.NewHashFromStr(hash)
	if err != nil {
		return nil, err
	}
	return client.GetRawTransactionVerbose(cHash)
}
