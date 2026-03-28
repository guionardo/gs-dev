package git

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		url     string
		want    GitURL
		wantUrl string
	}{
		{
			name: "gitlab ssh",
			url:  "git@gitlab.com:wee-ops/wee-api.git",
			want: GitURL{
				Success: true,
				Domain:  "gitlab.com",
				Repo:    "wee-ops/wee-api",
			},
			wantUrl: "https://gitlab.com/wee-ops/wee-api",
		},
		{
			name: "github ssh",
			url:  "git@github.com:guionardo/go-dev.git",
			want: GitURL{
				Success: true,
				Domain:  "github.com",
				Repo:    "guionardo/go-dev",
			},
			wantUrl: "https://github.com/guionardo/go-dev",
		},
		{
			name: "azure ssh",
			url:  "git@ssh.dev.azure.com:v3/CUSTOMER-SA/CUSTOMER-NS/ms-credit-api",
			want: GitURL{
				Success: true,
				Domain:  "dev.azure.com/CUSTOMER-SA/CUSTOMER-NS",
				Repo:    "ms-credit-api",
			},
			wantUrl: "https://dev.azure.com/CUSTOMER-SA/CUSTOMER-NS/_git/ms-credit-api",
		},
		{
			name: "gitlab http",
			url:  "https://gitlab.com/wee-ops/wee-api.git",
			want: GitURL{
				Success: true,
				Domain:  "gitlab.com",
				Repo:    "wee-ops/wee-api",
			},
			wantUrl: "https://gitlab.com/wee-ops/wee-api",
		},
		{
			name: "github http",
			url:  "https://github.com/guionardo/go-dev.git",
			want: GitURL{
				Success: true,
				Domain:  "github.com",
				Repo:    "guionardo/go-dev",
			},
			wantUrl: "https://github.com/guionardo/go-dev",
		},
		{
			name: "azure http",
			url:  "https://CUSTOMER-SA@dev.azure.com/CUSTOMER-SA/CUSTOMER-NS/_git/metric-api",
			want: GitURL{
				Success: true,
				Domain:  "dev.azure.com/CUSTOMER-SA/CUSTOMER-NS",
				Repo:    "metric-api",
			},
			wantUrl: "https://dev.azure.com/CUSTOMER-SA/CUSTOMER-NS/_git/metric-api",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Parse(tt.url)
			if tt.want.Success {
				require.NoError(t, err, "Parse() should not return error")
			}

			require.Equal(t, tt.want.Success, got.Success, "Parse() success expected")
			require.Equal(t, tt.want.Domain, got.Domain, "Parse() domain expected")
			require.Equal(t, tt.want.Repo, got.Repo, "Parse() repo expected")
			require.Equal(t, tt.wantUrl, got.GetURL(), "Parse() url expected")
		})
	}
}
