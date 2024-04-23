package taproot

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/ordinox/btc-service/config"
)

// Given a private key, derive all the obvious parameters for a P2TR transaction
func CreateP2TRMetaData(privateKey *btcec.PrivateKey, config config.Config) (*P2TRMetadata, error) {
	scriptBuilder := txscript.NewScriptBuilder()
	scriptBuilder.
		AddData(schnorr.SerializePubKey(privateKey.PubKey())).
		AddOp(txscript.OP_CHECKSIG)
	script, err := scriptBuilder.Script()
	if err != nil {
		return nil, err
	}
	leafNode := txscript.NewBaseTapLeaf(script)
	proof := txscript.TapscriptProof{
		TapLeaf:  leafNode,
		RootNode: leafNode,
	}
	controlBlock := proof.ToControlBlock(privateKey.PubKey())
	controlBlockWitness, err := controlBlock.ToBytes()
	if err != nil {
		return nil, err
	}

	tapHash := proof.RootNode.TapHash()
	address, err := btcutil.NewAddressTaproot(
		schnorr.SerializePubKey(
			txscript.ComputeTaprootOutputKey(
				privateKey.PubKey(),
				tapHash[:],
			),
		),
		config.BtcConfig.GetChainConfigParams(),
	)
	if err != nil {
		return nil, err
	}
	return &P2TRMetadata{
		Address:             address,
		ControlBlockWitness: controlBlockWitness,
		PkScript:            script,
		TapHash:             tapHash,
	}, nil
}
