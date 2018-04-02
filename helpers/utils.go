package helpers

import "unicode"

func IsNumber(slugOrID *string) bool {
	for _, char := range *slugOrID {
		if !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}
