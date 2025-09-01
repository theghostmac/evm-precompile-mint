// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @dev Interface for the Mint precompile contract.
 * This precompile allows authorized accounts to mint native Cosmos tokens.
 */
interface IMint {
    /**
     * @dev Emitted when tokens are minted to an account.
     * @param to The address that received the minted tokens
     * @param token The denomination of the token that was minted
     * @param value The amount of tokens that were minted
     */
    event Mint(address indexed to, string token, uint256 value);

    /**
     * @dev Mint native tokens to the specified address.
     * Can only be called by the authorized admin account.
     *
     * @param to The address to receive the minted tokens
     * @param token The token denomination to mint (e.g., "umint", "uatom")
     * @param value The amount of tokens to mint
     *
     * Requirements:
     * - Caller must be the authorized admin
     * - `to` cannot be the zero address
     * - `to` cannot be a smart contract or module account
     * - `token` must be a valid denomination
     * - `value` must be greater than zero
     *
     * Emits a {Mint} event.
     */
    function mint(address to, string calldata token, uint256 value) external;
}