package common

import "github.com/btcsuite/btcd/btcjson"

type Utxo interface {
	GetTxID() string
	GetVout() uint32
	GetValueInSats() uint64
}

type WebUtxoResponse struct {
	Jsonrpc string   `json:"jsonrpc"`
	Result  WebUtxos `json:"result"`
}

type WebUtxo struct {
	Height int    `json:"height"`
	TxHash string `json:"tx_hash"`
	Vout   uint32 `json:"tx_pos"`
	Value  uint64 `json:"value"`
}

type WebUtxos []WebUtxo

func (w WebUtxos) ToUtxo() []Utxo {
	var utxos []Utxo = make([]Utxo, len(w))
	for i := range w {
		utxos[i] = w[i]
	}
	return utxos
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

type BtcUnspents []BtcUnspent

func (w BtcUnspent) GetTxID() string {
	return w.TxID
}

func (w BtcUnspent) GetVout() uint32 {
	return w.Vout
}

func (w BtcUnspent) GetValueInSats() uint64 {
	return uint64(w.Amount * 100000000)
}

func (w BtcUnspents) ToUtxo() []Utxo {
	var utxos []Utxo = make([]Utxo, len(w))
	for i := range w {
		utxos[i] = w[i]
	}
	return utxos
}

func NewBtcUnspent(b btcjson.ListUnspentResult) BtcUnspent {
	return BtcUnspent{b}
}
