package utils

import (
	"bufio"
	"fmt"
	"strings"
)

func ReadYesNo(reader *bufio.Reader) (yes bool, err error) {
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(strings.ToLower(text))

	switch text {
	case "y":
		yes = true
		break
	case "yes":
		yes = true
		break
	case "n":
		break
	case "no":
		break
	default:
		err = fmt.Errorf("'%s' is not an accepted answer!\nYour answer must be Yes/No (or Y/N)", text)
		break
	}
	return
}
