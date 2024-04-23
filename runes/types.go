package runes

import (
	"fmt"
	"strconv"
	"strings"
)

type Rune struct {
	BlockNumber uint64
	TxIndex     uint32
}

func (r Rune) String() string {
	return fmt.Sprintf("%d:%d", r.BlockNumber, r.TxIndex)
}

func ParseRune(runeStr string) (*Rune, error) {
	split := strings.Split(runeStr, ":")
	blockNumber, err := strconv.Atoi(split[0])
	if err != nil {
		return nil, fmt.Errorf("error: invalid Rune string %s", runeStr)
	}

	txIdx, err := strconv.Atoi(split[1])
	if err != nil {
		return nil, fmt.Errorf("Error: invalid Rune ID %s", runeStr)
	}
	return &Rune{
		TxIndex:     uint32(txIdx),
		BlockNumber: uint64(blockNumber),
	}, nil
}
