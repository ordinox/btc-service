package config

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var (
	//go:embed default.yaml
	defaultConfigB []byte
	config         Config
)

type (
	Config struct {
		BtcConfig BtcConfig `mapstructure:"btc"`
		OpiConfig OpiConfig `mapstructure:"opi"`
	}

	BtcConfig struct {
		RpcHost    string `mapstructure:"rpc_host"`
		CookiePath string `mapstructure:"cookie_path"`
		WalletName string `mapstructure:"wallet_name"`
	}

	OpiConfig struct {
		Version   string       `mapstructure:"version"`
		Port      string       `mapstructure:"port"`
		Endpoints OpiEndpoints `mapstructure:"endpoints"`
	}

	OpiEndpoints struct {
		FetchEventsByInscriptionId string `mapstructure:"fetch_evts_by_inscription_id"`
	}
)

func init() {
	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)
	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(bytes.NewBuffer(defaultConfigB)); err != nil {
		log.Err(err).Msg("unable to load config")
	}
	if err := viper.Unmarshal(&config); err != nil {
		log.Err(err).Msg("unable to unmarshal config")
	}
}

func GetDefaultConfig() Config {
	return config
}

func (c BtcConfig) GetRpcHostWithWallet() string {
	return fmt.Sprintf("%s/wallet/%s", strings.TrimRight(c.RpcHost, "/"), c.WalletName)
}
