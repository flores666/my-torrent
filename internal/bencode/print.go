package bencode

import (
	"encoding/json"
	"fmt"
)

func PrintParsed(values []BValue) {
	bytes, _ := json.Marshal(values)
	fmt.Println(string(bytes))
}
