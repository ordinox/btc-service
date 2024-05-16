package brc20

import (
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ordinox/btc-service/common"
	"github.com/ordinox/btc-service/config"
)

/*
Transfer
{
  "commit": "a6105c9776fa8d5fd2bc1bea12a4a8c982d9a1f2d33b6e34e4a1120891e3e884",
  "inscriptions": [
    {
      "id": "2402cf6e1bcba65257bc4eccc6035cb306dbc79f7eb3843f52144ed1f554e66di0",
      "location": "2402cf6e1bcba65257bc4eccc6035cb306dbc79f7eb3843f52144ed1f554e66d:0:0"
    }
  ],
  "parent": null,
  "reveal": "2402cf6e1bcba65257bc4eccc6035cb306dbc79f7eb3843f52144ed1f554e66d",
  "total_fees": 59400
}
*/

func xTestTransf(t *testing.T) {
	senderPk := common.LoadPrivateKey(common.TEST_PRIV_KEY_HEX_1)
	destinationPk := common.LoadPrivateKey(common.TEST_PRIV_KEY_HEX_2)

	senderAddrPubkey, _ := btcutil.NewAddressPubKey(senderPk.PubKey().SerializeUncompressed(), &chaincfg.RegressionNetParams)
	destinationAddrPubkey, _ := btcutil.NewAddressPubKey(destinationPk.PubKey().SerializeUncompressed(), &chaincfg.RegressionNetParams)

	senderAddr := senderAddrPubkey.AddressPubKeyHash()
	destinationAddr := destinationAddrPubkey.AddressPubKeyHash()

	fmt.Println("Sender", senderAddr.EncodeAddress())
	fmt.Println("Destination", destinationAddr.EncodeAddress())

	utxos, err := common.GetUtxos(senderAddr.EncodeAddress(), config.GetDefaultConfig().BtcConfig)
	cUtxos := utxos.Result.ToUtxo()
	if err != nil {
		panic(err)
	}
	_, err = Transfer(cUtxos[0], cUtxos[2], senderAddr, destinationAddr, senderPk, senderPk.PubKey(), 25, config.GetDefaultConfig())
	if err != nil {
		panic(err)
	}
}
