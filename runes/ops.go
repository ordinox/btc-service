package runes

import (
	"math/big"

	"github.com/btcsuite/btcd/txscript"
	"github.com/multiformats/go-varint"
)

const (
	OP_RETURN = txscript.OP_RETURN
	OP_MAGIC  = txscript.OP_13
)

// RUNESTONE DATA PUSH CODES
var (
	MINT   uint64 = 20
	BODY   uint64 = 0
	AMOUNT uint64 = 10
)

// encodeULEB128 encodes a big.Int into a slice of bytes using unsigned LEB128 encoding.
func encodeULEB128(value *big.Int) []byte {
	var result []byte
	if value.Sign() < 0 {
		panic("encodeULEB128 only supports non-negative integers")
	}

	zero := big.NewInt(0)
	base := big.NewInt(128)
	mod := new(big.Int)

	for value.Cmp(zero) > 0 {
		// Compute the current byte
		mod.Mod(value, base)
		currentByte := mod.Uint64()
		value.Div(value, base)

		// If there are more bytes to encode, set the continuation bit (bit 7).
		if value.Cmp(zero) != 0 {
			currentByte |= 0x80
		}

		result = append(result, byte(currentByte))
	}

	if len(result) == 0 {
		result = append(result, 0)
	}

	return result
}

// Varint encoder, ported from
// https://github.com/ordinals/ord/blob/1e6cb641faf3b1eb0aba501a7a2822d7a3dc8643/crates/ordinals/src/varint.rs#L3-L39
// Using big.Int since go doesn't support u128
func encodeToSlice(n *big.Int) []byte {
	var result []byte
	var oneTwentyEight = big.NewInt(128)

	for n.Cmp(oneTwentyEight) >= 0 {
		temp := new(big.Int).Mod(n, oneTwentyEight)
		tempByte := byte(temp.Uint64()) | 0x80
		result = append(result, tempByte)
		n.Div(n, oneTwentyEight)
	}
	result = append(result, byte(n.Uint64()))
	return result
}

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

func NewEdict(rune Rune, amount *big.Int, output uint64) []byte {
	data := make([]byte, 0)

	data = append(data, encodeToSlice(big.NewInt(int64(rune.BlockNumber)))...)
	data = append(data, encodeToSlice(big.NewInt(int64(rune.TxIndex)))...)
	data = append(data, encodeToSlice(amount)...)
	data = append(data, encodeToSlice(big.NewInt(int64(1)))...)

	// Keep this for future reference & testing
	// data = append(data, ToVarInt(rune.BlockNumber)...)
	// data = append(data, ToVarInt(uint64(rune.TxIndex))...)
	// data = append(data, AppendSleb128(nil, amount.Int64())...)
	// data = append(data, ToVarInt(output)...)
	return data
}
