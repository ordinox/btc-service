package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ordinox/btc-service/config"
	"github.com/rs/zerolog/log"
)

type (
	OpiClient struct {
		endpoint string
		config   config.OpiConfig
	}

	Response[T any] struct {
		Error  any `json:"error"`
		Result T   `json:"result"`
	}

	Brc20Event struct {
		Tick           string `json:"tick"`
		Amount         string `json:"amount"`
		SourceWallet   string `json:"source_wallet"`
		SourcePkScript string `json:"source_pkScript"`
		EventType      string `json:"event_type"`
		UsingTxID      string `json:"using_tx_id,omitempty"`
		SpentWallet    string `json:"spent_wallet,omitempty"`
		SpentPkScript  string `json:"spent_pkScript,omitempty"`
	}

	Brc20Balance struct {
		OverallBalance   string `json:"overall_balance"`
		AvailableBalance string `json:"available_balance"`
		BlockHeight      int    `json:"block_height"`
	}
)

func NewOpiClient(c config.OpiConfig) *OpiClient {
	endpoint := fmt.Sprintf("http://localhost:%s", c.Port)
	resp, err := http.Get(endpoint + "/v1/brc20/ip")
	if err != nil {
		log.Fatal().Err(err).Msgf("error connecting to opi endpoint [%s]", endpoint)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal().Msgf("status error connecting to opi endpoint [%s]. status not 200", endpoint)
	}

	return &OpiClient{
		endpoint: fmt.Sprintf("http://localhost:%s", c.Port),
		config:   c,
	}
}

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
		log.Err(err).Msgf("http status not ok: [resp = %s]", string(bodyBytes))
		return nil, err
	}
	return bodyBytes, nil
}

func (c OpiClient) GetEventsByInscriptionId(inscriptionId string) ([]Brc20Event, error) {
	endpoint := fmt.Sprintf("%s%s?inscription_id=%s", c.endpoint, c.config.Endpoints.FetchEventsByInscriptionId, inscriptionId)
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

func (c OpiClient) GetBalance(address, ticker string) (*Brc20Balance, error) {
	endpoint := fmt.Sprintf("%s%s?address=%s&ticker=%s", c.endpoint, c.config.Endpoints.FetchBalance, address, ticker)
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
