package client

import (
	"fmt"
	"testing"

	"github.com/ordinox/btc-service/common"
	"github.com/ordinox/btc-service/config"
	"github.com/stretchr/testify/require"
)

func TestBisClient(t *testing.T) {
	config.Init()
	config := config.GetDefaultConfig()
	client := NewBISClient(config.BISConfig)
	events, err := client.GetEventsByTransactionId("e20ac63402f36e9eba2e2a27e3699e65ca2998e319e2d4de69f235efd032ff0a")
	require.NoError(t, err)
	fmt.Println(events.BlockHeight)

	for _, e := range events.Data {
		fmt.Println(e.EventType)
		fmt.Println(e.Event.Amount)
		fmt.Println(common.ParseStringFloat64(e.Event.Amount, 18))
		fmt.Println(e.Event.Tick)
		fmt.Println(e.Event.SourceWallet)
		fmt.Println(e.Event.SpentWallet)
	}
}
