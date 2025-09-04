package mint

import (
	"github.com/cosmos/evm/precompiles/mint"
	"github.com/cosmos/evm/testutil/integration/evm/factory"
	"github.com/cosmos/evm/testutil/integration/evm/grpc"
	"github.com/cosmos/evm/testutil/integration/evm/network"
	testkeyring "github.com/cosmos/evm/testutil/keyring"
	"github.com/stretchr/testify/suite"
)

// PrecompileTestSuite is the implementation of the TestSuite interface for Mint precompile unit tests.
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

	precompile *mint.Precompile
}

func NewPrecompileTestSuite(create network.CreateEvmApp, options ...network.ConfigOption) *PrecompileTestSuite {
	return &PrecompileTestSuite{
		create:  create,
		options: options,
	}
}

func (s *PrecompileTestSuite) SetupTest() {
	keyring := testkeyring.New(3) // we'd need 3 keys: admin, user1, and user2
	options := []network.ConfigOption{
		network.WithPreFundedAccounts(keyring.GetAllAccAddrs()...),
	}
	options = append(options, s.options...)
	integrationNetwork := network.NewUnitTestNetwork(s.create, options...)
	grpcHandler := grpc.NewIntegrationHandler(integrationNetwork)
	txFactory := factory.New(integrationNetwork, grpcHandler) // todo: implement missing methods.

	ctx := integrationNetwork.GetContext()
	sk := integrationNetwork.App.GetStakingKeeper()
	bondDenom, err := sk.BondDenom(ctx)
	s.Require().NoError(err)
	s.Require().NotEmpty(bondDenom, "bond denom cannot be empty")

	s.bondDenom = bondDenom
	s.factory = txFactory
	s.grpcHandler = grpcHandler
	s.keyring = keyring
	s.network = integrationNetwork

	// Set the authority as the first keyring address (admin)
	s.authority = keyring.GetKey(0).Addr.Hex()

	// Create the mint precompile instance
	s.precompile, err = mint.NewPrecompile(
		s.authority,
		s.network.App.GetBankKeeper(),
	)
	s.Require().NoError(err, "failed to create mint precomile")
}
