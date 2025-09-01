package mint

import "errors"

var (
	ErrUnauthorized = errors.New("caller is not authorized to mint tokens")
)
