package btc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ordinox/btc-service/common"
	"github.com/ordinox/btc-service/config"
)

func GetUtxos(address string) (*common.WebUtxoResponse, error) {
	url := fmt.Sprintf("%s/get_utxos?address=%s", config.GetDefaultConfig().BtcConfig.ElectrumProxy, address)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	var webUtxoResponse common.WebUtxoResponse
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &webUtxoResponse)
	if err != nil {
		return nil, err
	}
	return &webUtxoResponse, nil
}
