package encode

import (
	"bytes"
	"fmt"
)

func Percent(b []byte) string {
	var buf bytes.Buffer
	for _, c := range b {
		fmt.Fprintf(&buf, "%%%02X", c)
	}
	return buf.String()
}
