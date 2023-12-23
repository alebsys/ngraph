package localnetinfo

import (
	"path/filepath"
	"reflect"
	"testing"
)

func Test_GetPortRange(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		want    LocalPortRange
		wantErr bool
	}{
		{
			name: "file found, no error should come up",
			file: filepath.Join("..", "..", "testdata", "ip_local_port_range"),
			want: LocalPortRange{
				MinPort: 32768,
				MaxPort: 60999,
			},
			wantErr: false,
		},
		{
			name:    "error case - file not found",
			file:    "somewhere over the rainbow",
			want:    LocalPortRange{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPortRange(tt.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPortRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPortRange() = %v, want %v", got, tt.want)
			}
		})
	}
}
