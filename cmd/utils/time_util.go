package utils

import "time"

func NowStr() string {
	return time.Now().Format("2006-Jan-02 15:04:05")
}
