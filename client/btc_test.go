package client

import (
	"fmt"
	"testing"

	"github.com/ordinox/btc-service/config"
)

func TestBtcClient(t *testing.T) {
	config.Init()
	fmt.Println("???", config.GetDefaultConfig().BtcConfig.SandshrewApiKey)
	client := NewBitcoinClient(config.GetDefaultConfig())
	info, err := client.GetBlockChainInfo()
	if err != nil {
		panic(err)
	}
	fmt.Println(info.Chain, info.Blocks)
	fmt.Println()
}
