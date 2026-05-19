package openurl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_checkReachableUrl(t *testing.T) {
	t.Parallel()

	t.Run("checkReachableUrl_with_valid_url_should_return_nil", func(t *testing.T) {
		t.Parallel()

		url := "https://google.com"
		err := checkReachableUrl(url)
		require.NoError(t, err)
	})

	t.Run("checkReachableUrl_with_invalid_url_should_return_error", func(t *testing.T) {
		t.Parallel()

		url := "https://notfound.0000.11111.unexistent.com"
		err := checkReachableUrl(url)
		require.Error(t, err)
	})

	t.Run("checkReachableUrl_with_unsupported_url_scheme_should_return_error", func(t *testing.T) {
		t.Parallel()

		url := "ftp://google.com"
		err := checkReachableUrl(url)
		require.Error(t, err)
	})

	t.Run("checkReachableUrl_with_invalid_url_should_return_error", func(t *testing.T) {
		t.Parallel()

		url := "invalid-url"
		err := checkReachableUrl(url)
		require.Error(t, err)
	})

	t.Run("checkReachableUrl_with_unexistent_url_should_return_error", func(t *testing.T) {
		t.Parallel()

		url := "https://github.com/guionardo/unexistent"
		err := checkReachableUrl(url)
		require.Error(t, err)
	})
}
