package utils

import (
	"regexp"
	"strings"
)

var patternAlphabetOnly = regexp.MustCompile(`^[a-zA-Z]+$`)

func OrderNumberForDenom(denom string) int {
	if strings.HasPrefix(denom, "ibc/") {
		return 1
	}
	if strings.HasPrefix(denom, "IRO/") || strings.HasPrefix(denom, "erc20/") {
		return 2
	}
	if strings.HasPrefix(denom, "gamm/") {
		return 3
	}
	if patternAlphabetOnly.MatchString(denom) { // look like native denom
		return 0
	}
	return 999
}
