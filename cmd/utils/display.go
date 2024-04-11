package utils

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"strings"
)

func ConvertNumberIntoDisplayWithExponent(number *big.Int, exponent int) (display string, highNumber, lowNumber *big.Int, err error) {
	if exponent < 0 || exponent > 18 {
		err = fmt.Errorf("exponent must be in range 0 to 18, got %d", exponent)
		return
	}

	if number.Sign() < 0 {
		err = fmt.Errorf("number must be positive, got %s", number)
		return
	}

	if number.Sign() == 0 || exponent == 0 {
		highNumber = number
		lowNumber = common.Big0
		display = number.String() + ".0"
		return
	}

	oneHigh := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(exponent)), nil)

	highNumber = new(big.Int).Div(number, oneHigh)
	lowNumber = new(big.Int).Mod(number, oneHigh)

	if lowNumber.Sign() == 0 {
		display = highNumber.String() + ".0"
		return
	}

	displayLow := lowNumber.String()
	for len(displayLow) < exponent {
		displayLow = "0" + displayLow
	}

	displayLow = strings.TrimRightFunc(displayLow, func(r rune) bool {
		return r == '0'
	})

	display = fmt.Sprintf("%s.%s", highNumber, displayLow)
	return
}
