package compress

import (
	"bytes"
	"io"

	"github.com/andybalholm/brotli"
)

func WriteBrotli(data []byte, level int, buf io.Writer) error {
	brt := brotli.NewWriterLevel(buf, level)
	if _, err := brt.Write(data); err != nil {
		return err
	}
	if err := brt.Flush(); err != nil {
		return err
	}
	if err := brt.Close(); err != nil {
		return err
	}
	return nil
}

func ReadBrotli(reader io.Reader, writer io.Writer) error {
	brt := brotli.NewReader(reader)
	buf := bytes.NewBuffer(nil)
	temp := make([]byte, 4096)
	for {
		n, err := brt.Read(temp)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}
		if n == 0 {
			break
		}
		buf.Write(temp[:n])
	}
	return nil
}
