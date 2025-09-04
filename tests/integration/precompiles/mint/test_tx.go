package mint

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/cosmos/evm/precompiles/mint"
	"github.com/cosmos/evm/precompiles/testutil"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	toAddr = common.HexToAddress("0x742d35Cc6635C0532925a3b8D435C617d1C3d4E3")
)

func (s *PrecompileTestSuite) TestMint() {
	// Once setup for all test cases!
	s.SetupTest()

	method := s.precompile.Methods[mint.MintMethod]
	admin := s.keyring.GetKey(0)    // first key is the authority/admin
	nonAdmin := s.keyring.GetKey(1) // second key is a non-admin user

	testcases := []struct {
		name        string
		caller      common.Address
		malleate    func() []interface{}
		postCheck   func()
		expErr      bool
		errContains string
	}{
		{
			name:   "fail - unauthorized caller",
			caller: nonAdmin.Addr,
			malleate: func() []interface{} {
				return []interface{}{toAddr, "umint", big.NewInt(1000)}
			},
			expErr:      true,
			errContains: mint.ErrUnauthorized.Error(),
		},
		{
			name:   "fail - invalid number of arguments",
			caller: admin.Addr,
			malleate: func() []interface{} {
				return []interface{}{toAddr, "umint"} // Missing amount
			},
			expErr:      true,
			errContains: "invalid number of arguments",
		},
		{
			name:   "fail - invalid to address",
			caller: admin.Addr,
			malleate: func() []interface{} {
				return []interface{}{"invalid", "umint", big.NewInt(1000)}
			},
			expErr:      true,
			errContains: "invalid to address",
		},
		{
			name:   "fail - invalid token string",
			caller: admin.Addr,
			malleate: func() []interface{} {
				return []interface{}{toAddr, 123, big.NewInt(1000)} // Invalid token type
			},
			expErr:      true,
			errContains: "invalid token",
		},
		{
			name:   "fail - invalid amount",
			caller: admin.Addr,
			malleate: func() []interface{} {
				return []interface{}{toAddr, "umint", "invalid"} // Invalid amount type
			},
			expErr:      true,
			errContains: "invalid value",
		},
		{
			name:   "fail - zero amount",
			caller: admin.Addr,
			malleate: func() []interface{} {
				return []interface{}{toAddr, "umint", big.NewInt(0)}
			},
			expErr:      true,
			errContains: mint.ErrZeroAmount.Error(),
		},
		{
			name:   "fail - negative amount",
			caller: admin.Addr,
			malleate: func() []interface{} {
				return []interface{}{toAddr, "umint", big.NewInt(-100)}
			},
			expErr:      true,
			errContains: mint.ErrNegativeAmount.Error(),
		},
		{
			name:   "fail - invalid denomination",
			caller: admin.Addr,
			malleate: func() []interface{} {
				return []interface{}{toAddr, "", big.NewInt(1000)} // Empty denom
			},
			expErr:      true,
			errContains: "invalid token denomination",
		},
		{
			name:   "pass - successful mint",
			caller: admin.Addr,
			malleate: func() []interface{} {
				return []interface{}{toAddr, "umint", big.NewInt(1000)}
			},
			postCheck: func() {
				// Check that tokens were minted
				balance := s.network.App.GetBankKeeper().GetBalance(
					s.network.GetContext(),
					toAddr.Bytes(),
					"umint",
				)
				s.Require().Equal(big.NewInt(1000), balance.Amount.BigInt(), "expected tokens to be minted")
			},
			expErr: false,
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			stateDB := s.network.GetStateDB()

			var contract *vm.Contract
			contract, ctx := testutil.NewPrecompileContract(
				s.T(),
				s.network.GetContext(),
				tc.caller,
				s.precompile.Address(),
				0,
			)

			_, err := s.precompile.Mint(ctx, contract, stateDB, &method, tc.malleate())
			if tc.expErr {
				s.Require().Error(err, "expected mint transaction to fail")
				s.Require().Contains(err.Error(), tc.errContains, "expected mint transaction to fail with specific error")
			} else {
				s.Require().NoError(err, "expected mint transaction to succeed")
				if tc.postCheck != nil {
					tc.postCheck()
				}
			}
		})
	}
}

func (s *PrecompileTestSuite) TestIsAuthorized() {
	s.SetupTest()
	ctx := s.network.GetContext()

	admin := s.keyring.GetKey(0)
	nonAdmin := s.keyring.GetKey(1)

	testcases := []struct {
		name      string
		caller    common.Address
		authority string
		expected  bool
	}{
		{
			name:      "pass - admin hex address",
			caller:    admin.Addr,
			authority: admin.Addr.Hex(),
			expected:  true,
		},
		{
			name:      "pass - admin bech32 address",
			caller:    admin.Addr,
			authority: sdk.AccAddress(admin.Addr.Bytes()).String(),
			expected:  true,
		},
		{
			name:      "fail - non-admin",
			caller:    nonAdmin.Addr,
			authority: admin.Addr.Hex(),
			expected:  false,
		},
		{
			name:      "fail - zero address",
			caller:    common.Address{},
			authority: admin.Addr.Hex(),
			expected:  false,
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			// Create precompile with specific authority
			precompile, err := mint.NewPrecompile(
				tc.authority,
				s.network.App.GetBankKeeper(),
			)
			s.Require().NoError(err)

			result := precompile.IsAuthorized(ctx, tc.caller)
			s.Require().Equal(tc.expected, result, "unexpected authorization result")
		})
	}
}

func (s *PrecompileTestSuite) TestIsValidRecipient() {
	s.SetupTest()
	ctx := s.network.GetContext()

	testcases := []struct {
		name     string
		to       common.Address
		expected bool
	}{
		{
			name:     "pass - valid address",
			to:       toAddr,
			expected: true,
		},
		{
			name:     "pass - random address",
			to:       common.HexToAddress("0x1234567890123456789012345678901234567890"),
			expected: true,
		},
		{
			name:     "fail - zero address",
			to:       common.Address{},
			expected: false,
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			result := s.precompile.IsValidRecipient(ctx, tc.to)
			s.Require().Equal(tc.expected, result, "unexpected recipient validation result")
		})
	}
}
