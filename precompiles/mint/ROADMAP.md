## More features

- Burn precompile. `burn(address from, string token, uint256 amount)`, for complete token lifecycle.
- Transaction batching. `mintBatch()` and `burnBatch()` for gas efficiency
- Improve on token metadata mgt. set/get the token name/symbol/decimals via precompile.
- Improve `IsValidRecipient()` check. Currently only checking for zero address and basic validity. Enhancement ideas:
  - can make it check if address contains contract bytecode
  - can maintain a blacklist of known module accounts to prevent accidental minting to system account
  - can validate against known precompile addresses to prevent circular dependencies
  - **_all of these will prevent infinite lock problem._**
- Better event logging. Can help in analytics/alerts/monitoring infra  debugging.

## Recommended by assignment
- IBC middleware integration to auto-mit tokens when receiving IBC transfers. good for cross-chain token bridging, and expanding ecosystem.
- 