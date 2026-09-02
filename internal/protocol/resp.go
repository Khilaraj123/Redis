// RESP (REdis Serialization Protocol) layer
// RESP type definitions (SimpleString, Error, Bulk, Array, Integer)
package protocol

type Kind int

const (
	KindSimpleString Kind = iota
	KindError
	KindInt
	KindBulk
	KindArray
	KindNil
)

// Value is a single RESP value, a reply on the wire or part of one
type Value struct {
	Kind  Kind
	Str   string
	Int   int64
	Array []Value
}

func SimpleString(s string) Value {
	return Value{
		Kind: KindSimpleString,
		Str:  s,
	}
}

func NewError(s string) Value {
	return Value{
		Kind: KindError,
		Str:  s,
	}
}

func Integer(n int64) Value {
	return Value{
		Kind: KindInt,
		Int:  n,
	}
}

func Bulk(s string) Value {
	return Value{
		Kind: KindBulk,
		Str:  s,
	}
}
func Nil() Value {
	return Value{
		Kind: KindNil,
	}
}
func Array(items ...Value) Value {
	return Value{
		Kind:  KindArray,
		Array: items,
	}
}
