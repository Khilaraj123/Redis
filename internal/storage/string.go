// String data type ops
package storage

func (s *Store) SetString(key, value string) {
	s.Set(key, Value{
		Type: StringType,
		Data: value,
	})
}

func (s *Store) Set(key string, param any) {
	panic("unimplemented")
}
