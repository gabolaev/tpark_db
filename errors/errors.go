package errors

import (
	"github.com/pkg/errors"
)

var NotFoundError = errors.New("Not found")
var ConflictError = errors.New("Conflict")
var EmptySearchError = errors.New("Search result is empty")
var WrongParamsError = errors.New("Wrong params error")
var NothingChangedError = errors.New("Nothing changed")
