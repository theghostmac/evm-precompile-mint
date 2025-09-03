package mint

import (
	"testing"

	"github.com/cosmos/evm/evmd/tests/integration"
	"github.com/cosmos/evm/tests/integration/precompiles/mint"
	"github.com/stretchr/testify/suite"
)

func TestMintPrecompileTestSuite(t *testing.T) {
	s := mint.NewPrecompileTestSuite(integration.CreateEvmd)
	suite.Run(t, s)
}
