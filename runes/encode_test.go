package runes

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/multiformats/go-varint"
)

func TestEncoding(t *testing.T) {
	tests := []string{"1", "3", "11", "12", "15", "1337", "2000", "100000", "1000000000000000000000000"}

	t.Run("encodeToSlice", func(t *testing.T) {
		for i, c := range tests {
			num, _ := big.NewInt(0).SetString(c, 10)
			e1 := hex.EncodeToString(encodeToSlice(num))
			fmt.Println("encodeSlice", i, e1)
		}
		fmt.Println("-------")
	})

	t.Run("encodeToSlice2", func(t *testing.T) {
		for i, c := range tests {
			num, _ := big.NewInt(0).SetString(c, 10)
			e1 := hex.EncodeToString(encodeToSlice2(num))
			fmt.Println("encodeSlice2", i, e1)
		}
		fmt.Println("-------")
	})

	t.Run("uLEB", func(t *testing.T) {
		for i, c := range tests {
			num, _ := big.NewInt(0).SetString(c, 10)
			e1 := hex.EncodeToString(AppendUleb128(nil, num.Uint64()))
			fmt.Println("uleb128", i, e1)
		}
		fmt.Println("-------")
	})

	t.Run("varInt", func(t *testing.T) {
		for i, c := range tests {
			num, _ := big.NewInt(0).SetString(c, 10)
			e1 := hex.EncodeToString(varint.ToUvarint(num.Uint64()))
			fmt.Println("varInt", i, e1)
		}
		fmt.Println("-------")
	})

}
