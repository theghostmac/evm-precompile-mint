package mint

import "errors"

const (
	ErrCannotReceiveFunds = "cannot receive funds, received: %s"
	ErrInvalidRecipient   = "invalid recipient address: %s"
	ErrInvalidDenom       = "invalid token denomination: %s"
	ErrMintFailed         = "failed to mint tokens: %s"
	ErrTransferFailed     = "failed to transfer tokens: %s"
)

var (
	ErrUnauthorized   = errors.New("caller is not authorized to mint tokens")
	ErrZeroAmount     = errors.New("cannot mint zero amount of tokens")
	ErrNegativeAmount = errors.New("cannot mint negative amount of tokens")
)
