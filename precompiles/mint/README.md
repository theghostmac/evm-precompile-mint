## EVM Mint Precompile

## What I built
A custom EVM precompile that allows authorized account to mint native Cosmos tokens.

Features:
- admin-only minting with authorized checks. supports hex and bech32 formats
- input validation
- native Cosmos SDK token integration via BankKeeper
- EVM event logging


## How to Test/Run
Only testing applicable:
```bash
go test ./...
```

`TestMint()` contains the main minting logic with some failure and a success case.
`TestIsAuthorized()` handles authorization logic
`TestIsValidRecipient()` handles address validation (limited)
`TestEmitMintEvent()` checks that events are emitted.

## Challenges

- First, I tried using spawn twice to create Cosmos EVM app, didn't work. I carefully selected the evm option in the toggle menu.
- I picked the `erc20` precompile as a guide without comparing with others on which one comes close to `mint` due to time constraints.
    - Thankfully, it was easy to follow, and hopefully, my implementation is the best way to do this.
- Searching for the test suite all around the codebase was exhausting. Had to reverse-engineer patterns from the erc20 precompile code due to sparse precompile documentation.
    - Either I didn't find the documentation, or it didn't exist.
- The `BankKeeper` interface not including minting methods is somewhat surprising; but I guess that's the purpose of the assignment.
    - I extended the interface.
- Spent significant time debugging why the admin authorization was failing. Discovered its because `s.SetupTest()` creates a new keyring and precompile instance with a different authority address in each test case.
  - Moved `s.SetupTest()` outside the test loop to fix this address mismatch issue.
  - Should have spent time understanding the test suite and lifecycle.

## What went well
- Pattern recognition with the erc20 precompile as a template for adding the mint precompile - both logic and testing pattern.
- ABI integration was a breeze. loading and method dispatch system is well-designed and straightforward.
- supporting both address formats (hex/bech32) was great. I would have spent more time there.

## Design decisions

- using the 0x1111 address for now, for easy test and memorizing.
- hardcoded the authority admin address for now for simplicity.