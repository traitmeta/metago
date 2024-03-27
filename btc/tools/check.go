package tools

import "strings"

func DustCheck(address string) int64 {
	if strings.HasPrefix(address, "1") {
		return 546
	} else if strings.HasPrefix(address, "3") {
		return 540
	} else if strings.HasPrefix(address, "bc1q") {
		return 294
	} else if strings.HasPrefix(address, "bc1p") {
		return 330
	}

	return 546
}
