package brc20

import (
	"errors"
	"strings"

	"github.com/ordinox/btc-service/client"
	"github.com/ordinox/btc-service/common"
	"github.com/ordinox/btc-service/config"
)

var (
	ErrNoEventsFound        = errors.New("no events found")
	ErrInvalidBrc20Transfer = errors.New("error verifying brc20 transfer")
)

// Return the sender and reciever of a BRC20 transfer
func VerifyBrc20Deposit(config config.Config, inscriptionId string) (string, string, error) {
	// Regtest verification
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

type VerifyBrc20DepositData struct {
	InscriptionId, TxId, Tick    string
	FromWalletAddr, ToWalletAddr string
	Amount                       float64
	Decimals                     int // Decimals of the amount we are checking
}

func VerifyBrc20TransferV2(config config.Config, data VerifyBrc20DepositData) error {
	if config.BtcConfig.ChainConfig == "mainnet" {

		// Use BIS
		client := client.NewBISClient(config.BISConfig)
		events, err := client.GetEventsByTransactionId(data.TxId)
		if err != nil {
			return err
		}
		if len(events.Data) == 0 {
			return ErrNoEventsFound
		}
		for _, evt := range events.Data {
			if evt.EventType != "transfer-transfer" {
				continue
			}
			if !strings.EqualFold(evt.Event.SourceWallet, data.FromWalletAddr) {
				continue
			}
			if !strings.EqualFold(evt.Event.SpentWallet, data.ToWalletAddr) {
				continue
			}
			if !strings.EqualFold(evt.Event.Tick, data.Tick) {
				continue
			}
			amtParsed, err := common.ParseStringFloat64(evt.Event.Amount, data.Decimals)
			if err != nil {
				return err
			}
			if amtParsed != data.Amount {
				continue
			}
			// If control reaches here, then an event passed the check
			return nil
		}

		return ErrInvalidBrc20Transfer
	} else {
		client := client.NewOpiClient(config.OpiConfig)
		events, err := client.GetEventsByInscriptionId(data.InscriptionId)
		if err != nil {
			return err
		}

		for _, evt := range events {
			if evt.EventType == "transfer-transfer" {
				if evt.SourceWallet != data.FromWalletAddr {
					continue
				}
				if evt.SpentWallet != data.ToWalletAddr {
					continue
				}
				return nil
			}
		}
	}
	return ErrInvalidBrc20Transfer
}
