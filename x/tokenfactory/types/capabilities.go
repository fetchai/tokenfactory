package types

const (
	EnableSetMetadata   = "enable_metadata"
	EnableForceTransfer = "enable_force_transfer"
	// EnableBurnOwn Allows *ANY* owner of tokens (with *ANY* denomination) to burn these self-owned tokens.
	// If disabled (not present), only admin of a denom can execute the burn *IF* EnableBurnFrom is enabled.
	EnableBurnOwn = "enable_burn_own"
	// EnableBurnOwnUnregistered If enabled, token owner can burn its own tokens of any denomination which is *NOT*
	// registered in tokenfactory.
	// This EnableBurnOwnUnregistered has effect only if the EnableBurnOwn is enabled as well.
	EnableBurnOwnUnregistered = "enable_burn_unregistered"
	EnableBurnFrom            = "enable_burn_from"
	// EnableSudoMint is a High level enabler for minting and burning unbound denominations(= any denominations which
	// do *NOT* conform to the tokenfactory denomination format `factory/<CREATOR_ADDRESS>/<SUB_DENOM>`).
	// In order to enable minting of unbound denominations, such denominations *MUST* be registered in tokenfactory and
	// EnableSudoMint must be enabled.
	// By design, it *NOT* possible to register unbound denomination via tokenfactory MsgCreateDenom message, instead
	// it *MUST* be registered either in `genesis.json` file, or during chain software upgrade using the
	// `UnboundDenomCreator` interface.
	// In order to enable burning of unbound denominations, the EnableSudoMint must be enabled. It is possible to burn
	// unbound denominations which are not registered in token factory, however both - the EnableBurnOwn and
	// EnableBurnOwnUnregistered must be enabled.
	EnableSudoMint = "enable_admin_sudo_mint"

	// EnableCommunityPoolFeeFunding sends tokens to the community pool when a new fee is charged (if one is set in params).
	// This is useful for ICS chains, or networks who wish to just have the fee tokens burned (not gas fees, just the extra on top).
	EnableCommunityPoolFeeFunding = "enable_community_pool_fee_funding"
)
