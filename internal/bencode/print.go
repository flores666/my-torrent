package bencode

import "fmt"

func PrintParsed(values []BValue) {
	for _, val := range values {
		fmt.Println(SprintfBValue(val))
	}
}
