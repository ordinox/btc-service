package runes

import (
	"fmt"

	"github.com/ordinox/btc-service/client"
	"github.com/ordinox/btc-service/config"
)

// Since "/event" api returns a UTXO and UTXOs can have multiple outputs, it's impossible to extract one receiver
// Hence, this verifier requires sender, receiver & amountToVerify to ensure that the runes transfer has actually taken place
// TODO: Check if the blockheight of txid and the blockheight of the scan
func VerifyRunesDeposit(config config.Config, txId string, sender, receiver, amountToVerify string) error {
	var (
		client           = client.NewOpiClient(config.OpiConfig)
		senderVerified   = false
		receiverVerified = false
	)

	events, err := client.GetRunesEventsByTxID(txId)
	if err != nil {
		return err
	}

	// Loop through all the outputs and check if Sender, Receiver & Amount matches the inputs
	for _, evt := range events {
		if evt.EventType == "input" && evt.WalletAddr != nil {
			senderVerified = true
		} else if evt.EventType == "output" && evt.WalletAddr != nil && evt.Amount == amountToVerify && evt.WalletAddr.(string) == receiver {
			receiverVerified = true
		}
		if senderVerified && receiverVerified {
			return nil
		}
	}

	if !senderVerified {
		return fmt.Errorf("[Runes Verifier] Could not verify sender")
	}

	if !receiverVerified {
		return fmt.Errorf("[Runes Verifier] Could not verify receiver")
	}

	return nil
}
