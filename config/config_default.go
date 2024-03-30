//go:build !local_config
// +build !local_config

package config

var (
	config Config
)

func init() {
	config = Config{
		BtcConfig: BtcConfig{
			RpcHost:        "localhost:18443",
			CookiePath:     "/home/ubuntu/.bitcoin/regtest",
			WalletName:     "legacy",
			OrdPath:        "/home/ubuntu/OPI/ord/target/release",
			BitcoinDataDir: "/home/ubuntu/.bitcoin",
			OrdDataDir:     "/home/ubuntu/OPI/ord/target/release",
		},
		OpiConfig: OpiConfig{
			Version: "0.3.0",
			Port:    "8000",
			Endpoints: OpiEndpoints{
				FetchEventsByInscriptionId: "/v1/brc20/event",
				FetchBalance:               "/v1/brc20/get_current_balance_of_wallet",
			},
		},
	}
}
