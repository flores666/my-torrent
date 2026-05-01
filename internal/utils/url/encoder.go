package url

import (
	"bytes"
	"crypto/sha1"
	"fmt"
)

func Encode(s []byte) string {
	h := sha1.New()
	h.Write(s)
	hash := h.Sum(nil)

	var buf bytes.Buffer
	for _, b := range hash {
		fmt.Fprintf(&buf, "%%%02X", b)
	}

	return buf.String()
}
