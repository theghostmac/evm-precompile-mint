package mint

import (
	"math/big"

	"github.com/cosmos/evm/precompiles/mint"
)

func (s *PrecompileTestSuite) TestIsTransaction() {
	s.SetupTest()

	// The mint method should be seen as a transaction
	method := s.precompile.Methods[mint.MintMethod]
	s.Require().True(s.precompile.IsTransaction(&method), "mint should be identified as a transaction")
}

func (s *PrecompileTestSuite) TestRequiredGas() {
	s.SetupTest()

	testcases := []struct {
		name     string
		malleate func() []byte
		expGas   uint64
	}{
		{
			name: mint.MintMethod,
			malleate: func() []byte {
				to := s.keyring.GetAddr(0)
				token := "umint"
				value := big.NewInt(1000)

				bz, err := s.precompile.Pack(mint.MintMethod, to, token, value)
				s.Require().NoError(err, "expected no error packing ABI")
				return bz
			},
			expGas: mint.GasMint,
		},
		{
			name: "invalid method",
			malleate: func() []byte {
				return []byte("invalid method")
			},
			expGas: 0,
		},
		{
			name: "input bytes too short",
			malleate: func() []byte {
				return []byte{0x00, 0x00, 0x00}
			},
			expGas: 0,
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			input := tc.malleate()
			s.Require().Equal(tc.expGas, s.precompile.RequiredGas(input), "expected different gas requirement")
		})
	}
}

func (s *PrecompileTestSuite) TestAddress() {
	s.SetupTest()

	expectedAddr := "0x0000000000000000000000000000000000001111"
	s.Require().Equal(expectedAddr, s.precompile.Address().Hex(), "expected different precompile address")
}

func (s *PrecompileTestSuite) TestNewPrecompile() {
	testcases := []struct {
		name        string
		authority   string
		expectError bool
		errContains string
	}{
		{
			name:        "pass - valid hex authority",
			authority:   "0x742d35Cc6635C0532925a3b8D435C617d1C3d4E3",
			expectError: false,
		},
		{
			name:        "pass - valid bech32 authority",
			authority:   "cosmos1wskv7h9wrxdkzr8npq8g4juh4etm9khqxvcje4",
			expectError: false,
		},
		{
			name:        "pass - empty authority (for testing)",
			authority:   "",
			expectError: false,
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			s.SetupTest()

			precompile, err := mint.NewPrecompile(
				tc.authority,
				s.network.App.GetBankKeeper(),
			)

			if tc.expectError {
				s.Require().Error(err, "expected error creating precompile")
				s.Require().Contains(err.Error(), tc.errContains, "expected different error message")
				s.Require().Nil(precompile, "expected nil precompile")
			} else {
				s.Require().NoError(err, "expected no error creating precompile")
				s.Require().NotNil(precompile, "expected non-nil precompile")
				s.Require().Equal(tc.authority, precompile.GetAuthority(), "expected different authority")
			}
		})
	}
}
