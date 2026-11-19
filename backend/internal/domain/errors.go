package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrNoUpdateParams = errors.New("no set parameters for partial update")
)

func ErrInvalidParam(paramName string) error {
	return fmt.Errorf("invalid parameter %s", paramName)
}
