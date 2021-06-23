package fsnapshot

import (
	"crypto/sha256"
	"io"
	"os"
)

func checksum(d []byte) []byte {
	h := sha256.New()
	return h.Sum(d)
}

func checksumfile(f *os.File) ([]byte, error) {
	h := sha256.New()
	d := make([]byte, 1024*1000*10) // Read in 10MB chunks
	for {
		b, err := f.Read(d)
		if err != nil && err != io.EOF {
			return nil, err
		}

		if b == 0 {
			break
		}
		h.Write(d[:b])
	}
	return h.Sum(nil), nil
}
