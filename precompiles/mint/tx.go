package mint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/vm"
)

const (
	// MintMethod defines the ABI method name for the mint transaction
	MintMethod = "mint"
)

// Mint mints native tokens to the specified address
func (p *Precompile) Mint(ctx sdk.Context, contract *vm.Contract, stateDB vm.StateDB, method *abi.Method, args []interface{}) ([]byte, error) {
	// First check if the caller is authorized
	caller := contract.Caller()
	if !p.isAuthorized(ctx, caller) {
		return nil, ErrUnauthorized
	}
}
