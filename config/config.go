package config

import (
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
)

func GetDefaultConfig() Config {
	return config
}

func (c BtcConfig) GetRpcHostWithWallet() string {
	return fmt.Sprintf("%s/wallet/%s", strings.TrimRight(c.RpcHost, "/"), c.WalletName)
}

func (c BtcConfig) GetChainConfigParams() *chaincfg.Params {
	if c.ChainConfig == "mainnet" {
		return &chaincfg.MainNetParams
	} else if c.ChainConfig == "testnet" {
		return &chaincfg.SigNetParams
	} else {
		return &chaincfg.RegressionNetParams
	}
}

func (c BtcConfig) GetOrdChainConfigFlag() string {
	if c.GetChainConfigParams().Name == chaincfg.MainNetParams.Name {
		return ""
	} else if c.GetChainConfigParams().Name == chaincfg.SigNetParams.Name {
		return "-t"
	} else {
		return "-r"
	}
}
