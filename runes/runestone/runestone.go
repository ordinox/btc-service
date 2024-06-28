package runestone

import (
	"errors"
	"math"
	"math/big"
	"sort"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ordinox/btc-service/btc"
)

// Error definitions
var (
	ErrOverlong     = errors.New("overlong")
	ErrOverflow     = errors.New("overflow")
	ErrUnterminated = errors.New("unterminated")
)

type Runestone struct {
	Edicts  []Edict
	Mint    *RuneId
	Pointer *Uint32
}

const MAGIC_NUMBER = txscript.OP_13

func DecipherRunestone(tx *wire.MsgTx) Artifact {
	payload := payload(tx)
	if payload == nil {
		return Artifact{
			Cenotaph: &Cenotaph{Flaw: InvalidScript},
		}
	}
	if payload.Invalid != None {
		return Artifact{
			Cenotaph: &Cenotaph{Flaw: payload.Invalid},
		}
	}
	// Assuming that this is correct
	ints, err := integers(payload.Valid)
	if err != nil {
		return Artifact{Cenotaph: &Cenotaph{Flaw: Varint}}
	}

	msg := NewMessageFromIntegers(tx, ints)

	fields := msg.Fields
	flags, err := TakeFromTag(TagFlags, fields, 1, func(i []*big.Int) (*big.Int, error) { return i[0], nil })
	if err != nil {
		flags = new(big.Int).SetUint64(0)
	}
	// flags := TakeFromTag(TagFlags, fields, 1, func(i []big.Int) big.Int { return *new(big.Int) })

	mint, _ := TakeFromTag(TagMint, fields, 2, func(i []*big.Int) (*RuneId, error) {
		runeId := RuneId{}
		id, _ := runeId.Next(i[0], i[1])
		return &id, nil
	})

	pointer, _ := TakeFromTag(TagPointer, fields, 1, func(i []*big.Int) (*Uint32, error) {
		if !i[0].IsUint64() {
			return nil, nil
		}
		uPointer := i[0].Uint64()
		if uPointer > math.MaxUint32 {
			return nil, nil
		}
		u32Pointer := Uint32(uPointer)
		if int(u32Pointer) < len(tx.TxOut) {
			return &u32Pointer, nil
		}
		return nil, nil
	})

	flaw := None
	if flags != nil && flags.Uint64() != 0 {
		flaw = UnrecognizedFlag
	}

	for k := range fields {
		val, _ := new(big.Int).SetString(k, 10)
		if val.Uint64()%2 == 0 {
			if flaw == None {
				flaw = UnrecognizedEvenTag
			}
		}
	}

	if flaw != None {
		return Artifact{
			Cenotaph: &Cenotaph{
				Flaw:    Flaw(flaw),
				Mint:    mint,
				Etching: nil,
			},
		}
	}

	return Artifact{
		Runestone: &Runestone{
			Edicts:  msg.Edicts,
			Mint:    mint,
			Pointer: pointer,
		},
	}
}

func EncipherRunestone(runeStone Runestone) *txscript.ScriptBuilder {
	payload := make([]byte, 0)
	if runeStone.Mint != nil {
		payload = TagMint.Encode([]*big.Int{runeStone.Mint.Block.To64(), runeStone.Mint.Tx.To64()}, payload)
	}
	if runeStone.Pointer != nil {
		payload = TagPointer.Encode([]*big.Int{runeStone.Pointer.To64()}, payload)
	}

	if len(runeStone.Edicts) > 0 {
		payload = encodeToVec(TagBody.val.To64(), payload)
		edicts := make([]Edict, len(runeStone.Edicts))
		copy(edicts, runeStone.Edicts)
		// Sorting the edicts by the Block field of RuneID
		sort.Slice(edicts, func(i, j int) bool {
			if edicts[i].Id.Block == edicts[j].Id.Block {
				return edicts[i].Id.Tx < edicts[j].Id.Tx
			}
			return edicts[i].Id.Block < edicts[j].Id.Block
		})
		for _, e := range edicts {
			payload = encodeToVec(e.Id.Block.To64(), payload)
			payload = encodeToVec(e.Id.Tx.To64(), payload)
			payload = encodeToVec(e.Amount, payload)
			payload = encodeToVec(e.Output.To64(), payload)
		}
	}
	scriptBuilder := btc.NewScriptBuilder().AddOp(txscript.OP_RETURN).AddOp(MAGIC_NUMBER)

	for i := 0; i < len(payload); i += txscript.MaxScriptElementSize {
		end := i + txscript.MaxScriptElementSize
		if end > len(payload) {
			end = len(payload)
		}
		chunk := payload[i:end]
		scriptBuilder = scriptBuilder.AddData(chunk)
	}

	return scriptBuilder
}

type Payload struct {
	Valid   []byte
	Invalid Flaw
}

// payload searches transaction outputs for the payload
func payload(msg *wire.MsgTx) *Payload {
	for _, o := range msg.TxOut {
		if len(o.PkScript) == 0 || o.PkScript[0] != txscript.OP_RETURN {
			continue
		}

		if len(o.PkScript) < 2 || o.PkScript[1] != MAGIC_NUMBER {
			continue
		}

		var payload []byte
		i := 2

		for i < len(o.PkScript) {
			op := o.PkScript[i]
			i++

			if op <= txscript.OP_PUSHDATA4 {
				var dataSize int
				if op < txscript.OP_PUSHDATA1 {
					dataSize = int(op)
				} else if op == txscript.OP_PUSHDATA1 {
					if i >= len(o.PkScript) {
						return &Payload{Invalid: InvalidScript}
					}
					dataSize = int(o.PkScript[i])
					i++
				} else if op == txscript.OP_PUSHDATA2 {
					if i+1 >= len(o.PkScript) {
						return &Payload{Invalid: InvalidScript}
					}
					dataSize = int(o.PkScript[i]) | int(o.PkScript[i+1])<<8
					i += 2
				} else if op == txscript.OP_PUSHDATA4 {
					if i+3 >= len(o.PkScript) {
						return &Payload{Invalid: InvalidScript}
					}
					dataSize = int(o.PkScript[i]) | int(o.PkScript[i+1])<<8 | int(o.PkScript[i+2])<<16 | int(o.PkScript[i+3])<<24
					i += 4
				}

				if i+dataSize > len(o.PkScript) {
					return &Payload{Invalid: InvalidScript}
				}

				payload = append(payload, o.PkScript[i:i+dataSize]...)
				i += dataSize
			} else {
				return &Payload{Invalid: Opcode}
			}
		}

		return &Payload{Valid: payload, Invalid: None}
	}
	return nil
}

// integers decodes a byte slice into a slice of big.Int (simulating u128) using varint encoding
// AI Generated
// integers decodes a byte slice into a slice of *big.Int using varint encoding
func integers(payload []byte) ([]*big.Int, error) {
	var integers []*big.Int
	i := 0
	for i < len(payload) {
		integer, length, err := decode(payload[i:])
		if err != nil {
			return nil, err
		}
		integers = append(integers, integer)
		i += length
	}
	return integers, nil
}

// decode decodes a varint-encoded byte slice into a *big.Int and returns the integer and the number of bytes read
func decode(buffer []byte) (*big.Int, int, error) {
	n := big.NewInt(0)
	for i, byteValue := range buffer {
		if i > 18 {
			return nil, 0, ErrOverlong
		}
		value := big.NewInt(int64(byteValue & 0b0111_1111))
		if i == 18 && (value.Int64()&0b0111_1100) != 0 {
			return nil, 0, ErrOverflow
		}
		n.Or(n, new(big.Int).Lsh(value, uint(7*i)))
		if byteValue&0b1000_0000 == 0 {
			return n, i + 1, nil
		}
	}
	return nil, 0, ErrUnterminated
}
