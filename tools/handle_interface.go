package tools

func Num64(valueType any) int64 {
	switch value := valueType.(type) {
	case int:
		return int64(value)
	case int8:
		return int64(value)
	case int16:
		return int64(value)
	case int32:
		return int64(value)
	case int64:
		return int64(value)
	case uint:
		return int64(uint64(value))
	case uintptr:
		return int64(uint64(value))
	case uint8:
		return int64(uint64(value))
	case uint16:
		return int64(uint64(value))
	case uint32:
		return int64(uint64(value))
	case uint64:
		return int64(uint64(value))
	}
	return 0
}
