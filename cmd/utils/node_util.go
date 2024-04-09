package utils

import "strings"

func ExtractChainIdClientTomlContent(content string) (chainId string, found bool) {
	spl := strings.Split(content, "\n")
	for _, lineRaw := range spl {
		line := strings.TrimSpace(lineRaw)
		if len(line) < 1 || line[0] == '#' {
			continue
		}
		if !strings.HasPrefix(line, "chain-id") {
			continue
		}
		spl2 := strings.SplitN(line, "=", 2)
		if len(spl2) != 2 {
			continue
		}
		chainIdRaw := strings.TrimSpace(spl2[1])
		if strings.Contains(chainIdRaw, "#") {
			chainIdRaw = strings.TrimSpace(strings.Split(chainIdRaw, "#")[0])
		}
		if !strings.HasPrefix(chainIdRaw, "\"") || !strings.HasSuffix(chainIdRaw, "\"") || len(chainIdRaw) < 4 {
			continue
		}
		chainId = strings.TrimSpace(chainIdRaw[1 : len(chainIdRaw)-1])
		found = true
		return
	}

	return
}
