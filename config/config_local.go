//go:build local_config
// +build local_config

package config

import (
	"bytes"

	_ "embed"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var (
	//go:embed local.yaml
	localConfigB []byte
)

func init() {
	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)
	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(bytes.NewBuffer(localConfigB)); err != nil {
		log.Err(err).Msg("unable to load local config")
		panic("no config found")
	}
	if err := viper.Unmarshal(&config); err != nil {
		log.Err(err).Msg("unable to unmarshal config")
	}
}
