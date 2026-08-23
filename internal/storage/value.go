// Value wrapper (type tag + payload + TTL)
package storage

type DataType int

const (
	StringType DataType = iota
	HashType
	ListType
	SetType
	SortedSetType
)

type Value struct{
	Type DataType
	Data any
}