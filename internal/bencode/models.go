package bencode

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type BValue struct {
	Value    any
	Original []byte
}
type BInt int64
type BDict map[string]BValue
type BString []byte
type BList []BValue

type FoundValue struct {
	Value        BValue
	BytesVisited int
	Error        error
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
	if value.Value == nil {
		return "nil"
	}

	if dict, ok := value.Value.(BDict); ok {
		return dict.GetString()
	}

	if list, ok := value.Value.(BList); ok {
		return list.GetString()
	}

	if intVal, ok := value.Value.(BInt); ok {
		return fmt.Sprintf("%d", int64(intVal))
	}

	if strVal, ok := value.Value.(BString); ok {
		return string(strVal)
	}

	if arr, ok := value.Value.([]BValue); ok {
		sb := strings.Builder{}

		for _, v := range arr {
			fmt.Fprintf(&sb, "%s, ", SprintfBValue(v))
		}

		return sb.String()
	}

	return fmt.Sprintf("%s", value)
}

func GetTypedValue[T any](value BValue) (T, error) {
	val := GetUnderlyingTypedValue(value)

	if res, ok := val.(T); ok {
		return res, nil
	}

	var zeroValue T
	return zeroValue, fmt.Errorf("cannot convert %s to type %s", SprintfBValue(value), reflect.TypeFor[T]().String())
}

func GetUnderlyingTypedValue(value BValue) any {
	if value.Value == nil {
		return nil
	}

	if val, err := getSimpleTypedValue(value); err == nil {
		return val
	}

	if dict, ok := value.Value.(BDict); ok {
		return dict.GetUnderlyingTypedValue()
	}

	if list, ok := value.Value.(BList); ok {
		return list.GetUnderlyingTypedValue()
	}

	if intVal, ok := value.Value.(BInt); ok {
		return int64(intVal)
	}

	if strVal, ok := value.Value.(BString); ok {
		return []byte(strVal)
	}

	if arr, ok := value.Value.([]BValue); ok {
		res := make([]any, 0)

		for _, v := range arr {
			res = append(res, GetUnderlyingTypedValue(v))
		}

		return res
	}

	panic(fmt.Sprintf("Value = %s, Type = %s", SprintfBValue(value), reflect.TypeOf(value).String()))
}

func getSimpleTypedValue(value BValue) (any, error) {
	if val, ok := value.Value.(int64); ok {
		return val, nil
	}

	if val, ok := value.Value.(map[string]any); ok {
		return val, nil
	}

	if val, ok := value.Value.([]byte); ok {
		return val, nil
	}

	if val, ok := value.Value.(string); ok {
		return val, nil
	}

	if val, ok := value.Value.([]any); ok {
		return val, nil
	}

	return nil, errors.New("Not simple type")
}
