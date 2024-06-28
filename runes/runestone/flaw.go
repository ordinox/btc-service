package runestone

type Flaw int

const (
	EdictOutput = iota
	EdictRuneId
	InvalidScript
	Opcode
	SupplyOverflow
	TrailingIntegers
	TruncatedField
	UnrecognizedEvenTag
	UnrecognizedFlag
	Varint
	None
)
