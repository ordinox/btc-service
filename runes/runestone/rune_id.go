package runestone

import (
	"fmt"
	"math/big"
)

type (
	Uint64 uint64
	Uint32 uint32
)

type RuneId struct {
	Block Uint64
	Tx    Uint32
}

func (r RuneId) IsEmpty() bool {
	return r.Block == 0 && r.Tx == 0
}

func (u Uint32) To64() *big.Int {
	return new(big.Int).SetUint64(uint64(u))
}

func (u Uint64) To64() *big.Int {
	return new(big.Int).SetUint64(uint64(u))
}

var (
	EmptyRuneId = RuneId{Block: 0, Tx: 0}
)

func NewRuneId(block uint64, tx uint32) RuneId {
	if block == 0 && tx > 0 {
		return EmptyRuneId
	}
	return RuneId{Block: Uint64(block), Tx: Uint32(tx)}
}

func (r RuneId) Delta(next RuneId) (*big.Int, *big.Int) {
	return new(big.Int), new(big.Int)
}

func (r RuneId) Next(block *big.Int, tx *big.Int) (RuneId, error) {
	if !block.IsUint64() {
		return RuneId{}, fmt.Errorf("block value %v is out of range for uint64", block)
	}
	if !tx.IsUint64() {
		return RuneId{}, fmt.Errorf("tx value %v is out of range for uint64", tx)
	}

	blockUint64 := block.Uint64()
	txUint64 := tx.Uint64()

	if blockUint64 == 0 {
		txUint64 += uint64(r.Tx)
	}

	return NewRuneId(
		uint64(r.Block)+blockUint64,
		uint32(txUint64),
	), nil
}

func (r RuneId) String() string {
	return fmt.Sprintf("%d:%d", r.Block, r.Tx)
}
