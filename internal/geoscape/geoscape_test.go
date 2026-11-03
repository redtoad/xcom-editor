package geoscape_test

import (
	"encoding/hex"
	"strings"
	"unicode"
)

// MustLoadFromHex loads a file content from a hex-encoded string.
func MustLoadFromHex(str string) []byte {
	txt := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			// if the character is a space, drop it
			return -1
		}
		// else keep it in the string
		return r
	}, str)
	buffer, err := hex.DecodeString(txt)
	if err != nil {
		panic(err)
	}
	return buffer
}
