package brc20

// Given a private key (via a keychain), this package should help you
// 1. transfer BRC20 tokens out (handling utxo management for fee payment)
// 2. build txn for signing
// 3. with a transaction id, check if a brc20 token came into the address being monitored
import (
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ordinox/btc-service/client"
	"github.com/ordinox/btc-service/common"
	"github.com/ordinox/btc-service/config"
	"github.com/ordinox/btc-service/inscriptions"
	"github.com/rs/zerolog/log"
)

type transfer struct {
	P    string `json:"p"`
	Op   string `json:"op"`
	Tick string `json:"tick"`
	Amt  string `json:"amt"`
}

func InscribeTransfer(ticker string, from btcutil.Address, amt, feeRate uint, config config.Config) (*inscriptions.InscriptionResultRaw, error) {
	transfer := transfer{
		P:    "brc-20",
		Op:   "transfer",
		Tick: ticker,
		Amt:  fmt.Sprintf("%d", amt),
	}
	bz, _ := json.Marshal(transfer)
	return inscriptions.Inscribe(string(bz), from.String(), feeRate, config.BtcConfig)

	client := client.NewBitcoinClient(config)
	err := client.ImportAddress(from.String())
	if err != nil {
		log.Err(err).Msg("error importing address for tracking")
		return nil, err
	}
	utxos, err := client.GetUtxos(from)
	if err != nil {
		log.Err(err).Msgf("error getting utxos for address: %s", from.String())
	}
	var cardinal common.Utxo
	_ = cardinal
	for _, utxo := range utxos {
		fmt.Println(utxo.TxID)
		if utxo.Amount == 546e-8 {
			// Then this is an inscription
			continue
		}
		// First non-inscription
		cardinal = utxo
	}
	return nil, nil

}

func TransferBrc20(from, to btcutil.Address, inscriptionId string, amt uint, privKey btcec.PrivateKey, feeRate uint64, config config.Config) (*string, error) {
	// TODO: Calculate change
	client := client.NewBitcoinClient(config)
	utxos, err := client.GetUtxos(from)
	if err != nil {
		return nil, err
	}
	inscriptionTxId := inscriptionId[:len(inscriptionId)-2]
	var inscriptionUtxo common.Utxo
	var feeUtxo common.Utxo
	for i := range utxos {
		utxo := utxos[i]
		if feeUtxo != nil && inscriptionUtxo != nil {
			break
		}
		if utxo.Amount == 546e-8 && inscriptionTxId == utxo.TxID {
			inscriptionUtxo = &utxo
		} else if utxo.GetValueInSats() > 6500 {
			// TODO: Check if this can be potential inscription
			feeUtxo = &utxo
		}
	}
	if inscriptionUtxo == nil {
		return nil, fmt.Errorf("inscription utxo not found")
	}
	if feeUtxo == nil {
		return nil, fmt.Errorf("no fee utxo found")
	}

	err = Transfer(feeUtxo, inscriptionUtxo, from, to, &privKey, privKey.PubKey(), feeRate)
	if err != nil {
		return nil, err
	}
	res := "transfer complete"

	// fmt.Println(hash.String())

	return &res, nil
}

// Use UTXOs of the given wallet to transfer an inscription
func Transfer(cUtxo, iUtxo common.Utxo, sendderAddr, destAddr btcutil.Address, senderPk *btcec.PrivateKey, senderPubKey *btcec.PublicKey, feeRate uint64) error {
	tx, err := BuildTransferTx(cUtxo, iUtxo, sendderAddr, destAddr)
	if err != nil {
		log.Err(err).Msg("error building MsgTx")
		return err
	}
	gas, err := tx.EstimateGas(feeRate)
	if err != nil {
		log.Err(err).Msg("error estimating gas")
		return err
	}

	change := cUtxo.GetValueInSats() - gas

	changeTxOut := wire.NewTxOut(int64(change), tx.SenderPkScript)
	tx.AddTxOut(changeTxOut)

	pkData := senderPubKey.SerializeCompressed()
	sigHash0, err := tx.SigHash(0)
	if err != nil {
		log.Err(err).Msg("error generating sighash0")
		return err
	}

	sig0 := ecdsa.Sign(senderPk, sigHash0).Serialize()
	sig0 = append(sig0, byte(txscript.SigHashAll))
	sigScript0, err := txscript.NewScriptBuilder().AddData(sig0).AddData(pkData).Script()
	if err != nil {
		log.Err(err).Msg("error building sigScript0")
		return err
	}

	// Prepare sigscript - uncompressed

	sigHash1, err := tx.SigHash(1)
	if err != nil {
		log.Err(err).Msg("error generating sighash1")
	}
	sig1 := ecdsa.Sign(senderPk, sigHash1).Serialize()
	sig1 = append(sig1, byte(txscript.SigHashAll))
	sigScript1, err := txscript.NewScriptBuilder().AddData(sig1).AddData(pkData).Script()
	if err != nil {
		log.Err(err).Msg("error building sigScript1")
	}

	tx.TxIn[0].SignatureScript = sigScript0
	tx.TxIn[1].SignatureScript = sigScript1

	client := client.NewBitcoinClient(config.GetDefaultConfig())
	h, err := client.SendRawTransaction(tx.MsgTx, true)
	if err != nil {
		log.Err(err).Msg("error broadcasting txn")
		return err
	}
	log.Info().Msgf("hash: %s", h.String())
	return nil
}

// 89e68ee66bbed960bd2ac69159bce2d188c8a1e19c6196de7ce3e7dfe91ecb9e

func BuildTransferTx(cardinalUtxo, inscriptionUtxo common.Utxo, senderAddr, destinationAddr btcutil.Address) (*common.WrappedTx, error) {
	tx := wire.NewMsgTx(wire.TxVersion)
	destinationAddrScript, err := txscript.PayToAddrScript(destinationAddr)
	if err != nil {
		return nil, err
	}
	senderAddrScript, err := txscript.PayToAddrScript(senderAddr)
	if err != nil {
		return nil, err
	}

	utxo0, err := wire.NewOutPointFromString(fmt.Sprintf("%s:%d", inscriptionUtxo.GetTxID(), inscriptionUtxo.GetVout()))
	utxo1, err := wire.NewOutPointFromString(fmt.Sprintf("%s:%d", cardinalUtxo.GetTxID(), cardinalUtxo.GetVout()))

	txin0 := wire.NewTxIn(utxo0, nil, [][]byte{})
	tx.AddTxIn(txin0)
	txin1 := wire.NewTxIn(utxo1, nil, [][]byte{})
	tx.AddTxIn(txin1)
	txout := wire.NewTxOut(546, destinationAddrScript)
	tx.AddTxOut(txout)
	return &common.WrappedTx{
		MsgTx:          tx,
		SenderPkScript: senderAddrScript,
	}, nil
}

func ListInscriptionUtxos() {}

func ListCardinalUtxo() {}
