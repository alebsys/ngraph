package collector

import (
	"testing"

	"github.com/prometheus/procfs"
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

func Test_shouldMatchBySubnets(t *testing.T) {
	tests := []struct {
		name    string
		subnets []string
		addr    string
		want    bool
	}{
		{
			name:    "Matching subnet, should return true",
			subnets: []string{"none", "10.32.68"},
			addr:    "10.32.68.1",
			want:    true,
		},
		{
			name:    "Non-matching subnet, should return false",
			subnets: []string{"10.32.68"},
			addr:    "192.168.1.1",
			want:    false,
		},
		{
			name:    "TODO",
			subnets: []string{"", "10.32.68"},
			addr:    "192.168.1.1",
			want:    false,
		},
		{
			name:    "TODO",
			subnets: []string{" ", "nOne", "10.32"},
			addr:    "10.32.68.10",
			want:    true,
		},
		{
			name:    "TODO",
			subnets: []string{""},
			addr:    "10.32.68.10",
			want:    false,
		},
		{
			name:    "TODO",
			subnets: []string{" "},
			addr:    "10.32.68.10",
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldMatchBySubnets(tt.subnets, tt.addr)
			if got != tt.want {
				t.Errorf("\nExpected '%v', got '%v' (subnets: %v, address %s)", tt.want, got, tt.subnets, tt.addr)
			}
		})
	}
}

func Test_selectNetworkNamespaceInode(t *testing.T) {
	tests := []struct {
		name       string
		namespaces procfs.Namespaces
		want       uint32
	}{
		{
			name: "success",
			namespaces: procfs.Namespaces{
				"mnt": {Type: "mnt", Inode: 4026531840},
				"net": {Type: "net", Inode: 4026531991},
			},
			want: 4026531991,
		},
		{
			name: "error",
			namespaces: procfs.Namespaces{
				"mnt": {Type: "mnt", Inode: 4026531840},
				"ipc": {Type: "ipc", Inode: 4026531993},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := selectNetworkNamespaceInode(tt.namespaces)
			if got != tt.want {
				t.Errorf("Expected '%v', got '%v'", tt.want, got)
			}
		})
	}
}
