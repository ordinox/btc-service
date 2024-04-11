package taproot

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/ordinox/btc-service/config"
)

func CreateP2TRInscriptionMetaData(inscription InscriptionData, privateKey *btcec.PrivateKey, config config.Config) (*P2TRMetadata, error) {
	scriptBuilder := txscript.NewScriptBuilder()
	scriptBuilder.
		AddData(schnorr.SerializePubKey(privateKey.PubKey())).
		AddOp(txscript.OP_CHECKSIG).
		AddOp(txscript.OP_FALSE).
		AddOp(txscript.OP_IF).
		AddData([]byte("ord")).
		AddOp(txscript.OP_1).
		AddData([]byte(inscription.ContentType)).
		AddOp(txscript.OP_0).
		AddData([]byte(inscription.Data)).
		AddOp(txscript.OP_ENDIF)

	script, err := scriptBuilder.Script()
	if err != nil {
		return nil, err
	}
	// script = append(script, txscript.OP_ENDIF)
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
	pkScript, err := txscript.PayToAddrScript(address)
	if err != nil {
		return nil, err
	}
	return &P2TRMetadata{
		Address:             address,
		ControlBlockWitness: controlBlockWitness,
		PkScript:            pkScript,
		TapHash:             tapHash,
		LockScript:          script,
	}, nil
}
