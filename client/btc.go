package client

import (
	"fmt"
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
	host := config.BtcConfig.GetRpcHostWithWallet()
	cookiePath := config.BtcConfig.CookiePath
	user := ""
	pass := ""
	if config.BtcConfig.SandshrewApiKey != "" {
		host = fmt.Sprintf("%s/%s", "mainnet.sandshrew.io/v1", config.BtcConfig.SandshrewApiKey)
		cookiePath = ""
		user = "user"
		pass = "pass"
	}

	connConfig := &rpcclient.ConnConfig{
		Host:         host,
		HTTPPostMode: true, // Bitcoin Core
		DisableTLS:   true,
		CookiePath:   cookiePath,
		User:         user,
		Pass:         pass,
	}
	client, err := rpcclient.New(connConfig, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create btc rpc client")
	}
	if config.BtcConfig.SandshrewApiKey == "" {
		if _, err := client.LoadWallet(config.BtcConfig.WalletName); err != nil {
			if !strings.Contains(err.Error(), "already loaded") {
				log.Fatal().Err(err).Msgf("error loading wallet: %s", config.BtcConfig.WalletName)
			}
		}
	}
	return &BtcRpcClient{
		Client:          client,
		TrackedAddreses: make(map[string]bool),
	}
}
