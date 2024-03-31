package common

type WebUtxoResponse struct {
	ID      string    `json:"id"`
	Jsonrpc string    `json:"jsonrpc"`
	Result  []WebUtxo `json:"result"`
}

type WebUtxo struct {
	Height int    `json:"height"`
	TxHash string `json:"tx_hash"`
	TxPos  int    `json:"tx_pos"`
	Value  int    `json:"value"`
}
