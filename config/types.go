package config

type (
	Config struct {
		BtcConfig BtcConfig `mapstructure:"btc"`
		OpiConfig OpiConfig `mapstructure:"opi"`
	}

	BtcConfig struct {
		RpcHost        string `mapstructure:"rpc_host"`
		CookiePath     string `mapstructure:"cookie_path"`
		WalletName     string `mapstructure:"wallet_name"`
		ChainConfig    string `mapstructure:"chain_cfg"`
		OrdPath        string `mapstructure:"ord_path"`
		BitcoinDataDir string `mapstructure:"bitcoin_data_dir"`
		OrdDataDir     string `mapstructure:"ord_data_dir"`
		ElectrumProxy  string `mapstructure:"electrum_proxy"`
	}

	OpiConfig struct {
		Version   string       `mapstructure:"version"`
		Brc20Url  string       `mapstructure:"brc20_url"`
		Endpoints OpiEndpoints `mapstructure:"endpoints"`
		RunesUrl  string       `mapstructure:"runes_url"`
	}

	OpiEndpoints struct {
		FetchEventsByInscriptionId      string `mapstructure:"fetch_brc20_evts_by_inscription_id"`
		FetchRunesEventsByTransactionId string `mapstructure:"fetch_runes_evts_by_txid"`
		FetchBrc20Balance               string `mapstructure:"fetch_brc20_balance"`
		FetchRunesBalance               string `mapstructure:"fetch_runes_balance"`
		FetchRunesUnspentOutpoint       string `mapstructure:"fetch_runes_unspent_outpoints"` // These names have techincal meaning and are not to be changed
	}
)
