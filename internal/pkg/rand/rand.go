package rand

import (
	rd "crypto/rand"
	"encoding/hex"
	"io"
)

func RandomBytes(n int) []byte {
	b := make([]byte, n)
	bs, err := io.ReadFull(rd.Reader, b)
	if err != nil {
		return nil
	}
	if bs < n {
		return nil
	}

	return b
}

func RandomString(n int) string {
	b := RandomBytes(n)
	return hex.EncodeToString(b)
}
