package compression

type Compressor interface {
	IsCandidate(data []byte) bool
	Compress(data []byte) (compressed []byte, err error)
	Decompress(compressed []byte) (data []byte, err error)
}

const (
	CompressorDummy = "dummy"
	CompressorFlate = "flate"
)

func Compress(data []byte) (compressed []byte, compressor string) {
	compressed, err := NewFlateCompressor().Compress(data)
	if err == nil && len(compressed) < len(data) {
		return compressed, CompressorFlate
	}

	return data, CompressorDummy
}

func Decompress(compressed []byte, compressorName string) (data []byte, err error) {
	if compressorName == CompressorFlate {
		return NewFlateCompressor().Decompress(compressed)
	}

	return compressed, nil
}
