package inscriptions

import (
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/ordinox/btc-service/common"
	"github.com/ordinox/btc-service/config"
	"github.com/ordinox/btc-service/taproot"
)

const PK_HEX = "f80c5f4802ac331207fe47f5a36cb4a0c17a0dbd4fe57f4d1243d56e4000e79a"

const DEPLOY_TXT = "{\"p\":\"brc-20\",\"op\":\"deploy\",\"tick\":\"ORDX\",\"max\":\"10000000\",\"lim\":\"10000000\"}"
const MINT_TXT = "{\"p\":\"brc-20\",\"op\":\"mint\",\"tick\":\"ORDX\",\"amt\":\"10\"}"

var _ = fmt.Println

func TestInscriptionCommit(t *testing.T) {
	config.Init()
	config := config.GetDefaultConfig()
	pk := common.LoadPrivateKey(PK_HEX)
	pk2, _ := btcec.NewPrivateKey()
	addr2, _ := btcutil.NewAddressTaproot(schnorr.SerializePubKey(pk2.PubKey()), config.BtcConfig.GetChainConfigParams())
	_, err := InscribeNative(
		addr2,
		pk,
		taproot.NewInscriptionData(MINT_TXT, taproot.ContentTypeText),
		20,
		config,
	)
	if err != nil {
		panic(err)
	}
}
