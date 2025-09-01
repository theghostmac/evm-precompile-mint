package mint

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// EventMint defines the event data for the Mint event
type EventMint struct {
	To    common.Address
	Token string
	Value *big.Int
}

// ParseMintArgs parses the arguments from the mint method and returns
// the recipient address, token denomination, and amount.
func ParseMintArgs(args []interface{}) (
	to common.Address, token string, value *big.Int, err error,
) {
	if len(args) != 3 {
		return common.Address{}, "", nil, fmt.Errorf("invalid number of arguments; expected 3; got: %d", len(args))
	}

	// Parse "to" address
	to, ok := args[0].(common.Address)
	if !ok {
		return common.Address{}, "", nil, fmt.Errorf("invalid `to` address: %v", args[0])
	}

	// Parse 'token' string
	token, ok = args[1].(string)
	if !ok {
		return common.Address{}, "", nil, fmt.Errorf("invalid `token`: %v", args[1])
	}

	// And finally parse 'value' big.Int
	value, ok = args[2].(*big.Int)
	if !ok {
		return common.Address{}, "", nil, fmt.Errorf("invalid `value`: %v", args[2])
	}

	return to, token, value, nil
}
