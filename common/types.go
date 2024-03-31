package common

import "github.com/btcsuite/btcd/btcjson"

type Utxo interface {
	GetTxID() string
	GetVout() uint32
	GetValueInSats() uint64
}

type WebUtxoResponse struct {
	ID      string    `json:"id"`
	Jsonrpc string    `json:"jsonrpc"`
	Result  []WebUtxo `json:"result"`
}

type WebUtxo struct {
	Height int    `json:"height"`
	TxHash string `json:"tx_hash"`
	Vout   uint32 `json:"tx_pos"`
	Value  uint64 `json:"value"`
}

func (w WebUtxo) GetTxID() string {
	return w.TxHash
}

func (w WebUtxo) GetVout() uint32 {
	return w.Vout
}

func (w WebUtxo) GetValueInSats() uint64 {
	return w.Value
}

type BtcUnspent struct {
	btcjson.ListUnspentResult
}

func (w BtcUnspent) GetTxID() string {
	return w.TxID
}

func (w BtcUnspent) GetVout() uint32 {
	return w.Vout
}

func (w BtcUnspent) GetValueInSats() uint64 {
	return uint64(w.Amount * 100000000)
}

func NewBtcUnspent(b btcjson.ListUnspentResult) BtcUnspent {
	return BtcUnspent{b}
}
