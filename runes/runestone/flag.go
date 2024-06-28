package runestone

import "math/big"

type Flag struct {
	val uint64
}

func (f Flag) Mask() *big.Int {
	return new(big.Int).Lsh(big.NewInt(1), uint(f.val))
}

func (f Flag) Take(flags *big.Int) bool {
	mask := f.Mask()
	set := new(big.Int).And(flags, mask).Cmp(big.NewInt(0)) != 0
	flags.And(flags, new(big.Int).Not(mask))
	return set
}

func (f Flag) Set(flags *big.Int) {
	flags.Or(flags, f.Mask())
}
