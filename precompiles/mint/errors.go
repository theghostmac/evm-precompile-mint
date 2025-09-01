package mint

import "errors"

const (
	ErrInvalidRecipient = "invalid recipient address: %s"
)

var (
	ErrUnauthorized = errors.New("caller is not authorized to mint tokens")
)
