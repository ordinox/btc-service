package client

import (
	"strings"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/ordinox/btc-service/config"
	"github.com/rs/zerolog/log"
)

type BtcRpcClient struct {
	*rpcclient.Client
	TrackedAddreses map[string]bool
}

func NewBitcoinClient(config config.Config) *BtcRpcClient {
	connConfig := &rpcclient.ConnConfig{
		Host:         config.BtcConfig.GetRpcHostWithWallet(),
		HTTPPostMode: true, // Bitcoin Core
		DisableTLS:   true,
		CookiePath:   config.BtcConfig.CookiePath,
	}
	client, err := rpcclient.New(connConfig, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create btc rpc client")
	}
	if _, err := client.LoadWallet(config.BtcConfig.WalletName); err != nil {
		if !strings.Contains(err.Error(), "already loaded") {
			log.Fatal().Err(err).Msgf("error loading wallet: %s", config.BtcConfig.WalletName)
		}
	}
	return &BtcRpcClient{
		Client:          client,
		TrackedAddreses: make(map[string]bool),
	}
}
