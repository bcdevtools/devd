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

func ConvertDisplayWithExponentIntoRaw(display string, exponent int, decimalsPoint rune) (number, highNumber, lowNumber *big.Int, err error) {
	var markupRune rune
	if decimalsPoint == '.' {
		markupRune = ','
	} else {
		markupRune = '.'
	}

	if exponent < 0 || exponent > 18 {
		err = fmt.Errorf("exponent must be in range 0 to 18, got %d", exponent)
		return
	}

	// remove markup runes
	display = strings.ReplaceAll(display, string(markupRune), "")

	spl := strings.Split(display, string(decimalsPoint))
	if len(spl) > 2 {
		err = fmt.Errorf("input contains multiple decimals points")
		return
	}

	var ok bool
	highNumber, ok = new(big.Int).SetString(spl[0], 10)
	if !ok {
		err = fmt.Errorf("failed to read, left part %s is not a number", spl[0])
		return
	}

	if len(spl) == 1 || strings.ReplaceAll(spl[1], "0", "") == "" {
		lowNumber = common.Big0
	} else {
		low := spl[1]
		for len(low) < exponent {
			low += "0"
		}

		lowNumber, ok = new(big.Int).SetString(low, 10)
		if !ok {
			err = fmt.Errorf("failed to read, right part %s is not a number", low)
			return
		}
	}

	oneHigh := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(exponent)), nil)
	number = new(big.Int).Add(new(big.Int).Mul(highNumber, oneHigh), lowNumber)

	return
}
