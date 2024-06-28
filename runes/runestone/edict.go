package runestone

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/wire"
)

type Edict struct {
	Id     RuneId
	Amount *big.Int
	Output Uint32
}

var (
	EmptyEdict       = Edict{Amount: big.NewInt(0)}
	ErrInvalidOutput = errors.New("invalid output")
)

func (e Edict) IsEmpty() bool {
	return e.Id.IsEmpty() && e.Output == 0 && e.Amount.Cmp(big.NewInt(0)) == 0
}

// NewEdictFromIntegers creates an Edict from the given parameters
func NewEdictFromIntegers(tx *wire.MsgTx, id RuneId, amount *big.Int, output *big.Int) (Edict, error) {
	if output == nil {
		return EmptyEdict, fmt.Errorf("rekt")
	}
	if !output.IsUint64() {
		return EmptyEdict, fmt.Errorf("output value %v is out of range for uint64", output)
	}

	o := int(output.Uint64())
	if o >= len(tx.TxOut) { // Changed to >= to handle out-of-bounds correctly
		return EmptyEdict, fmt.Errorf("%w: output index %d out of range", ErrInvalidOutput, o)
	}

	return Edict{Id: id, Amount: amount, Output: Uint32(o)}, nil
}
