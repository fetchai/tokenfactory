package types

const (
	EnableSetMetadata   = "enable_metadata"
	EnableForceTransfer = "enable_force_transfer"
	// Allows to *ANY* owner of tokens (with *ANY* denomination) to burn these self-owned tokens.
	// If disabled (not present), only admin of a denom or sudoer can execute the burn.
	EnableBurnOwn  = "enable_burn_own"
	EnableBurnFrom = "enable_burn_from"
	// Allows addresses registered as sudo admins imn genesis store to mint tokens of *ANY* denominations.
	// NOTE: with SudoMint enabled, the sudo admin can mint `any` token, not just tokenfactory tokens.
	// This is intended behavior as requested by other teams, rather than having its own module with very minor logic.
	// If you do not wish for this behavior, then either do NOT enable this capability, or implement your own logic.
	EnableSudoMint = "enable_admin_sudo_mint"
	// EnableCommunityPoolFeeFunding sends tokens to the community pool when a new fee is charged (if one is set in params).
	// This is useful for ICS chains, or networks who wish to just have the fee tokens burned (not gas fees, just the extra on top).
	EnableCommunityPoolFeeFunding = "enable_community_pool_fee_funding"
)
