package utils

import "testing"

func TestTryExtractUserNameFromHomeDir(t *testing.T) {
	tests := []struct {
		homeDir      string
		wantUsername string
		wantSuccess  bool
	}{
		{
			homeDir:      "/root",
			wantUsername: "root",
			wantSuccess:  true,
		},
		{
			homeDir:      "/root/",
			wantUsername: "root",
			wantSuccess:  true,
		},
		{
			homeDir:      "/ro0t/",
			wantUsername: "",
			wantSuccess:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.homeDir, func(t *testing.T) {
			gotUsername, gotSuccess := TryExtractUserNameFromHomeDir(tt.homeDir)
			if gotUsername != tt.wantUsername {
				t.Errorf("TryExtractUserNameFromHomeDir() gotUsername = %v, want %v", gotUsername, tt.wantUsername)
			}
			if gotSuccess != tt.wantSuccess {
				t.Errorf("TryExtractUserNameFromHomeDir() gotSuccess = %v, want %v", gotSuccess, tt.wantSuccess)
			}
		})
	}
}
