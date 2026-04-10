package bencode

import (
	"fmt"
	"strings"
)

type BValue interface{}
type BInt int
type BDict map[string]BValue
type BString []byte
type BList []BValue

func (bs *BString) ToString() string {
	return string(*bs)
}

func (d *BDict) GetString() string {
	sb := strings.Builder{}

	for k, v := range *d {
		fmt.Fprintf(&sb, "key = %s, value = %s\n", k, v)
	}

	return sb.String()
}
