package bencode

import (
	"fmt"
	"strings"
)

type BValue any
type BInt int
type BDict map[string]BValue
type BString []byte
type BList []BValue

type FoundValue struct {
	Value        BValue
	BytesVisited int
	Error        error
}

func (bs *BString) ToString() string {
	return string(*bs)
}

func (d *BDict) GetString() string {
	sb := strings.Builder{}

	for k, v := range *d {
		fmt.Fprintf(&sb, "key = %s, value = %s\n", k, SprintfBValue(v))
	}

	return sb.String()
}

func (d *BList) GetString() string {
	sb := strings.Builder{}

	for _, v := range *d {
		fmt.Fprintf(&sb, "%s", SprintfBValue(v))
	}

	return sb.String()
}

func SprintfBValue(value BValue) string {
	if dict, ok := value.(BDict); ok {
		return dict.GetString()
	}

	if list, ok := value.(BList); ok {
		return list.GetString()
	}

	if intVal, ok := value.(BInt); ok {
		return fmt.Sprintf("%d", intVal)
	}

	if strVal, ok := value.(BString); ok {
		return fmt.Sprintf("%s", strVal)
	}

	if arr, ok := value.([]BValue); ok {
		sb := strings.Builder{}

		for _, v := range arr {
			fmt.Fprintf(&sb, "%s, ", SprintfBValue(v))
		}

		return sb.String()
	}

	return "unknown"
}
