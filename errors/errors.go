package errors

import (
	"github.com/pkg/errors"
)

var NotFoundError = errors.New("Not found")
var ConflictError = errors.New("Conflict")
var EmptySearchError = errors.New("Search result is empty")
