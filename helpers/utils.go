package helpers

import (
	"bytes"
	"fmt"
	"unicode"
)

func IsNumber(slugOrID *string) bool {
	for _, char := range *slugOrID {
		if !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}

func lsdBuilder(queryStringBuffer *bytes.Buffer, limit, since, desc []byte, sinceString string, equality bool) (sinceExists bool) {
	faseDescChecker := false
	if len(since) != 0 {
		sinceExists = true
		sign := ">"
		if desc != nil && bytes.Equal([]byte("true"), desc) {
			faseDescChecker = true
			sign = "<"
		}
		if equality {
			sign += "="
		}
		queryStringBuffer.WriteString(fmt.Sprintf(" AND %s %s $2", sinceString, sign))
	}

	queryStringBuffer.WriteString(fmt.Sprintf("\nORDER BY %s", sinceString))
	if faseDescChecker || desc != nil && bytes.Equal([]byte("true"), desc) {
		queryStringBuffer.WriteString(" DESC")
	}

	if limit != nil {
		queryStringBuffer.WriteString(fmt.Sprintf("\nLIMIT %s", limit))
	}
	return
}
