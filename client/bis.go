package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ordinox/btc-service/config"
	"github.com/rs/zerolog/log"
)

type BISClient struct {
	baseUrl string
	apiKey  string
}

var _ RunesUnspentOutput = BISRunesUnspentOutput{}

func (u BISRunesUnspentOutput) GetPkScript() string {
	return u.Pkscript
}

func (u BISRunesUnspentOutput) GetWalletAddr() string {
	return u.WalletAddr
}

func (u BISRunesUnspentOutput) GetOutpoint() string {
	return u.Output
}

func (u BISRunesUnspentOutput) GetRuneIds() []string {
	return u.RuneIds
}

func (u BISRunesUnspentOutput) GetRuneNames() []string {
	return u.RuneNames
}

type BISRunesUnspentOutputList = []BISRunesUnspentOutput

func NewBISClient(c config.BISConfig) *BISClient {
	if len(c.APIKey) == 0 {
		panic("BIS API_KEY is not defined")
	}
	return &BISClient{"https://api.bestinslot.xyz", c.APIKey}
}

func authenticatedBisGetRequest(endpoint, headerKey, headerValue string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Err(err).Msgf("error creating request [url = %s]", endpoint)
		return nil, err
	}

	req.Header.Set(headerKey, headerValue)

	resp, err := client.Do(req)
	if err != nil {
		log.Err(err).Msgf("error fetching events by inscription id [url = %s]", endpoint)
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Err(err).Msg("error reading msg body")
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Info().Msgf("GET %s", string(endpoint))
		log.Err(err).Msgf("http status not ok: [resp = %s]", string(bodyBytes))
		return nil, err
	}
	return bodyBytes, nil
}

// Get all brc20 events with the given transaction ID
func (b BISClient) GetEventsByTransactionId(txId string) ([]BISBrc20Event, error) {
	endpoint := "https://api.bestinslot.xyz/v3/brc20/event_from_txid?=" + txId
	res, err := authenticatedBisGetRequest(endpoint, "x-api-key", b.apiKey)
	if err != nil {
		return nil, err
	}
	bisResponse := BISResponseWrapper[[]BISBrc20Event]{}
	if err := json.Unmarshal(res, &bisResponse); err != nil {
		log.Err(err).Msgf("error unmarshalling response: [resp = %s]", string(res))
		return nil, err
	}
	return bisResponse.Data, nil
}

// Fetch runes UTXOs from BIS API
func (b BISClient) FetchRunesUtxos(address string) ([]RunesUnspentOutput, error) {
	endpoint := fmt.Sprintf("%s/v3/runes/wallet_valid_outputs?address=%s&order=asc&offset=0&count=2000&sort_by=output", b.baseUrl, address)
	res, err := authenticatedBisGetRequest(endpoint, "x-api-key", b.apiKey)
	if err != nil {
		return nil, err
	}

	bisResponse := BISResponseWrapper[[]BISRunesUnspentOutput]{}
	if err := json.Unmarshal(res, &bisResponse); err != nil {
		log.Err(err).Msgf("error unmarshalling response: [resp = %s]", string(res))
		return nil, err
	}

	data := make([]RunesUnspentOutput, 0)
	for _, i := range bisResponse.Data {
		data = append(data, i)
	}

	return data, nil
}
