package runestone

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/ordinox/btc-service/client"
	"github.com/ordinox/btc-service/config"
)

var (
	ErrNotEnoughConfirmations = errors.New("not enough confirmations")
	ErrNoRunestoneFound       = errors.New("no runestones found")
	ErrInvalidRunestone       = errors.New("invalid runestone")
	ErrParsingPkScript        = errors.New("error parsing pkscript")
)

type RunesDepositRequest struct {
	TxId             string
	FromAddr, ToAddr string
	RuneId           RuneId
	Amount           *big.Int
}

func VerifyRunesDeposit(request RunesDepositRequest, cfg config.Config) (bool, error) {
	btcClient := client.NewBitcoinClient(cfg)
	hash, err := chainhash.NewHashFromStr(request.TxId)

	if err != nil {
		return false, err
	}

	tx, err := btcClient.GetRawTransaction(hash)
	if err != nil {
		return false, err
	}

	runestone := DecipherRunestone(tx.MsgTx()).Runestone
	if runestone == nil {
		return false, fmt.Errorf("runestone is nil, %w", ErrNoRunestoneFound)
	}

	txVerbose, err := btcClient.GetRawTransactionVerbose(hash)
	if err != nil {
		return false, err
	}

	if txVerbose.Confirmations < 1 {
		return false, fmt.Errorf("less than 1 confirmation, %w", ErrNotEnoughConfirmations)
	}
	msgTx := tx.MsgTx()

	for _, e := range runestone.Edicts {
		if len(msgTx.TxOut) < int(e.Output) {
			return false, fmt.Errorf("invalid vout length: %w", ErrInvalidRunestone)
		}
		if len(msgTx.TxIn) < 1 {
			return false, fmt.Errorf("invalid vin length: %w", ErrInvalidRunestone)
		}

		if !e.Id.Equals(request.RuneId) {
			continue
		}
		if e.Amount.Cmp(request.Amount) != 0 {
			continue
		}
		pkscript, err := txscript.ParsePkScript(msgTx.TxOut[e.Output].PkScript)
		if err != nil {
			continue
		}
		destinationAddr, err := pkscript.Address(cfg.BtcConfig.GetChainConfigParams())
		if err != nil {
			continue
		}
		if destinationAddr.EncodeAddress() != request.ToAddr {
			continue
		}
		prevOutpoint := msgTx.TxIn[0].PreviousOutPoint
		tx, err := btcClient.GetRawTransaction(&prevOutpoint.Hash)
		if err != nil {
			continue
		}
		senderPkScript, err := txscript.ParsePkScript(tx.MsgTx().TxOut[prevOutpoint.Index].PkScript)
		if err != nil {
			continue
		}
		senderAddr, err := senderPkScript.Address(cfg.BtcConfig.GetChainConfigParams())
		if err != nil {
			continue
		}
		if senderAddr.EncodeAddress() != request.FromAddr {
			continue
		}
		return true, nil
	}

	return false, nil
}
