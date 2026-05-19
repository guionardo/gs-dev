package compression_test

import (
	"testing"

	lorem "github.com/derektata/lorem/ipsum"
	"github.com/guionardo/gs-dev/internal/compression"
	"github.com/stretchr/testify/require"
)

var (
	reallySmallData, smallData, longData []byte
)

func init() {
	// Create a new generator
	g := lorem.NewGenerator()
	g.WordsPerSentence = 10     // Customize how many words per sentence
	g.SentencesPerParagraph = 5 // Customize how many sentences per paragraph
	g.CommaAddChance = 3        // Customize the chance of a comma being added to a sentence

	reallySmallString := g.Generate(5)
	smallString := g.Generate(200)
	longString := g.Generate(5000)

	reallySmallData = []byte(reallySmallString)
	smallData = []byte(smallString)
	longData = []byte(longString)
}

func TestCompression(t *testing.T) {
	t.Parallel()

	t.Run("compress_small_data_should_return_dummy_compressor", func(t *testing.T) {
		t.Parallel()

		compressed, compressor := compression.Compress(reallySmallData)

		require.NotNil(t, compressed)
		require.Equal(t, compression.CompressorDummy, compressor)

		decompressed, err := compression.Decompress(compressed, compressor)
		require.NoError(t, err)
		require.NotNil(t, decompressed)
		require.Equal(t, reallySmallData, decompressed)
	})

	t.Run("compress_small_data_should_return_flates_compressor", func(t *testing.T) {
		t.Parallel()

		compressed, compressor := compression.Compress(smallData)

		require.NotNil(t, compressed)
		require.Equal(t, compression.CompressorFlate, compressor)

		decompressed, err := compression.Decompress(compressed, compressor)
		require.NoError(t, err)
		require.NotNil(t, decompressed)
		require.Equal(t, smallData, decompressed)
	})

	t.Run("compress_long_data_should_return_flates_compressor", func(t *testing.T) {
		t.Parallel()

		compressed, compressor := compression.Compress(longData)

		require.NotNil(t, compressed)
		require.Equal(t, compression.CompressorFlate, compressor)

		decompressed, err := compression.Decompress(compressed, compressor)
		require.NoError(t, err)
		require.NotNil(t, decompressed)
		require.Equal(t, longData, decompressed)
	})
}
