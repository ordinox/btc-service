package client

import (
	"math/big"

	"github.com/ordinox/btc-service/config"
)

type (
	OpiClient struct {
		brc20Host string
		runesHost string
		config    config.OpiConfig
	}

	// HTTP API Responses
	Response[T any] struct {
		Error  any `json:"error"`
		Result T   `json:"result"`
	}

	// BRC20 Transfer Events
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

	// Runes Transfer Events
	RunesEvent struct {
		EventType  string      `json:"event_type"`
		Outpoint   interface{} `json:"outpoint"`
		Pkscript   interface{} `json:"pkscript"`
		WalletAddr interface{} `json:"wallet_addr"`
		RuneID     string      `json:"rune_id"`
		Amount     string      `json:"amount"`
	}

	// BRC20 Balance Response
	Brc20Balance struct {
		OverallBalance   string `json:"overall_balance"`
		AvailableBalance string `json:"available_balance"`
		BlockHeight      int    `json:"block_height"`
	}

	// Runes Balance Response Wrapper
	// Keeping this because we may need to verify blockheight
	RunesResponseData[T any] struct {
		Error         interface{} `json:"error"`
		Result        T           `json:"result"`
		DbBlockHeight int         `json:"db_block_height"`
	}

	// Runes Balance Response Wrapper
	RunesBalance struct {
		Pkscript     string `json:"pkscript"`
		WalletAddr   string `json:"wallet_addr"`
		RuneID       string `json:"rune_id"`
		RuneName     string `json:"rune_name"`
		TotalBalance string `json:"total_balance"`
	}

	RunesUnspentOutput struct {
		Pkscript   string     `json:"pkscript"`
		WalletAddr string     `json:"wallet_addr"`
		Outpoint   string     `json:"outpoint"`
		RuneIds    []string   `json:"rune_ids"`
		Balances   []*big.Int `json:"balances"`
	}
)
