package enums

import (
	"fmt"
	"strconv"
)

type IntEnum interface {
	// ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
	~int
}

// EnumMap maps an integer enum value to a human-readable string
type IntEnumMap[T IntEnum] map[T]string

// Values retrieves a slice of all human-readable strings in the map
func (m IntEnumMap[T]) Values() []string {
	values := make([]string, len(m))
	i := 0
	for _, v := range m {
		values[i] = v
		i++
	}
	return values
}

// Keys retrieves a slice of all enum values in the map
func (m IntEnumMap[T]) Keys() []T {
	keys := make([]T, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

// GetText retrieves the human-readable string for an enum value
func (m IntEnumMap[T]) GetText(e T) (string, error) {
	text, ok := m[e]
	if !ok {
		return "", fmt.Errorf("invalid enum value %v", e)
	}
	return text, nil
}

// GetByText retrieves the enum value for a given human-readable string
func (m IntEnumMap[T]) GetByText(text string) (T, error) {
	for k, v := range m {
		if v == text {
			return k, nil
		}
	}

	var t T
	return t, fmt.Errorf("text %q does not map to any enum value", text)
}

func (m IntEnumMap[T]) ByIndex(keyValue any) (T, error) {
	var targetIndex int

	switch v := keyValue.(type) {
	case string:
		val, err := strconv.Atoi(v)
		if err != nil {
			return T(0), fmt.Errorf("int value %s does not map to any enum value", keyValue)
		}
		targetIndex = val
	case int:
		targetIndex = v
	default:
		return T(0), fmt.Errorf("int value %s does not map to any enum value", keyValue)
	}

	t := T(targetIndex)
	_, ok := m[t]
	if !ok {
		return t, fmt.Errorf("int value %d does not map to any enum value", targetIndex)
	}
	return t, nil
}
