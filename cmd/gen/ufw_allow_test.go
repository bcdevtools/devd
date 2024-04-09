package gen

import "testing"

func Test_isValidFirewallIpAddress(t *testing.T) {
	tests := []struct {
		ip   string
		want bool
	}{
		{
			ip:   ":::1",
			want: false,
		},
		{
			ip:   "127.0.0.1",
			want: true,
		},
		{
			ip:   "localhost",
			want: false,
		},
		{
			ip:   "127.0.0.1/0",
			want: false,
		},
		{
			ip:   "127.0.0.1/8",
			want: true,
		},
		{
			ip:   "127.0.0.1/16",
			want: true,
		},
		{
			ip:   "127.0.0.1/24",
			want: true,
		},
		{
			ip:   "127.0.0.1/32",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.ip, func(t *testing.T) {
			if got := isValidFirewallIpAddress(tt.ip); got != tt.want {
				t.Errorf("isValidFirewallIpAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
