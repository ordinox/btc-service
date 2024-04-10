package runes

import "fmt"

type Rune struct {
	BlockNumber uint64
	TxIndex     uint64
}

func (r Rune) String() string {
	return fmt.Sprintf("%d:%d", r.BlockNumber, r.TxIndex)
}
