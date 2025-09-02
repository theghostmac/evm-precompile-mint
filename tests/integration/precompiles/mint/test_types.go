package mint

import (
	"math/big"

	"github.com/cosmos/evm/precompiles/mint"
	utiltx "github.com/cosmos/evm/testutil/tx"
)

func (s *PrecompileTestSuite) TestParseMintArgs() {
	to := utiltx.GenerateAddress()
	token := "umint"
	value := big.NewInt(1000)

	testcases := []struct {
		name        string
		args        []interface{}
		expPass     bool
		errContains string
	}{
		{
			name: "pass - correct arguments",
			args: []interface{}{
				to,
				token,
				value,
			},
			expPass: true,
		},
		{
			name: "fail - invalid number of arguments (too few)",
			args: []interface{}{
				to,
				token,
			},
			errContains: "invalid number of arguments; expected 3; got: 2",
		},
		{
			name: "fail - invalid number of arguments (too many)",
			args: []interface{}{
				to,
				token,
				value,
				"extra",
			},
			errContains: "invalid number of arguments; expected 3; got: 4",
		},
		{
			name: "fail - invalid to address",
			args: []interface{}{
				"invalid address",
				token,
				value,
			},
			errContains: "invalid to address",
		},
		{
			name: "fail - invalid token",
			args: []interface{}{
				to,
				123, // should be string
				value,
			},
			errContains: "invalid token",
		},
		{
			name: "fail - invalid value",
			args: []interface{}{
				to,
				token,
				"invalid amount", // should be *big.Int
			},
			errContains: "invalid value",
		},
		{
			name: "fail - nil address",
			args: []interface{}{
				nil,
				token,
				value,
			},
			errContains: "invalid to address",
		},
		{
			name: "fail - nil token",
			args: []interface{}{
				to,
				nil,
				value,
			},
			errContains: "invalid token",
		},
		{
			name: "fail - nil value",
			args: []interface{}{
				to,
				token,
				nil,
			},
			errContains: "invalid value",
		},
	}

	for _, tc := range testcases {
		s.Run(tc.name, func() {
			parsedTo, parsedToken, parsedValue, err := mint.ParseMintArgs(tc.args)
			if tc.expPass {
				s.Require().NoError(err, "unexpected error while parsing mint arguments")
				s.Require().Equal(to, parsedTo, "expected different `to` address")
				s.Require().Equal(token, parsedToken, "expected different `token` address")
				s.Require().Equal(value, parsedValue, "expected different `value` address")
			} else {
				s.Require().Error(err, "expected an error parsing the mint arguments")
				s.Require().ErrorContains(err, tc.errContains, "expected different error message")
			}
		})
	}
}
