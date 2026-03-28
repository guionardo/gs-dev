package colors

import (
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Parallel()

	phrase := "{warning}test{normal} Resto da frase"
	color.NoColor = false
	parsedPhrase := Parse(phrase)
	t.Logf("Parsed phrase:\n%s", parsedPhrase)
	require.Equal(t, "\x1b[43;30mtest\x1b[0m\x1b[37;40m Resto da frase\x1b[0m", parsedPhrase)
}
