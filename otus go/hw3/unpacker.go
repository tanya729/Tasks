package hw3

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Unpack will return unpacked string or error from packed string
func Unpack(s string) (string, error) {
	var unpack strings.Builder
	var flagEscape bool
	for i, value := range s {
		switch {
		case string(value) == `\` && !flagEscape:
			flagEscape = true
		case flagEscape:
			unpack.WriteRune(value)
			flagEscape = false
		case unicode.IsDigit(value) && len(unpack.String()) > 0:
			number, _ := strconv.Atoi(string(value))
			if number == 0 {
				str := []rune(unpack.String())
				str = str[:len(str)-1]
				unpack.Reset()
				unpack.WriteString(string(str))
			} else {
				w := []rune(unpack.String())
				w = w[len(w)-1:]
				unpack.WriteString(strings.Repeat(string(w), number-1))
			}
		case !unicode.IsDigit(value):
			unpack.WriteRune(value)
		default:
			return "", errors.New(fmt.Sprintf("Something strange with unpacking string '%s' on character - %s, position - %d.", s, string(value), i))
		}
	}
	return unpack.String(), nil
}
