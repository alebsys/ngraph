package collector

import (
	"testing"
)

func Test_checkConnectDirection(t *testing.T) {
	tests := []struct {
		name string
		port int
		min  int
		max  int
		got  string
		want string
	}{
		{
			name: "port within range, expect 'input'",
			port: 444,
			min:  32768,
			max:  65535,
			want: "input",
		},
		{
			name: "port within range, expect 'output'",
			port: 33445,
			min:  32768,
			max:  65535,
			want: "output",
		},
		{
			name: "port at the lower bound, expect 'output'",
			port: 32768,
			min:  32768,
			max:  65535,
			want: "output",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checkConnectDirection(tt.port, tt.min, tt.max)
			if got != tt.want {
				t.Errorf("\nExpected '%s', got '%s' (port: %d, min: %d, max: %d)", tt.want, got, tt.port, tt.min, tt.max)
			}
		})
	}
}
