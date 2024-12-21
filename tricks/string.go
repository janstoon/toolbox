package tricks

import "strings"

func StringToRunes(src string) []rune {
	return []rune(src)
}

func IsEmptyString(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}
