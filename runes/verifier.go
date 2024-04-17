package runes

import (
	"fmt"

	"github.com/ordinox/btc-service/client"
	"github.com/ordinox/btc-service/config"
)

// Return the sender and the receiver of a Runes transfer
func VerifyRunesDeposit(config config.Config, txId string) (sender, receiver string, err error) {
	client := client.NewOpiClient(config.OpiConfig)
	events, err := client.GetRunesEventsByTxID(txId)
	if err != nil {
		return "", "", err
	}
	for _, evt := range events {
		if evt.EventType == "input" && evt.WalletAddr != nil {
			sender = evt.WalletAddr.(string)
		} else if evt.EventType == "output" && evt.WalletAddr != nil {
			receiver = evt.WalletAddr.(string)
		}
	}
	if len(sender) == 0 || len(receiver) == 0 {
		return "", "", fmt.Errorf("invalid runes transfer txid")
	}
	return
}
