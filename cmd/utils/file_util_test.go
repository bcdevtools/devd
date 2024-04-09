package utils

import (
	"fmt"
	"io/fs"
	"testing"
)

func TestValidatePasswordFileMode(t *testing.T) {
	tests := []struct {
		mode    int
		wantErr bool
	}{
		{
			mode:    0o200,
			wantErr: true,
		},
		{
			mode:    0o220,
			wantErr: true,
		},
		{
			mode:    0o222,
			wantErr: true,
		},
		{
			mode:    0o202,
			wantErr: true,
		},
		{
			mode:    0o400,
			wantErr: false,
		},
		{
			mode:    0o440,
			wantErr: true,
		},
		{
			mode:    0o444,
			wantErr: true,
		},
		{
			mode:    0o404,
			wantErr: true,
		},
		{
			mode:    0o500,
			wantErr: false,
		},
		{
			mode:    0o550,
			wantErr: true,
		},
		{
			mode:    0o555,
			wantErr: true,
		},
		{
			mode:    0o505,
			wantErr: true,
		},
		{
			mode:    0o600,
			wantErr: false,
		},
		{
			mode:    0o660,
			wantErr: true,
		},
		{
			mode:    0o666,
			wantErr: true,
		},
		{
			mode:    0o606,
			wantErr: true,
		},
		{
			mode:    0o700,
			wantErr: false,
		},
		{
			mode:    0o770,
			wantErr: true,
		},
		{
			mode:    0o777,
			wantErr: true,
		},
		{
			mode:    0o707,
			wantErr: true,
		},
		{
			mode:    0o411,
			wantErr: true,
		},
		{
			mode:    0o640,
			wantErr: true,
		},
		{
			mode:    0o750,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%o", tt.mode), func(t *testing.T) {
			if err := ValidatePasswordFileMode(fs.FileMode(tt.mode)); (err != nil) != tt.wantErr {
				t.Errorf("ValidatePasswordFileMode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExtractPermissionParts(t *testing.T) {
	tests := []struct {
		perm      fs.FileMode
		wantOwner int
		wantGroup int
		wantOther int
	}{
		{
			perm:      0o750,
			wantOwner: 7,
			wantGroup: 5,
			wantOther: 0,
		},
		{
			perm:      0o756,
			wantOwner: 7,
			wantGroup: 5,
			wantOther: 6,
		},
		{
			perm:      0o777,
			wantOwner: 7,
			wantGroup: 7,
			wantOther: 7,
		},
		{
			perm:      0o600,
			wantOwner: 6,
			wantGroup: 0,
			wantOther: 0,
		},
		{
			perm:      0o060,
			wantOwner: 0,
			wantGroup: 6,
			wantOther: 0,
		},
		{
			perm:      0o066,
			wantOwner: 0,
			wantGroup: 6,
			wantOther: 6,
		},
		{
			perm:      0o123,
			wantOwner: 1,
			wantGroup: 2,
			wantOther: 3,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%o", tt.perm), func(t *testing.T) {
			gotOwner, gotGroup, gotOther := ExtractPermissionParts(tt.perm)
			if gotOwner != tt.wantOwner {
				t.Errorf("ExtractPermissionParts() gotOwner = %v, want %v", gotOwner, tt.wantOwner)
			}
			if gotGroup != tt.wantGroup {
				t.Errorf("ExtractPermissionParts() gotGroup = %v, want %v", gotGroup, tt.wantGroup)
			}
			if gotOther != tt.wantOther {
				t.Errorf("ExtractPermissionParts() gotOther = %v, want %v", gotOther, tt.wantOther)
			}
		})
	}
}
