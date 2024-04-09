package files

import "testing"

func Test_isOrContainsRsyncRecursiveFlag(t *testing.T) {
	//goland:noinspection SpellCheckingInspection
	tests := []struct {
		option string
		want   bool
	}{
		{
			option: "recursive",
			want:   false,
		},
		{
			option: "archive",
			want:   false,
		},
		{
			option: "r",
			want:   false,
		},
		{
			option: "a",
			want:   false,
		},
		{
			option: "--recursive",
			want:   true,
		},
		{
			option: "--archive",
			want:   true,
		},
		{
			option: "-r",
			want:   true,
		},
		{
			option: "-a",
			want:   true,
		},
		{
			option: "-av",
			want:   true,
		},
		{
			option: "-recursive",
			want:   true, // contains 'r' and that's enough
		},
		{
			option: "--r",
			want:   false,
		},
		{
			option: "--a",
			want:   false,
		},
		{
			option: "--recursive-mixed",
			want:   false,
		},
		{
			option: "--archive-mixed",
			want:   false,
		},
		{
			option: "-rlptgoD",
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.option, func(t *testing.T) {
			if got := isOrContainsRsyncRecursiveFlag(tt.option); got != tt.want {
				t.Errorf("isOrContainsRsyncRecursiveFlag() = %v, want %v", got, tt.want)
			}
		})
	}
}
