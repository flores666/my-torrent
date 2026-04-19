package bencode

import (
	"errors"
	"maps"
	"strconv"
	"unicode"
)

func Decode(bencode []byte) ([]BValue, error) {
	result := make([]BValue, 0)

	for i := 0; i < len(bencode); i++ {
		found := getBValue(bencode, i)
		if found.Error != nil {
			return result, found.Error
		}

		i += found.BytesVisited
		result = append(result, found.Value)
	}

	return result, nil
}

func getBValue(bencode []byte, i int) FoundValue {
	ch := bencode[i]

	switch {
	case ch == 'e':
		return FoundValue{
			Value:        nil,
			BytesVisited: -1,
			Error:        errors.New("End of value"),
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
		Value:        nil,
		BytesVisited: -1,
		Error:        errors.New("Value does not found"),
	}
}

func appendDictionary(from, to *BDict) {
	maps.Copy(*to, *from)
}

func findDictionary(i int, s []byte) FoundValue {
	return FoundValue{}
}

func findList(from int, bencode []byte) FoundValue {
	result := make([]BValue, 0)
	visited := 1 // skip l
	bencode = bencode[from+visited : len(bencode)-1]

	for i := 0; i < len(bencode); i++ {
		found := getBValue(bencode, i)
		if found.Error != nil {
			return found
		}

		i += found.BytesVisited
		visited += found.BytesVisited
		result = append(result, found.Value)
	}

	return FoundValue{
		Value:        result,
		BytesVisited: visited + 1,
		Error:        nil,
	}
}

func findString(i int, s []byte) FoundValue {
	len, idx, err := findStringLength(s[i:])
	if err != nil {
		return FoundValue{
			Value:        nil,
			BytesVisited: idx,
			Error:        err,
		}
	}

	idx += 1
	return FoundValue{
		Value:        BString(string(s[i+idx : i+idx+len])),
		BytesVisited: len + idx - 1,
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
			res, err := strconv.Atoi(string(result))
			return FoundValue{
				Value:        BInt(res),
				BytesVisited: i,
				Error:        err,
			}
		} else if unicode.IsDigit(rune(ch)) {
			result = append(result, ch)
		}
	}

	return FoundValue{
		Value:        BInt(-1),
		BytesVisited: -1,
		Error:        errors.New("Could not find int"),
	}
}
