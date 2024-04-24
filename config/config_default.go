//go:build !local_config
// +build !local_config

package config

var (
	config Config
)

func Init() {
	config = Config{
		BtcConfig: BtcConfig{
			RpcHost:        "localhost:18443",
			CookiePath:     "/home/ubuntu/.bitcoin/regtest",
			WalletName:     "w1",
			OrdPath:        "/home/ubuntu/OPI/ord/target/release",
			BitcoinDataDir: "/home/ubuntu/.bitcoin",
			OrdDataDir:     "/home/ubuntu/OPI/ord/target/release",
			ElectrumProxy:  "http://localhost:6789",
		},
		OpiConfig: OpiConfig{
			Version:  "0.3.0",
			Brc20Url: "http://localhost:8000",
			Endpoints: OpiEndpoints{
				FetchEventsByInscriptionId: "/v1/brc20/event",
				FetchBrc20Balance:          "/v1/brc20/get_current_balance_of_wallet",
			},
		},
	}
}
