package mint

import (
	"fmt"

	"cosmossdk.io/math"
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

	// Convert amount to SDK Int
	amount := math.NewIntFromBigInt(value)
	if amount.IsZero() {
		return nil, ErrZeroAmount
	}
	if amount.IsNegative() {
		return nil, ErrNegativeAmount
	}

	// Create coin to mint
	coin := sdk.NewCoin(token, amount)
	coins := sdk.NewCoins(coin)

	// Validate coin denomination
	if err := coin.Validate(); err != nil {
		return nil, fmt.Errorf(ErrInvalidDenom, token)
	}

	// MINTING TOKENS USING THE BANK KEEPER

	// Mint coins to temporary module account.
	if err := p.bankKeeper.MintCoins(ctx, "mint", coins); err != nil {
		return nil, fmt.Errorf(ErrMintFailed, err.Error())
	}

	// Send minted coins to recipient
	recipientAddr := sdk.AccAddress(to.Bytes())
	if err := p.bankKeeper.SendCoinsFromModuleToAccount(ctx, "mint", recipientAddr, coins); err != nil {
		return nil, fmt.Errorf(ErrTransferFailed, err.Error())
	}

	// Emit mint event
	if err := p.EmitMintEvent(ctx, stateDB, to, token, value); err != nil {
		return nil, err
	}

	return method.Outputs.Pack()
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
