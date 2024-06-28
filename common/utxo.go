package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/ordinox/btc-service/config"
)

func GetUtxos(address string, config config.BtcConfig) (*WebUtxoResponse, error) {
	// Use Sandshrew when possible
	if config.SandshrewApiKey != "" {
		return GetEsploraUtxos(address, config)
	}

	url := fmt.Sprintf("%s/getunspent?address=%s&network=%s", config.ElectrumProxy, address, config.GetChainConfigParams().Name)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	var webUtxoResponse WebUtxoResponse
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &webUtxoResponse)
	if err != nil {
		return nil, err
	}
	return &webUtxoResponse, nil
}

// Given an address and a min value, find an eligible UTXO for the address
func SelectOneUtxo(addr string, minValue uint64, config config.BtcConfig) (*btcjson.ListUnspentResult, error) {
	utxos, err := GetUtxos(addr, config)
	if err != nil {
		return nil, err
	}

	if len(utxos.Result) == 0 {
		return nil, fmt.Errorf("no utxos available for addr=%s", addr)
	}

	for _, utxo := range utxos.Result {
		if utxo.Value < minValue {
			continue
		}
		return &btcjson.ListUnspentResult{
			Amount:  float64(utxo.Value),
			TxID:    utxo.TxHash,
			Address: addr,
			Vout:    utxo.Vout,
		}, nil
	}
	return nil, fmt.Errorf("no eligible utxo found")
}
