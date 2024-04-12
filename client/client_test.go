package client

import (
	"fmt"
	"testing"

	"github.com/ordinox/btc-service/config"
)

var _ = fmt.Println

func TestOpiClient(t *testing.T) {
	opiClient := NewOpiClient(config.GetDefaultConfig().OpiConfig)
	evts, err := opiClient.GetEventsByInscriptionId("5ccae5f073c372698108f4fc65b56d8f82d0857828a8987d1e6f8b10052f05ffi0")
	if err != nil {
		panic(err)
	}
	for _, evt := range evts {
		fmt.Println(evt.Tick, evt.EventType, evt.SourceWallet, evt.SpentWallet)
	}
}
