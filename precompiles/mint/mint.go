package mint

import (
	"embed"
	"fmt"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cmn "github.com/cosmos/evm/precompiles/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/vm"
)

// Precompile defines the precompiled contract for Mint.
type Precompile struct {
	cmn.Precompile
	authority  string
	bankKeeper cmn.BankKeeper
}

const (
	abiPath = "abi.json"
	GasMint = 25_000
)

// Embed abi json file to the executable binary.
//
//go:embed abi.json
var f embed.FS

var _ vm.PrecompiledContract = &Precompile{}

// NewPrecompile creates a new Mint Precompile instance as a PrecompiledContract interface.
func NewPrecompile(authority string, bankKeeper cmn.BankKeeper) (*Precompile, error) {
	newABI, err := cmn.LoadABI(f, abiPath)
	if err != nil {
		return nil, err
	}

	p := &Precompile{
		Precompile: cmn.Precompile{
			ABI:                  newABI,
			KvGasConfig:          storetypes.GasConfig{},
			TransientKVGasConfig: storetypes.GasConfig{},
		},
		authority:  authority,
		bankKeeper: bankKeeper,
	}

	// Set the address of the Mint precompile contract. would use 0x1111 for easy testing
	// TODO:
	p.SetAddress(common.HexToAddress("0x0000000000000000000000000000000000001111"))

	return p, nil
}

// RequiredGas calculates the contract gas used.
func (p *Precompile) RequiredGas(input []byte) uint64 {
	if len(input) < 4 {
		return 0
	}

	methodID := input[:4]
	method, err := p.MethodById(methodID)
	if err != nil {
		return 0
	}

	switch method.Name {
	case MintMethod:
		return GasMint
	default:
		return 0
	}
}

// Run executes the precompiled contract Mint method defined in the ABI.
func (p *Precompile) Run(evm *vm.EVM, contract *vm.Contract, readonly bool) (bz []byte, err error) {
	bz, err = p.run(evm, contract, readonly)
	if err != nil {
		return cmn.ReturnRevertError(evm, err)
	}

	return bz, nil
}

func (p *Precompile) run(evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	// Mint precompile cannot receive funds.
	if value := contract.Value(); value.Sign() == 1 {
		return nil, fmt.Errorf(ErrCannotReceiveFunds, contract.Value().String())
	}

	ctx, stateDB, method, initialGas, args, err := p.RunSetup(evm, contract, readonly, p.IsTransaction)
	if err != nil {
		return nil, err
	}

	// If there's any error so far during execution, don't panic, just return OOG error so EVM can continue running.
	defer cmn.HandleGasError(ctx, contract, initialGas, &err)()

	bz, err := p.HandleMethod(ctx, contract, stateDB, method, args)
	if err != nil {
		return nil, err
	}

	cost := ctx.GasMeter().GasConsumed() - initialGas

	if !contract.UseGas(cost, nil, tracing.GasChangeCallPrecompiledContract) {
		return nil, vm.ErrOutOfGas
	}

	return bz, nil
}

// IsTransaction checks if the given method name corresponds to a transaction or query.
func (p *Precompile) IsTransaction(method *abi.Method) bool {
	switch method.Name {
	case MintMethod:
		return true
	default:
		return false
	}
}

// HandleMethod handles the execution of each of the Mint methods.
func (p *Precompile) HandleMethod(
	ctx sdk.Context,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	args []interface{},
) (bz []byte, err error) {
	switch method.Name {
	case MintMethod:
		bz, err = p.Mint(ctx, contract, stateDB, method, args)
	default:
		return nil, fmt.Errorf(cmn.ErrUnknownMethod, method.Name)
	}

	return bz, err
}

// GetAuthority returns the authority address for testing purposes
func (p *Precompile) GetAuthority() string {
	return p.authority
}
