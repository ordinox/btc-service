package runestone

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ordinox/btc-service/client"
	"github.com/ordinox/btc-service/common"
	"github.com/ordinox/btc-service/config"
	"github.com/ordinox/btc-service/runes"
)

type RuneTx struct {
	Edict   Edict
	Tx      *btcutil.Tx
	Senders []btcutil.Address
}

// Test rune balance
func TestEdicts(t *testing.T) {
	config.Init()
	config := config.GetDefaultConfig()

	// txId := "3d544a30daa2309fd6b8e90769e07b1fbc88883853b2ebc45ed31b183e0c956f"

	addr := "1JCX1jsCuiPsPj9nJWrZYCdYnauF99Z1mU"
	// txId := "1415e790442f0e6ed585d1b163620b1c8635b285ff9d9fad29cdc511995d9fb1"
	// client := client.NewBitcoinClient(config)

	utxos, err := common.GetUtxos(addr, config.BtcConfig)
	if err != nil {
		panic(err)
	}
	txIds := make([]string, 0)
	for _, u := range utxos.Result.ToUtxo() {
		txId := u.GetTxID()
		txIds = append(txIds, txId)
	}
	txs, err := client.GetRawTxs(config, txIds)
	if err != nil {
		panic(err)
	}
	runes := make([]RuneTx, 0)
	for _, tx := range txs {
		artifact := DecipherRunestone(tx.MsgTx())
		if artifact.Runestone != nil {
			for _, e := range artifact.Runestone.Edicts {
				output := tx.MsgTx().TxOut[e.Output]
				pkScript, _ := txscript.ParsePkScript(output.PkScript)
				_ = pkScript
				runes = append(runes, RuneTx{Edict: e, Tx: tx})
			}
		}
	}
	hashes := make([]string, 0)
	for _, rune := range runes {
		for _, in := range rune.Tx.MsgTx().TxIn {
			hashes = append(hashes, in.PreviousOutPoint.Hash.String())
		}
	}
	txs, err = client.GetRawTxs(config, hashes)
	if err != nil {
		panic(err)
	}
	txMap := make(map[string]*btcutil.Tx)
	for _, tx := range txs {
		txMap[tx.Hash().String()] = tx
	}
	for i, rune := range runes {
		for _, in := range rune.Tx.MsgTx().TxIn {
			tx := txMap[in.PreviousOutPoint.Hash.String()]
			out := tx.MsgTx().TxOut[in.PreviousOutPoint.Index]
			_, addrs, _, err := txscript.ExtractPkScriptAddrs(out.PkScript, &chaincfg.MainNetParams)
			if err != nil {
				panic(err)
			}
			rune.Senders = addrs
			runes[i] = rune
		}
	}
	for _, rune := range runes {
		fmt.Printf("%s - %s - %v\n", rune.Edict.Id, rune.Edict.Amount, rune.Senders)
	}
}

func TestEncipherRune(t *testing.T) {
	script1, _ := runes.CreateTransferScript(runes.Rune{BlockNumber: 100, TxIndex: 100}, big.NewInt(100), 0, true)
	edict := Edict{Id: NewRuneId(100, 100), Output: 0, Amount: big.NewInt(100)}
	runeStone := Runestone{Edicts: []Edict{edict}}
	runeScript, _ := EncipherRunestone(runeStone).Script()
	fmt.Println((script1), (runeScript))
	tx := wire.NewMsgTx(wire.TxVersion)
	tx.AddTxOut(wire.NewTxOut(0, script1))
	for _, i := range DecipherRunestone(tx).Runestone.Edicts {
		fmt.Printf("%v", i)
	}
}
