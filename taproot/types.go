package taproot

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

type P2TRMetadata struct {
	Address             btcutil.Address
	ControlBlockWitness []byte
	PkScript            []byte
	TapHash             chainhash.Hash
	LockScript          []byte
}

type InscriptionData struct {
	Data        string
	ContentType string
}

var ContentTypeText = "text/plain;charset=utf-8"

func NewInscriptionData(data, contentType string) InscriptionData {
	return InscriptionData{
		Data:        data,
		ContentType: contentType,
	}
}
