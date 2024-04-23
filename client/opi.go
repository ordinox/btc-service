package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ordinox/btc-service/config"
	"github.com/rs/zerolog/log"
)

func (r RunesBalance) Copy() *RunesBalance {
	return &RunesBalance{
		Pkscript:     r.Pkscript,
		WalletAddr:   r.WalletAddr,
		RuneID:       r.RuneID,
		RuneName:     r.RuneName,
		TotalBalance: r.TotalBalance,
	}
}

// Create a new OPI client and check if the API is live
func NewOpiClient(c config.OpiConfig) *OpiClient {
	if len(c.Brc20Port) == 0 {
		panic("OPI_CONFIG_ERROR: BRC20 PORT UNDEFINED")
	}
	if len(c.RunesPort) == 0 {
		panic("OPI_CONFIG_ERROR: RUNES PORT UNDEFINED")
	}

	brc20Endpoint := fmt.Sprintf("http://localhost:%s", c.Brc20Port)
	resp1, err := http.Get(brc20Endpoint + "/v1/brc20/ip")
	if err != nil {
		log.Fatal().Err(err).Msgf("error connecting brc20 to opi endpoint [%s]", brc20Endpoint)
	}
	defer resp1.Body.Close()

	if resp1.StatusCode != http.StatusOK {
		log.Fatal().Msgf("status error connecting to brc20 opi endpoint [%s]. status not 200, but [%d]", brc20Endpoint, resp1.StatusCode)
	}

	runesEndpoint := fmt.Sprintf("http://localhost:%s", c.RunesPort)
	resp2, err := http.Get(runesEndpoint + "/v1/runes/ip")
	if err != nil {
		log.Fatal().Err(err).Msgf("error connecting runes to opi endpoint [%s]", brc20Endpoint)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		log.Fatal().Msgf("status error connecting to runes opi endpoint [%s]. status not 200", brc20Endpoint)
	}

	return &OpiClient{
		brc20Host: brc20Endpoint,
		runesHost: runesEndpoint,
		config:    c,
	}
}

// HTTP GET request wrapper
func getRequest(endpoint string) ([]byte, error) {
	resp, err := http.Get(endpoint)
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

// Get all events by Inscription ID
// Events prove if a BRC20 was transferred or not
func (c OpiClient) GetEventsByInscriptionId(inscriptionId string) ([]Brc20Event, error) {
	endpoint := fmt.Sprintf("%s%s?inscription_id=%s", c.brc20Host, c.config.Endpoints.FetchEventsByInscriptionId, inscriptionId)
	data := Response[[]Brc20Event]{}
	bodyBytes, err := getRequest(endpoint)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		log.Err(err).Msgf("error unmarshalling response: [resp = %s]", string(bodyBytes))
		return nil, err
	}

	return data.Result, nil
}

// Get BRC20 balance
func (c OpiClient) GetBrc20Balance(address, ticker string) (*Brc20Balance, error) {
	endpoint := fmt.Sprintf("%s%s?address=%s&ticker=%s", c.brc20Host, c.config.Endpoints.FetchBrc20Balance, address, ticker)
	bodyBytes, err := getRequest(endpoint)
	data := Response[Brc20Balance]{}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		log.Err(err).Msgf("error unmarshalling response: [resp = %s]", string(bodyBytes))
		return nil, err
	}
	return &data.Result, nil
}

// Get Runes Balance
func (c OpiClient) GetRunesBalance(address string) ([]RunesBalance, error) {
	endpoint := fmt.Sprintf("%s%s?address=%s", c.runesHost, c.config.Endpoints.FetchRunesBalance, address)
	bodyBytes, err := getRequest(endpoint)
	data := Response[[]RunesBalance]{}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		log.Err(err).Msgf("error unmarshalling response: [resp = %s]", string(bodyBytes))
		return nil, err
	}
	return data.Result, nil
}

// Get Runes Events by TxID
func (c OpiClient) GetRunesEventsByTxID(txId string) ([]RunesEvent, error) {
	endpoint := fmt.Sprintf("%s%s?transaction_id=%s", c.runesHost, c.config.Endpoints.FetchRunesEventsByTransactionId, txId)
	data := Response[[]RunesEvent]{}
	bodyBytes, err := getRequest(endpoint)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		log.Err(err).Msgf("error unmarshalling response: [resp = %s]", string(bodyBytes))
		return nil, err
	}

	return data.Result, nil
}

// Get all Rune UTXOs
func (c OpiClient) GetRunesUnspentOutpoints(address string) ([]RunesUnspentOutput, error) {
	endpoint := fmt.Sprintf("%s%s?address=%s", c.runesHost, c.config.Endpoints.FetchRunesUnspentOutpoint, address)
	data := RunesResponseData[[]RunesUnspentOutput]{}
	bodyBytes, err := getRequest(endpoint)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		log.Err(err).Msgf("error unmarshalling response: [resp = %s]", string(bodyBytes))
		return nil, err
	}
	return data.Result, nil
}
