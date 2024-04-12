package brc20

// Given a private key (via a keychain), this package should help you
// 1. transfer BRC20 tokens out (handling utxo management for fee payment)
// 2. build txn for signing
// 3. with a transaction id, check if a brc20 token came into the address being monitored
import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ordinox/btc-service/client"
	"github.com/ordinox/btc-service/common"
	"github.com/ordinox/btc-service/config"
	"github.com/ordinox/btc-service/inscriptions"
	"github.com/ordinox/btc-service/taproot"
	"github.com/rs/zerolog/log"
)

type transfer struct {
	P    string `json:"p"`
	Op   string `json:"op"`
	Tick string `json:"tick"`
	Amt  string `json:"amt"`
}

func InscribeTransfer(ticker string, amt uint64, destination btcutil.Address, privateKey *btcec.PrivateKey, feeRate uint64, config config.Config) (*inscriptions.SingleInscriptionResult, error) {
	transfer := transfer{
		P:    "brc-20",
		Op:   "transfer",
		Tick: ticker,
		Amt:  fmt.Sprintf("%d", amt),
	}
	bz, _ := json.Marshal(transfer)
	inscription := taproot.NewInscriptionData(string(bz), taproot.ContentTypeText)
	return inscriptions.InscribeNative(destination, privateKey, inscription, feeRate, config)
}

func getUtxos(client *client.BtcRpcClient, from btcutil.Address, inscriptionTxId string, config config.Config) (iUtxo, fUtxo common.Utxo, err error) {
	var utxos []common.Utxo
	mUtxos, err := common.GetUtxos(from.EncodeAddress(), config.BtcConfig)
	if err != nil {
		return nil, nil, err
	}
	utxos = mUtxos.Result.ToUtxo()
	for i := range utxos {
		utxo := utxos[i]
		if fUtxo != nil && iUtxo != nil {
			break
		}
		if inscriptionTxId == utxo.GetTxID() {
			iUtxo = utxo
		} else if utxo.GetValueInSats() > 6500 {
			// TODO: Check if this can be potential inscription
			fUtxo = utxo
		}
	}
	return
}

func TransferBrc20(from, to btcutil.Address, inscriptionId string, privKey *btcec.PrivateKey, feeRate uint64, config config.Config) (*string, error) {
	fmt.Printf("--transferring brc20 from=%s to=%s", from.String(), to.String())
	client := client.NewBitcoinClient(config)

	inscriptionTxId := strings.TrimRight(inscriptionId, "i0")
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
			fmt.Printf("utxos found")
			break
		}
		count = count + 1
		time.Sleep(1 * time.Second)
	}

	if feeUtxo == nil {
		return nil, fmt.Errorf("no fee utxo found")
	}

	hash, err := Transfer(feeUtxo, inscriptionUtxo, from, to, privKey, privKey.PubKey(), feeRate, config)
	if err != nil {
		return nil, err
	}
	// fmt.Println(hash.String())

	return &hash, nil
}

// Use UTXOs of the given wallet to transfer an inscription
func Transfer(cUtxo, iUtxo common.Utxo, senderAddr, destAddr btcutil.Address, senderPk *btcec.PrivateKey, senderPubKey *btcec.PublicKey, feeRate uint64, config config.Config) (string, error) {
	fmt.Println("Transfer Called ")
	senderAddr, senderPkData, err := common.VerifyPrivateKey(senderPk, senderAddr, config.BtcConfig.GetChainConfigParams())
	if err != nil {
		return "", fmt.Errorf("error verifying privatekey: %s", err.Error())
	}
	tx, err := BuildTransferTx(cUtxo, iUtxo, senderAddr, destAddr)
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

	sigHash0, err := tx.SigHash(0)
	if err != nil {
		log.Err(err).Msg("error generating sighash0")
		return "", err
	}

	sig0 := ecdsa.Sign(senderPk, sigHash0).Serialize()
	sig0 = append(sig0, byte(txscript.SigHashAll))
	sigScript0, err := txscript.NewScriptBuilder().AddData(sig0).AddData(senderPkData).Script()
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
	sigScript1, err := txscript.NewScriptBuilder().AddData(sig1).AddData(senderPkData).Script()
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

func SendBrc20(ticker string, from, to btcutil.Address, amt, feeRate uint64, inscriberPrivateKey, senderPrivateKey *btcec.PrivateKey, config config.Config) (inscriptionId, hash string, err error) {
	res, err := InscribeTransfer(ticker, amt, from, inscriberPrivateKey, feeRate, config)
	if err != nil {
		return "", "", err
	}
	hashPtr, err := TransferBrc20(from, to, res.RevealTx, senderPrivateKey, feeRate, config)
	if err != nil {
		return inscriptionId, "", err
	}
	inscriptionId = res.RevealTx
	hash = *hashPtr
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
