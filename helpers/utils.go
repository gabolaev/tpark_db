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

func lsdBuilder(queryStringBuffer *bytes.Buffer, limit, since, desc []byte, sinceField, orderField string, equality bool) (sinceExists bool) {
	faseDescChecker := false
	if len(since) != 0 {
		sinceExists = true
		var sign string
		if desc != nil && bytes.Equal([]byte("true"), desc) {
			faseDescChecker = true
			sign = "<"
		} else {
			sign = ">"
		}
		if equality {
			sign += "="
		}
		queryStringBuffer.WriteString(fmt.Sprintf(" AND %s %s $2", sinceField, sign))
	}

	queryStringBuffer.WriteString(fmt.Sprintf("\nORDER BY %s", orderField))
	if faseDescChecker || desc != nil && bytes.Equal([]byte("true"), desc) {
		queryStringBuffer.WriteString(" DESC")
	}

	var strLimit = string(limit)
	if limit != nil && IsNumber(&strLimit) {
		queryStringBuffer.WriteString(fmt.Sprintf("\nLIMIT %s", limit))
	}
	return
}
