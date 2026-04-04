package urlparsers_test

import (
	"testing"

	urlparsers "github.com/guionardo/gs-dev/internal/git/url_parsers"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		url     string
		want    string
		wantErr bool
	}{
		{
			name: "git ssh",
			url:  "git@github.com:guionardo/gs-dev.git",
			want: "https://github.com/guionardo/gs-dev",
		},
		{
			name: "git https",
			url:  "https://github.com/guionardo/gs-dev.git",
			want: "https://github.com/guionardo/gs-dev",
		},
		{
			name: "git private",
			url:  "ssh://domain.com.br:22/Company/Project/_git/sample-project",
			want: "https://domain.com.br/Company/Project/_git/sample-project",
		},
		{
			name: "gitlab ssh",
			url:  "git@gitlab.com:wee-ops/wee-api.git",
			want: "https://gitlab.com/wee-ops/wee-api",
		},
		{
			name: "github ssh",
			url:  "git@github.com:guionardo/go-dev.git",
			want: "https://github.com/guionardo/go-dev",
		},
		{
			name: "azure ssh",
			url:  "git@ssh.dev.azure.com:v3/CUSTOMER-SA/CUSTOMER-NS/ms-credit-api",
			want: "https://dev.azure.com/CUSTOMER-SA/CUSTOMER-NS/_git/ms-credit-api",
		},
		{
			name: "gitlab http",
			url:  "https://gitlab.com/wee-ops/wee-api.git",
			want: "https://gitlab.com/wee-ops/wee-api",
		},
		{
			name: "github http",
			url:  "https://github.com/guionardo/go-dev.git",
			want: "https://github.com/guionardo/go-dev",
		},
		{
			name: "azure http",
			url:  "https://CUSTOMER-SA@dev.azure.com/CUSTOMER-SA/CUSTOMER-NS/_git/metric-api",
			want: "https://dev.azure.com/CUSTOMER-SA/CUSTOMER-NS/_git/metric-api",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, gotErr := urlparsers.Parse(tt.url)
			if tt.wantErr {
				require.Error(t, gotErr, "Parse() should return error")
			} else {
				require.NoError(t, gotErr, "Parse() should not return error")
			}

			require.Equal(t, tt.want, got, "Parse() should return the correct URL")
		})
	}
}
