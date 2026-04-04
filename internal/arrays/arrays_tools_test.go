package arrays

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetArraysDiffs(t *testing.T) {
	t.Parallel()

	a := []int{1, 2, 3, 4, 5}
	b := []int{1, 2, 3, 7, 8, 9, 10}

	t.Run("Find difference between two int arrays", func(t *testing.T) {
		t.Parallel()

		diffs := GetArraysDiffs(a, b)

		want := []string{"-4", "-5", "+7", "+8", "+9", "+10"}
		require.Equal(t, want, diffs, "GetArraysDiffs() should return the correct diffs")
	})
}
