package runestone

import (
	"math/big"

	"github.com/btcsuite/btcd/wire"
)

type Message struct {
	Flaw   Flaw
	Edicts []Edict
	Fields map[string][]*big.Int
}

type BigInts []*big.Int

func (b BigInts) String() (s string) {
	for _, i := range b {
		s += i.String() + " "
	}
	return
}

func NewMessageFromIntegers(tx *wire.MsgTx, payload BigInts) Message {
	edicts := make([]Edict, 0)
	fields := make(map[string][]*big.Int)
	flaw := None

	for i := 0; i < len(payload); i += 2 {
		tag := payload[i].Uint64()

		if tag == uint64(TagBody.val) {
			id := RuneId{}
			for j := i + 1; j < len(payload); j += 4 {
				if len(payload[j:]) < 4 {
					if flaw == None {
						flaw = TrailingIntegers
						break
					}
				}
				chunk := payload[j : j+4]
				newId, err := id.Next(chunk[0], chunk[1])
				if err != nil {
					if flaw == None {
						flaw = EdictRuneId
						break
					}
				}
				newEdict, err := NewEdictFromIntegers(tx, newId, chunk[2], chunk[3])
				if err != nil {
					if flaw == None {
						flaw = EdictOutput
						break
					}
				}
				id = newId
				edicts = append(edicts, newEdict)
			}
			break
		}

		if i+1 >= len(payload) {
			if flaw == None {
				flaw = TruncatedField
				break
			}
		}
		value := payload[i+1]
		idx := new(big.Int).SetUint64((tag)).String()
		v, exists := fields[idx]
		if !exists {
			v = make([]*big.Int, 0)
		}
		v = append(v, value)
		fields[idx] = v
	}
	return Message{
		Edicts: edicts,
		Flaw:   Flaw(flaw),
		Fields: fields,
	}
}
