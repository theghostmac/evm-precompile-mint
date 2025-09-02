package mint

import (
	"math/big"

	cmn "github.com/cosmos/evm/precompiles/common"
	"github.com/cosmos/evm/precompiles/mint"
	utiltx "github.com/cosmos/evm/testutil/tx"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func (s *PrecompileTestSuite) TestEmitMintEvent() {
	testcases := []struct {
		name  string
		to    common.Address
		token string
		value *big.Int
	}{
		{
			name:  "pass - basic mint event",
			to:    utiltx.GenerateAddress(),
			token: "umint",
			value: big.NewInt(1000),
		},
		{
			name:  "pass - large amount",
			to:    utiltx.GenerateAddress(),
			token: "uatom",
			value: big.NewInt(1000000000000),
		},
		{
			name:  "pass - different token",
			to:    utiltx.GenerateAddress(),
			token: "stake",
			value: big.NewInt(500),
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			s.SetupTest()
			stateDB := s.network.GetStateDB()

			err := s.precompile.EmitMintEvent(
				s.network.GetContext(), stateDB, tc.to, tc.token, tc.value,
			)
			s.Require().NoError(err, "expected mint event to be emitted successfully")

			log := stateDB.Logs()[0]
			s.Require().Equal(log.Address, s.precompile.Address())

			// Check that the event signature matches the one emitted
			event := s.precompile.Events[mint.EventTypeMint]
			s.Require().Equal(crypto.Keccak256Hash([]byte(event.Sig)), common.HexToHash(log.Topics[0].Hex()))
			s.Require().Equal(log.BlockNumber, uint64(s.network.GetContext().BlockHeight()))

			// Check that the  unpacked event matches the one emitted.
			var mintEvent mint.EventMint
			err = cmn.UnpackLog(s.precompile.ABI, &mintEvent, mint.EventTypeMint, *log)
			s.Require().NoError(err, "unable to unpack log into mint event")

			s.Require().Equal(tc.to, mintEvent.To, "expected different to address")
			s.Require().Equal(tc.token, mintEvent.Token, "expected different token")
			s.Require().Equal(tc.value, mintEvent.Value, "expected different value")
		})
	}
}
