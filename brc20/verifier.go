package brc20

import (
	"errors"
	"github.com/ordinox/btc-service/client"
	"github.com/ordinox/btc-service/config"
)

// Return the sender and reciever of a BRC20 transfer
func VerifyBrc20Deposit(config config.Config, inscriptionId string) (string, string, error) {
	// TODO: Add amount checking
	client := client.NewOpiClient(config.OpiConfig)
	events, err := client.GetEventsByInscriptionId(inscriptionId)
	if err != nil {
		return "", "", err
	}
	for _, evt := range events {
		if evt.EventType == "transfer-transfer" {
			return evt.SourceWallet, evt.SpentWallet, nil
		}
	}
	return "", "", errors.New("unable to fetch transfer-inscription transfer")
}
