package geoscape_test

import (
	"encoding/hex"
	"strings"
	"unicode"
)

func loadHex(str string) ([]byte, error) {
	txt := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			// if the character is a space, drop it
			return -1
		}
		// else keep it in the string
		return r
	}, str)
	return hex.DecodeString(txt)
}
