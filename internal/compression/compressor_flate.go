package compression

import (
	"bytes"
	"io"

	errs "github.com/guionardo/gs-dev/internal/errors"
	"github.com/klauspost/compress/flate"
)

type FlateCompressor struct {
}

func NewFlateCompressor() *FlateCompressor {
	return &FlateCompressor{}
}

func (c *FlateCompressor) Compress(data []byte) (compressed []byte, err error) {
	in := bytes.NewReader(data)
	out := bytes.NewBuffer(make([]byte, 0, len(data)))

	enc, err := flate.NewWriter(out, flate.DefaultCompression)
	if err != nil {
		return nil, errs.NewError(err, "error creating encoder", false)
	}

	if _, err = io.Copy(enc, in); err != nil {
		enc.Close()
		return nil, errs.NewError(err, "error copying data", false)
	}

	enc.Close()

	return out.Bytes(), nil
}

func (c *FlateCompressor) Decompress(compressed []byte) (data []byte, err error) {
	in := bytes.NewReader(compressed)
	out := bytes.NewBuffer(make([]byte, 0, len(compressed)))
	dec := flate.NewReader(in)

	defer dec.Close()

	_, err = io.Copy(out, dec)
	if err != nil {
		return nil, errs.NewError(err, "error copying data", false)
	}

	return out.Bytes(), nil
}
