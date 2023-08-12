package url

import "testing"

func Test_checkReachableUrl(t *testing.T) {

	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"Valid URL", "https://google.com", false},
		{"Invalid URL", "https://notfound.0000.11111.unexistent.com", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkReachableUrl(tt.url); (err != nil) != tt.wantErr {
				t.Errorf("checkReachableUrl() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
