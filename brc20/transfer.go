package brc20

// Given a private key (via a keychain), this package should help you
// 1. transfer BRC20 tokens out (handling utxo management for fee payment)
// 2. build txn for signing
// 3. with a transaction id, check if a brc20 token came into the address being monitored
import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ordinox/btc-service/btc"
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

func InscribeTransfer(ticker string, from btcutil.Address, amt, feeRate uint64, config config.Config) (*inscriptions.InscriptionResultRaw, error) {
	transfer := transfer{
		P:    "brc-20",
		Op:   "transfer",
		Tick: ticker,
		Amt:  fmt.Sprintf("%d", amt),
	}
	bz, _ := json.Marshal(transfer)
	return inscriptions.Inscribe(string(bz), from.String(), feeRate, config.BtcConfig)
}

func getUtxos(client *client.BtcRpcClient, from btcutil.Address, inscriptionTxId string, config config.Config) (iUtxo, fUtxo common.Utxo, err error) {
	var utxos []common.Utxo

	if config.BtcConfig.ChainConfig == "mainnet" {
		mUtxos, err := btc.GetUtxos(from.EncodeAddress(), config.BtcConfig)
		if err != nil {
			return nil, nil, err
		}
		utxos = mUtxos.Result.ToUtxo()
	} else {
		rUtxos, err := client.GetUtxos(from)
		if err != nil {
			return nil, nil, err
		}
		utxos = rUtxos.ToUtxo()
	}
	for i := range utxos {
		utxo := utxos[i]
		if fUtxo != nil && iUtxo != nil {
			break
		}
		if utxo.GetValueInSats() == 546 && inscriptionTxId == utxo.GetTxID() {
			iUtxo = utxo
		} else if utxo.GetValueInSats() > 6500 {
			// TODO: Check if this can be potential inscription
			fUtxo = utxo
		}
	}
	return
}

func TransferBrc20(from, to btcutil.Address, inscriptionId string, amt uint64, privKey btcec.PrivateKey, feeRate uint64, config config.Config) (*string, error) {
	fmt.Printf("--transferring brc20 from=%s to=%s", from.String(), to.String())
	client := client.NewBitcoinClient(config)

	inscriptionTxId := inscriptionId[:len(inscriptionId)-2]
	var inscriptionUtxo common.Utxo
	var feeUtxo common.Utxo

	count := 0
	// 120 second backoff till utxo is found in the mempool
	for {
		fmt.Println("Finding UTXOs - Attempt ", count+1)
		if count == 120 {
			fmt.Printf("-- Err InscriptionUtxoFound? %t  FeeUtxoFound? %t \n", inscriptionUtxo != nil, feeUtxo != nil)
			return nil, fmt.Errorf("couldn't finalise inscription/fee UTXO within the backoff time (120s)")
		}
		var err error
		inscriptionUtxo, feeUtxo, err = getUtxos(client, from, inscriptionTxId, config)
		if err != nil {
			fmt.Printf("-- err getting utxos InscriptionUtxoFound? %t FeeUtxoFound? %t \n", inscriptionUtxo != nil, feeUtxo != nil)
			return nil, err
		}
		if inscriptionUtxo != nil && feeUtxo != nil {
			fmt.Printf("-- BREAKING InscriptionUtxoFound? %t  FeeUtxoFound? %t \n", inscriptionUtxo != nil, feeUtxo != nil)
			break
		}
		count = count + 1
		time.Sleep(1 * time.Second)
	}

	if feeUtxo == nil {
		return nil, fmt.Errorf("no fee utxo found")
	}

	hash, err := Transfer(feeUtxo, inscriptionUtxo, from, to, &privKey, privKey.PubKey(), feeRate, config)
	if err != nil {
		return nil, err
	}
	res := "transfer complete with hash: " + hash

	// fmt.Println(hash.String())

	return &res, nil
}

// Use UTXOs of the given wallet to transfer an inscription
func Transfer(cUtxo, iUtxo common.Utxo, sendderAddr, destAddr btcutil.Address, senderPk *btcec.PrivateKey, senderPubKey *btcec.PublicKey, feeRate uint64, config config.Config) (string, error) {
	fmt.Println("Transfer Called ")
	tx, err := BuildTransferTx(cUtxo, iUtxo, sendderAddr, destAddr)
	if err != nil {
		log.Err(err).Msg("error building MsgTx")
		return "", err
	}
	gas, err := tx.EstimateGas(feeRate)
	if err != nil {
		log.Err(err).Msg("error estimating gas")
		return "", err
	}

	change := cUtxo.GetValueInSats() - gas

	changeTxOut := wire.NewTxOut(int64(change), tx.SenderPkScript)
	tx.AddTxOut(changeTxOut)

	pkData := senderPubKey.SerializeCompressed()
	sigHash0, err := tx.SigHash(0)
	if err != nil {
		log.Err(err).Msg("error generating sighash0")
		return "", err
	}

	sig0 := ecdsa.Sign(senderPk, sigHash0).Serialize()
	sig0 = append(sig0, byte(txscript.SigHashAll))
	sigScript0, err := txscript.NewScriptBuilder().AddData(sig0).AddData(pkData).Script()
	if err != nil {
		log.Err(err).Msg("error building sigScript0")
		return "", err
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

	client := client.NewBitcoinClient(config)
	h, err := client.SendRawTransaction(tx.MsgTx, true)
	if err != nil {
		log.Err(err).Msg("error broadcasting txn")
		return "", err
	}
	log.Info().Msgf("hash: %s", h.String())
	log.Info().Msgf("fee: %d", gas)
	return h.String(), nil
}

func SendBrc20(ticker string, from, to btcutil.Address, amt, feeRate uint64, privKey btcec.PrivateKey, config config.Config) (inscriptionId, hash string, err error) {
	res, err := InscribeTransfer(ticker, from, amt, feeRate, config)
	if err != nil {
		return "", "", err
	}
	if res.Inscriptions == nil || len(res.Inscriptions) == 0 {
		// This usually means we failed to acquire ord lock
		// Sleep and try again
		time.Sleep(100 * time.Millisecond)
		return SendBrc20(ticker, from, to, amt, feeRate, privKey, config)
	}
	res2, err := TransferBrc20(from, to, res.Inscriptions[0].Id, amt, privKey, feeRate, config)
	if err != nil {
		return inscriptionId, "", err
	}
	hash = *res2
	err = nil
	return
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
