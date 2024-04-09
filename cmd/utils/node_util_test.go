package utils

import (
	"strconv"
	"testing"
)

func TestExtractChainIdClientTomlContent(t *testing.T) {
	tests := []struct {
		content     string
		wantChainId string
		wantFound   bool
	}{
		{
			content: `
# This is a comment
chain-id = "test-chain-id"
`,
			wantChainId: "test-chain-id",
			wantFound:   true,
		},
		{
			content: `
# This is a comment
chain-id= "test-chain-id"
# This is a comment
`,
			wantChainId: "test-chain-id",
			wantFound:   true,
		},
		{
			content: `
# This is a comment
chain-id ="test-chain-id"
# This is a comment
`,
			wantChainId: "test-chain-id",
			wantFound:   true,
		},
		{
			content: `
# This is a comment
chain-id ="test-chain-id-1"#comment
chain-id =  "test-chain-id-2"
# This is a comment
`,
			wantChainId: "test-chain-id-1",
			wantFound:   true,
		},
		{
			content: `
# This is a comment
	# chain-id ="test-chain-id-1"
	chain-id =  "test-chain-id-2"
# This is a comment
`,
			wantChainId: "test-chain-id-2",
			wantFound:   true,
		},
		{
			content: `
# This is a comment
	# chain-id ="test-chain-id-1"
	chain-id =  "test-chain-id-2" # This is a comment
# This is a comment
`,
			wantChainId: "test-chain-id-2",
			wantFound:   true,
		},
		{
			content: `
# This is a comment
	# chain-id ="test-chain-id-1"
	chain-id =  "2" # This is a comment
# This is a comment
`,
			wantFound: false,
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			gotChainId, gotFound := ExtractChainIdClientTomlContent(tt.content)
			if gotChainId != tt.wantChainId {
				t.Errorf("ExtractChainIdClientTomlContent() gotChainId = %v, want %v", gotChainId, tt.wantChainId)
			} else if gotFound != tt.wantFound {
				t.Errorf("ExtractChainIdClientTomlContent() gotFound = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}
