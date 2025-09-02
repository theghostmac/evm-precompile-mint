package mint

import (
	"github.com/cosmos/evm/precompiles/mint"
	"github.com/cosmos/evm/testutil/integration/base/grpc"
	"github.com/cosmos/evm/testutil/integration/evm/factory"
	"github.com/cosmos/evm/testutil/integration/evm/network"
	testkeyring "github.com/cosmos/evm/testutil/keyring"
	"github.com/stretchr/testify/suite"
)

// PrecomileTestSuite is the implementation of the TestSuite interface for Mint precompile unit tests.
type PrecompileTestSuite struct {
	suite.Suite

	create    network.CreateEvmApp
	options   []network.ConfigOption
	bondDenom string
	// authority is the address of the admin who can call the mint function.
	authority   string
	network     *network.UnitTestNetwork
	factory     factory.TxFactory
	grpcHandler grpc.Handler
	keyring     testkeyring.Keyring

	precomppile *mint.Precompile
}

func NewPrecompileTestSuite(create network.CreateEvmApp, options ...network.ConfigOption) *PrecompileTestSuite {
	return &PrecompileTestSuite{
		create:  create,
		options: options,
	}
}
