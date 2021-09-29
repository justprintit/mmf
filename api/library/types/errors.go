package types

import (
	"go.sancus.dev/core/errors"
)

func ErrInvalidValue(v interface{}) error {
	return errors.New("Invalid Value: %v", v)
}
