package buffer

import (
	"bytes"
	"errors"
)

func BytesBuffer2Bytes(dst []byte, src *bytes.Buffer) error {
	buf := src.Bytes()
	if len(dst) < len(buf) {
		return errors.New("destination is too small")
	}
	for i := 0; i < len(buf); i++ {
		dst[i] = buf[i]
	}
	return nil
}
