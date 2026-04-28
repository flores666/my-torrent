package bencode

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type BValue any
type BInt int
type BDict map[string]BValue
type BString []byte
type BList []BValue

type FoundValue struct {
	Value        any
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

func (d *BDict) GetUnderlyingTypedValue() map[string]any {
	res := make(map[string]any)

	for k, v := range *d {
		res[k] = GetUnderlyingTypedValue(v)
	}

	return res
}

func (d *BList) GetString() string {
	sb := strings.Builder{}

	for _, v := range *d {
		fmt.Fprintf(&sb, "%s", SprintfBValue(v))
	}

	return sb.String()
}

func (d *BList) GetUnderlyingTypedValue() []any {
	res := make([]any, 0)

	for _, v := range *d {
		res = append(res, GetUnderlyingTypedValue(v))
	}

	return res
}

func SprintfBValue(value BValue) string {
	if value == nil {
		return "nil"
	}

	if dict, ok := value.(BDict); ok {
		return dict.GetString()
	}

	if list, ok := value.(BList); ok {
		return list.GetString()
	}

	if intVal, ok := value.(BInt); ok {
		return fmt.Sprintf("%d", int(intVal))
	}

	if strVal, ok := value.(BString); ok {
		return strVal.ToString()
	}

	if arr, ok := value.([]BValue); ok {
		sb := strings.Builder{}

		for _, v := range arr {
			fmt.Fprintf(&sb, "%s, ", SprintfBValue(v))
		}

		return sb.String()
	}

	return fmt.Sprintf("%s", value)
}

func GetUnderlyingTypedValue(value BValue) any {
	if value == nil {
		return nil
	}

	if val, err := getSimpleTypedValue(value); err == nil {
		return val
	}

	if dict, ok := value.(BDict); ok {
		return dict.GetUnderlyingTypedValue()
	}

	if list, ok := value.(BList); ok {
		return list.GetUnderlyingTypedValue()
	}

	if intVal, ok := value.(BInt); ok {
		return int(intVal)
	}

	if strVal, ok := value.(BString); ok {
		return strVal.ToString()
	}

	if arr, ok := value.([]BValue); ok {
		res := make([]any, 0)

		for _, v := range arr {
			res = append(res, GetUnderlyingTypedValue(v))
		}

		return res
	}

	panic(fmt.Sprintf("Value = %s, Type = %s", SprintfBValue(value), reflect.TypeOf(value).String()))
}

func getSimpleTypedValue(value BValue) (any, error) {
	if val, ok := value.(int); ok {
		return val, nil
	}

	if val, ok := value.(string); ok {
		return val, nil
	}

	if val, ok := value.([]any); ok {
		return val, nil
	}

	if val, ok := value.(map[string]any); ok {
		return val, nil
	}

	return nil, errors.New("Not simple type")
}
