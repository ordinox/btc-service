package cmd

import (
	"github.com/ordinox/btc-service/config"
	"github.com/spf13/cobra"
)

func Execute() {
	root := cobra.Command{
		Use:   "btc-service",
		Short: "cli for interacting with tokens on bitcoin",
	}

	config.Init()

	config := config.GetDefaultConfig()

	root.AddCommand(
		brc20Cmd(config),
		getKeyPairCmd(),
		genBlocksCmd(config.BtcConfig),
		getUtxosCmd(),
		transferBtcCmd(config),
		satsToBtcCmd(),
		runesCmd(config),
	)
	err := root.Execute()
	if err != nil {
		panic(err)
	}
}
