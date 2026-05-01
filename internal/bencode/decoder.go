package bencode

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"unicode"
)

func Decode(bencode []byte) ([]BValue, error) {
	result := make([]BValue, 0)

	for i := 0; i < len(bencode); {
		found := getBValue(bencode, i)
		if found.Error != nil {
			return result, found.Error
		}

		if found.Value.Value == nil {
			return result, nil
		}

		i += found.BytesVisited
		result = append(result, found.Value)
	}

	return result, nil
}

func getBValue(bencode []byte, i int) (res FoundValue) {
	ch := bencode[i]
	switch {
	case ch == 'e':
		return FoundValue{
			Value:        BValue{},
			BytesVisited: -1,
			Error:        nil,
		}
	case ch == 'd':
		return findDictionary(i, bencode)
	case ch == 'l':
		return findList(i, bencode)
	case ch == 'i':
		return findInt(i, bencode)
	case unicode.IsDigit(rune(ch)):
		return findString(i, bencode)
	}

	return FoundValue{
		Value:        BValue{},
		BytesVisited: -1,
		Error:        fmt.Errorf("Value does not found for char = %c, index - %d", ch, i),
	}
}

func findDictionary(from int, bencode []byte) FoundValue {
	result := make(BDict)
	isKey := true
	key := ""
	visited := 1
	bencode = bencode[from:]

	for i := 1; i < len(bencode); {
		if bencode[i] == 'e' {
			return FoundValue{
				Value: BValue{
					Value:    result,
					Original: bencode[:visited+1],
				},
				BytesVisited: i + 1,
				Error:        nil,
			}
		}

		found := getBValue(bencode, i)
		if found.Error != nil {
			return found
		}

		i += found.BytesVisited
		visited += found.BytesVisited

		if isKey {
			bStringKey, ok := found.Value.Value.(BString)
			if ok {
				key = string(bStringKey)
				result[key] = BValue{}
				isKey = false
			} else {
				return FoundValue{
					Value: BValue{
						Value:    result,
						Original: bencode[:visited+1],
					},
					BytesVisited: -1,
					Error:        fmt.Errorf("Invalid key type - %s, must be string, Value - %s", reflect.TypeFor[BValue](), SprintfBValue(found.Value)),
				}
			}
		} else {
			result[key] = found.Value
			isKey = true
		}
	}

	return FoundValue{
		Value: BValue{
			Value:    result,
			Original: bencode[:visited+1],
		},
		BytesVisited: visited + 1,
		Error:        nil,
	}
}

func findList(from int, bencode []byte) FoundValue {
	result := make([]any, 0)
	visited := 1
	bencode = bencode[from:]

	for i := 1; i < len(bencode); {
		if bencode[i] == 'e' {
			return FoundValue{
				Value: BValue{
					Value:    result,
					Original: bencode[:i+1],
				},
				BytesVisited: i + 1,
				Error:        nil,
			}
		}

		found := getBValue(bencode, i)

		if found.Error != nil {
			return found
		}

		i += found.BytesVisited
		visited += found.BytesVisited
		result = append(result, GetUnderlyingTypedValue(found.Value))
	}

	return FoundValue{
		Value: BValue{
			Value:    result,
			Original: bencode[:visited+1],
		},
		BytesVisited: visited + 1,
		Error:        nil,
	}
}

func findString(i int, s []byte) FoundValue {
	len, idx, err := findStringLength(s[i:])
	if err != nil {
		return FoundValue{
			Value:        BValue{},
			BytesVisited: idx,
			Error:        err,
		}
	}

	idx += 1
	return FoundValue{
		Value: BValue{
			Value:    BString(s[i+idx : i+idx+len]),
			Original: s[i+idx : i+idx+len],
		},
		BytesVisited: len + idx,
		Error:        nil,
	}
}

func findStringLength(s []byte) (int, int, error) {
	len := 0

	for _, ch := range s {
		if ch == ':' {
			res, err := strconv.Atoi(string(s[:len]))
			return res, len, err
		}

		len++
	}

	return -1, -1, errors.New("Could not compute string length")
}

func findInt(from int, s []byte) FoundValue {
	result := make([]byte, 0, 10)

	for i, ch := range s[from:] {
		if ch == 'e' {
			res, err := strconv.ParseInt(string(result), 10, 64)
			return FoundValue{
				Value: BValue{
					Value:    BInt(res),
					Original: result,
				},
				BytesVisited: i + 1,
				Error:        err,
			}
		} else if unicode.IsDigit(rune(ch)) {
			result = append(result, ch)
		}
	}

	return FoundValue{
		Value: BValue{
			Value: BInt(-1),
		},
		BytesVisited: -1,
		Error:        errors.New("Could not find int"),
	}
}
