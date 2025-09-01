package mint

import (
	"embed"

	storetypes "cosmossdk.io/store/types"
	cmn "github.com/cosmos/evm/precompiles/common"
	"github.com/ethereum/go-ethereum/common"
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

var f embed.FS

var _ vm.PrecompiledContract = &Precompile{}

// NewPRecompile creates a new Mint Precompile instance as a PrecompiledContract interface.
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
	p.SetAddress(common.HexToAddress("0x0000000000000000000000000000000000001111"))

	return p, nil
}

func (p Precompile) RequiredGas(input []byte) uint64 {
	//TODO implement me
	panic("implement me")
}

func (p Precompile) Run(evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}
