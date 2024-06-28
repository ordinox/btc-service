package client

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/ordinox/btc-service/config"
)

type SandshrewMulticallParams struct {
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
	Jsonrpc string        `json:"jsonrpc"`
}

type GetRawTxsResponse struct {
	Result []struct {
		Result string `json:"result"`
	} `json:"result"`
}

func GetRawTxs(config config.Config, txHashes []string) ([]*btcutil.Tx, error) {
	url := "https://mainnet.sandshrew.io/v1/" + config.BtcConfig.SandshrewApiKey
	method := "POST"

	params := make([]interface{}, len(txHashes))
	for i, txHash := range txHashes {
		params[i] = []interface{}{"btc_getrawtransaction", []any{txHash}}
	}

	payload := SandshrewMulticallParams{
		Method:  "sandshrew_multicall",
		Params:  params,
		ID:      0,
		Jsonrpc: "2.0",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	data := new(GetRawTxsResponse)
	if err := json.Unmarshal(body, data); err != nil {
		panic(err)
	}
	retData := make([]*btcutil.Tx, 0)
	for _, i := range data.Result {
		txB, err := hex.DecodeString(i.Result)
		if err != nil {
			panic(err)
		}
		tx, err := btcutil.NewTxFromBytes(txB)
		if err != nil {
			panic(err)
		}
		retData = append(retData, tx)
	}
	return retData, nil
}

func GetTxInfo(txId string) {

}
