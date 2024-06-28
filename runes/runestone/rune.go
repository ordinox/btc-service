package runestone

import (
	"errors"
	"fmt"
	"math/big"
	"math/bits"
	"strings"
	"unicode"
)

type Rune struct {
	val big.Int
}

type SpacedRune struct {
	Rune    Rune
	Spacers uint32
}

var (
	ErrLeadingSpacer  = errors.New("leading spacer")
	ErrDoubleSpacer   = errors.New("double spacer")
	ErrCharacter      = errors.New("invalid character")
	ErrTrailingSpacer = errors.New("trailing spacer")
	ErrRuneParse      = errors.New("failed to parse rune")
)

func fromStr(s string) (SpacedRune, error) {
	var runeStr strings.Builder
	var spacers uint32

	for _, c := range s {
		if unicode.IsUpper(c) {
			runeStr.WriteRune(c)
		} else if c == '.' || c == '•' {
			flag := uint32(1) << (runeStr.Len() - 1)
			if runeStr.Len() == 0 {
				return SpacedRune{}, ErrLeadingSpacer
			}
			if spacers&flag != 0 {
				return SpacedRune{}, ErrDoubleSpacer
			}
			spacers |= flag
		} else {
			return SpacedRune{}, fmt.Errorf("%w: %c", ErrCharacter, c)
		}
	}

	if 32-bits.LeadingZeros32(spacers) >= runeStr.Len() {
		return SpacedRune{}, ErrTrailingSpacer
	}

	var runeVal big.Int
	if _, ok := runeVal.SetString(runeStr.String(), 10); !ok {
		return SpacedRune{}, ErrRuneParse
	}

	return SpacedRune{
		Rune:    Rune{val: runeVal},
		Spacers: spacers,
	}, nil
}
