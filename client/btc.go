package client

import (
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/ordinox/btc-service/common"
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

func (r *BtcRpcClient) ImportAddress(address string) error {
	log.Info().Msgf("importing new addr: %s", address)
	r.TrackedAddreses[strings.ToLower(address)] = true
	err := r.Client.ImportAddress(address)
	if err != nil {
		log.Err(err).Msg("error importing address")
	}
	return err
}

func GetWebUtxo(addr string) []common.WebUtxo {
	return nil
}

func (r *BtcRpcClient) GetUtxos(addr btcutil.Address) (common.BtcUnspents, error) {
	if r.TrackedAddreses == nil {
		r.TrackedAddreses = make(map[string]bool)
	}
	if _, ok := r.TrackedAddreses[strings.ToLower(addr.String())]; !ok {
		if err := r.ImportAddress(addr.String()); err != nil {
			return nil, err
		}
	}
	unspent, err := r.ListUnspent()
	if err != nil {
		return nil, err
	}
	utxos := make([]common.BtcUnspent, 0)
	for _, utxo := range unspent {
		if utxo.Address == addr.String() {
			utxos = append(utxos, common.NewBtcUnspent(utxo))
		}
	}
	return utxos, nil
}
