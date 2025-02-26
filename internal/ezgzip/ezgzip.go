package ezgzip

import (
	"bytes"
	"compress/gzip"
	"io"
)

func MustDecompressBytes(b []byte) []byte {
	r, err := DecompressBytes(b)
	if err != nil {
		panic(err)
	}
	return r
}

func DecompressBytes(b []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}
