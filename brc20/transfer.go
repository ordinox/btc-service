package brc20

// Given a private key (via a keychain), this package should help you
// 1. transfer BRC20 tokens out (handling utxo management for fee payment)
// 2. build txn for signing
// 3. with a transaction id, check if a brc20 token came into the address being monitored
import (
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ordinox/btc-service/client"
	"github.com/ordinox/btc-service/common"
	"github.com/ordinox/btc-service/config"
	"github.com/rs/zerolog/log"
)

// Use UTXOs of the given wallet to transfer an inscription
func Transfer(cUtxo, iUtxo btcjson.ListUnspentResult, sendderAddr, destAddr btcutil.Address, senderPk *btcec.PrivateKey, senderPubKey *btcec.PublicKey) error {
	tx, err := BuildTransferTx(cUtxo, iUtxo, sendderAddr, destAddr)
	if err != nil {
		log.Err(err).Msg("error building MsgTx")
		return err
	}
	pkData := senderPubKey.SerializeUncompressed()
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

	client := client.CreateBitcoinClient(config.GetDefaultConfig())
	h, err := client.SendRawTransaction(tx.MsgTx, true)
	if err != nil {
		log.Err(err).Msg("error broadcasting txn")
		return err
	}
	log.Info().Msgf("hash: %s", h.String())
	return nil
}

// 89e68ee66bbed960bd2ac69159bce2d188c8a1e19c6196de7ce3e7dfe91ecb9e

func BuildTransferTx(cardinalUtxo, inscriptionUtxo btcjson.ListUnspentResult, senderAddr, destinationAddr btcutil.Address) (*common.WrappedTx, error) {
	tx := wire.NewMsgTx(wire.TxVersion)
	destinationAddrScript, err := txscript.PayToAddrScript(destinationAddr)
	if err != nil {
		return nil, err
	}
	senderAddrScript, err := txscript.PayToAddrScript(senderAddr)
	if err != nil {
		return nil, err
	}

	utxo0, err := wire.NewOutPointFromString(fmt.Sprintf("%s:%d", inscriptionUtxo.TxID, inscriptionUtxo.Vout))
	utxo1, err := wire.NewOutPointFromString(fmt.Sprintf("%s:%d", cardinalUtxo.TxID, cardinalUtxo.Vout))

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
