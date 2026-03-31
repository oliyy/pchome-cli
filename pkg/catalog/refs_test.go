package catalog

import "testing"

func TestNormalizeProductRef(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "raw id", input: "DRAA5K-A900JOK9O", want: "DRAA5K-A900JOK9O"},
		{name: "suffixed id", input: "DRAA5K-A900JOK9O-000", want: "DRAA5K-A900JOK9O"},
		{name: "url", input: "https://24h.pchome.com.tw/prod/DRAA5K-A900JOK9O?foo=bar", want: "DRAA5K-A900JOK9O"},
		{name: "invalid", input: "not-a-product", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NormalizeProductRef(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error for %q", tc.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tc.input, err)
			}
			if got != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, got)
			}
		})
	}
}
