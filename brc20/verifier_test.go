package brc20

import (
	"testing"

	"github.com/ordinox/btc-service/config"
	"github.com/stretchr/testify/require"
)

func TestBis(t *testing.T) {
	config.Init()
	config := config.GetDefaultConfig()
	config.BtcConfig.ChainConfig = "mainnet"
	err := VerifyBrc20TransferV2(config, VerifyBrc20DepositData{
		FromWalletAddr: "1JCX1jsCuiPsPj9nJWrZYCdYnauF99Z1mU",
		ToWalletAddr:   "bc1pwfk592dwzz4t9wwp7yxs74tctctl9v2aez07j5lkhzx938hvur6sakvwvf",
		TxId:           "af6a79c5c9d124bdeca189dd8375c6e17f2f24258a485e5515faec0ba9070180",
		Amount:         1,
		Decimals:       18,
		Tick:           "bzrk",
	})
	require.NoError(t, err)
}
