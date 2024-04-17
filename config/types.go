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
		Brc20Port string       `mapstructure:"brc20_port"`
		Endpoints OpiEndpoints `mapstructure:"endpoints"`
		RunesPort string       `mapstructure:"runes_port"`
	}

	OpiEndpoints struct {
		FetchEventsByInscriptionId      string `mapstructure:"fetch_brc20_evts_by_inscription_id"`
		FetchRunesEventsByTransactionId string `mapstructure:"fetch_runes_evts_by_txid"`
		FetchBrc20Balance               string `mapstructure:"fetch_brc20_balance"`
		FetchRunesBalance               string `mapstructure:"fetch_runes_balance"`
	}
)
