package runestone

import (
	"fmt"
	"math/big"

	"github.com/multiformats/go-varint"
)

type Tag struct {
	val Uint64
}

var (
	TagBody        = Tag{0}
	TagFlags       = Tag{2}
	TagRune        = Tag{4}
	TagPremine     = Tag{6}
	TagCap         = Tag{8}
	TagAmount      = Tag{10}
	TagHeightStart = Tag{12}
	TagHeightEnd   = Tag{14}
	TagOffsetStart = Tag{16}
	TagOffsetEnd   = Tag{18}
	TagMint        = Tag{20}
	TagPointer     = Tag{22}
	TagCenotaph    = Tag{126}

	TagDivisibility = Tag{1}
	TagSpacers      = Tag{3}
	TagSymbol       = Tag{5}
	TagNop          = Tag{127}
)

/**
  let field = fields.get_mut(&self.into())?;

  let mut values: [u128; N] = [0; N];

  for (i, v) in values.iter_mut().enumerate() {
    *v = *field.get(i)?;
  }

  let value = with(values)?;

  field.drain(0..N);

  if field.is_empty() {
    fields.remove(&self.into()).unwrap();
  }

  Some(value)
*/

func TakeFromTag[T any](t Tag, fields map[string][]*big.Int, N int, with func([]*big.Int) (T, error)) (T, error) {
	// Convert the tag value to a string key
	val := new(big.Int).SetUint64(uint64(t.val)).String()
	field, ok := fields[val]
	if !ok || len(field) < N {
		return *new(T), fmt.Errorf("rekt: field not found or not enough elements")
	}

	// Prepare the values slice with the first N elements
	values := make([]*big.Int, N)
	for i := range values {
		values[i] = field[i]
	}

	// Call the 'with' function with the first 'N' elements
	res, err := with(values)
	if err != nil {
		return *new(T), err
	}

	// Remove the first 'N' elements from 'field'
	fields[val] = field[N:]

	// If the field is empty after draining, remove it from the map
	if len(fields[val]) == 0 {
		delete(fields, val)
	}

	return res, nil
}
func (t Tag) Encode(values []*big.Int, payload []byte) (res []byte) {
	for _, value := range values {
		payload = encodeToVec(new(big.Int).SetUint64(uint64(t.val)), payload)
		payload = encodeToVec(value, payload)
	}
	return payload
}

// encodeToVec encodes a big.Int (simulating u128) into a byte slice using a similar varint encoding method
func encodeToVec(n *big.Int, v []byte) []byte {
	return append(v, varint.ToUvarint(n.Uint64())...)
}
