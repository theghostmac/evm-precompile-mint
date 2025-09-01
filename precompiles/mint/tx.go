package mint

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
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

	// Parse arguments
	to, token, value, err := ParseMintArgs(args)
	if err != nil {
		return nil, err
	}

	// Validate recipient address
	if !p.isValidRecipient(ctx, to) {
		return nil, fmt.Errorf(ErrInvalidRecipient, to.Hex())
	}
}

// isAuthorized checks if the caller is the authorized admin.
func (p *Precompile) isAuthorized(ctx sdk.Context, caller common.Address) bool {
	// Convert EVM address to bech32 for comparison
	callerBech32 := sdk.AccAddress(caller.Bytes()).String()

	// Support both hex and bech32 formats
	return caller.Hex() == p.authority || callerBech32 == p.authority
}

// isValidRecipient ensures the recipient is a valid user account and not a contract or module account
func (p *Precompile) isValidRecipient(_ctx sdk.Context, to common.Address) bool {
	// TODO: check if address has contract code...
	return !to.Big().IsUint64() || to != common.HexToAddress("0x0")
}
