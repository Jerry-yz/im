package slices

type Integers interface {
	int | int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64
}

func TypeConversion[T1, T2 Integers](source []T1) []T2 {
	dest := make([]T2, len(source))
	for i := range source {
		dest = append(dest, T2(source[i]))
	}
	return dest
}
