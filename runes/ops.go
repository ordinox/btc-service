package runes

import (
	"github.com/btcsuite/btcd/txscript"
	"github.com/multiformats/go-varint"
)

const (
	OP_RETURN = txscript.OP_RETURN
	OP_MAGIC  = txscript.OP_13
)

var (
	MINT   uint64 = 20
	BODY   uint64 = 0
	AMOUNT uint64 = 10
)

func TagToVarInt(tag uint64, values ...uint64) []byte {
	data := make([]byte, 0)
	for _, v := range values {
		data = append(data, ToVarInt(tag)...)
		data = append(data, ToVarInt(v)...)
	}
	return data
}

func ToVarInt(i uint64) []byte {
	return varint.ToUvarint(i)
}

func NewEdict(rune Rune, amount, output uint64) []byte {
	data := make([]byte, 0)
	data = append(data, ToVarInt(rune.BlockNumber)...)
	data = append(data, ToVarInt(uint64(rune.TxIndex))...)
	data = append(data, ToVarInt(amount)...)
	data = append(data, ToVarInt(output)...)
	return data
}
