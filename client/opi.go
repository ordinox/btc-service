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

	GetEventsByInscriptionIdResponse struct {
		Error  any          `json:"error"`
		Result []Brc20Event `json:"result"`
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

func (c OpiClient) GetEventsByInscriptionId(inscriptionId string) ([]Brc20Event, error) {
	endpoint := fmt.Sprintf("%s%s?inscription_id=%s", c.endpoint, c.config.Endpoints.FetchEventsByInscriptionId, inscriptionId)
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
	data := GetEventsByInscriptionIdResponse{}
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		log.Err(err).Msgf("error unmarshalling response: [resp = %s]", string(bodyBytes))
		return nil, err
	}

	return data.Result, nil
}
