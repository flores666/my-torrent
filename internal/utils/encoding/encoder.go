package encoding

import (
	"bytes"
	"crypto/sha1"
	"fmt"
)

func Sha1(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	hash := h.Sum(nil)

	var buf bytes.Buffer
	for _, b := range hash {
		fmt.Fprintf(&buf, "%%%02X", b)
	}

	return buf.String()
}
