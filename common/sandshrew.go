package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ordinox/btc-service/config"
)

type EsploraResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  []struct {
		Txid   string `json:"txid"`
		Vout   uint32 `json:"vout"`
		Status struct {
			Confirmed   bool   `json:"confirmed"`
			BlockHeight int    `json:"block_height"`
			BlockHash   string `json:"block_hash"`
			BlockTime   uint32 `json:"block_time"`
		} `json:"status"`
		Value uint64 `json:"value"`
	} `json:"result"`
}

func GetEsploraUtxos(address string, config config.BtcConfig) (*WebUtxoResponse, error) {
	url := "https://mainnet.sandshrew.io/v1/" + config.SandshrewApiKey
	method := "POST"

	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "esplora_address::utxo",
		"params":  []string{address},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling payload:", err)
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	var data EsploraResponse
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Error parsing esplora response", err)
		return nil, err
	}

	webUtxoResponse := WebUtxoResponse{
		Jsonrpc: data.Jsonrpc,
		Result:  make(WebUtxos, len(data.Result)),
	}
	for i, utxo := range data.Result {
		wUtxo := WebUtxo{
			Height: utxo.Status.BlockHeight,
			TxHash: utxo.Txid,
			Vout:   uint32(utxo.Vout),
			Value:  utxo.Value,
		}
		webUtxoResponse.Result[i] = wUtxo
	}
	return &webUtxoResponse, nil
}
