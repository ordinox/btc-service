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
	}

	OpiConfig struct {
		Version   string       `mapstructure:"version"`
		Port      string       `mapstructure:"port"`
		Endpoints OpiEndpoints `mapstructure:"endpoints"`
	}

	OpiEndpoints struct {
		FetchEventsByInscriptionId string `mapstructure:"fetch_evts_by_inscription_id"`
		FetchBalance               string `mapstructure:"fetch_balance"`
	}
)
