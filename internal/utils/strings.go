package utils

import "unicode"

func MixCase(s string) string {
	r := []rune(s)
	for i := range r {
		if i%2 == 0 {
			r[i] = unicode.ToUpper(r[i])
		} else {
			r[i] = unicode.ToLower(r[i])
		}
	}
	return string(r)
}
