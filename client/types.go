package client

import (
	"math/big"
	"time"

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

	// OPI Runes Unspent Output
	OPIRunesUnspentOutput struct {
		Pkscript   string     `json:"pkscript"`
		WalletAddr string     `json:"wallet_addr"`
		Outpoint   string     `json:"outpoint"`
		RuneIds    []string   `json:"rune_ids"`
		Balances   []*big.Int `json:"balances"`
	}

	// BestInSlot Runes Unspent Output
	BISRunesUnspentOutput struct {
		Pkscript        string   `json:"pkscript"`
		WalletAddr      string   `json:"wallet_addr"`
		Output          string   `json:"output"`
		RuneIds         []string `json:"rune_ids"`
		Balances        []int64  `json:"balances"`
		RuneNames       []string `json:"rune_names"`
		SpacedRuneNames []string `json:"spaced_rune_names"`
		Decimals        []int    `json:"decimals"`
	}

	BISResponseWrapper[T any] struct {
		Data        T   `json:"data"`
		BlockHeight int `json:"block_height"`
	}

	BISBrc20Event struct {
		InscriptionID string `json:"inscription_id"`
		EventType     string `json:"event_type"`
		Event         struct {
			Tick            string `json:"tick"`
			Amount          string `json:"amount"`
			UsingTxID       string `json:"using_tx_id"`
			SpentWallet     string `json:"spent_wallet"`
			SourceWallet    string `json:"source_wallet"`
			SpentPkScript   string `json:"spent_pkScript"`
			SourcePkScript  string `json:"source_pkScript"`
			Price           int    `json:"price"`
			MarketplaceType string `json:"marketplace_type"`
		} `json:"event"`
	}

	BISRuneEvent struct {
		EventType      string    `json:"event_type"`
		Txid           string    `json:"txid"`
		Outpoint       string    `json:"outpoint"`
		Pkscript       string    `json:"pkscript"`
		WalletAddr     string    `json:"wallet_addr"`
		RuneID         string    `json:"rune_id"`
		Amount         string    `json:"amount"`
		BlockHeight    int       `json:"block_height"`
		BlockTimestamp time.Time `json:"block_timestamp"`
		RuneName       string    `json:"rune_name"`
		SpacedRuneName string    `json:"spaced_rune_name"`
		Decimals       int       `json:"decimals"`
		SaleInfo       struct {
			SalePrice        int    `json:"sale_price"`
			SoldToPkscript   string `json:"sold_to_pkscript"`
			SoldToWalletAddr string `json:"sold_to_wallet_addr"`
			Marketplace      string `json:"marketplace"`
		} `json:"sale_info"`
	}
)

type RunesUnspentOutput interface {
	GetPkScript() string
	GetWalletAddr() string
	GetOutpoint() string
	GetRuneIds() []string
	GetRuneNames() []string
}
