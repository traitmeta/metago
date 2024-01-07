package common

import "strings"

func Has0xPrefix(str string) bool {
	if len(str) < 2 {
		return false
	}
	return strings.ToLower(str[:2]) == Prefix0x
}
