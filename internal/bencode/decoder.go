package bencode

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

func Decode(bencode []byte) (BDict, error) {
	result := make(BDict)
	isKey := true
	var key string

	for i := 0; i < len(bencode); i++ {
		token := bencode[i]

		if equal(token, 'd') {
			continue
		} else if equal(token, 'i') { // only string can be a key, so always set value
			value, idx, err := findInt(bencode[i+1:])

			if err != nil {
				return BDict{}, fmt.Errorf("Error while decoding, index = %d, error = %s", i, err.Error())
			}

			i += idx

			result[key] = value
			isKey = true
		} else if unicode.IsDigit(rune(token)) {
			len, idx, err := findStringLength(bencode[i:])
			if err != nil {
				fmt.Printf("%s\n", result.GetString())
				return BDict{}, fmt.Errorf("Error while decoding, index = %d, error = %s", i, err.Error())
			}

			i += idx
			value := string(bencode[i : i+len])

			if isKey {
				key = value
				result[key] = struct{}{}
				isKey = false
			} else {
				result[key] = value
			}

			continue
		} else if equal(token, 'L') {
			// isKey = true
		}
	}

	return result, nil
}

func findStringLength(s []byte) (int, int, error) {
	len := 0

	for idx := range s {
		if s[idx] == ':' {
			res, err := strconv.Atoi(string(s[:len]))
			return res, len + 1, err
		}

		len++
	}

	return -1, -1, errors.New("Could not compute string length")
}

func findInt(s []byte) (BInt, int, error) {
	result := make([]byte, 0, 10)

	for i := range s {
		if s[i] == 'e' {
			res, err := strconv.Atoi(string(result))
			return BInt(res), i + 1, err
		}

		result = append(result, s[i])
	}

	return BInt(-1), -1, errors.New("Could not find int")
}

func equal(token byte, symbol rune) bool {
	return rune(token) == rune(symbol)
}
