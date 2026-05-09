package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// UnboundDenomCreator provides the only supported way to create tokenfactory-unbound denominations externally.
type UnboundDenomCreator interface {
	CreateDenom(ctx sdk.Context, creatorAddr string, denom string) error
}

type unboundDenomCreator struct {
	Keeper
}

func NewUnboundDenomCreator(keeper Keeper) UnboundDenomCreator {
	return &unboundDenomCreator{Keeper: keeper}
}

// CreateDenom registers authority metadata for tokenfactory-unbound denomination under a creator.
//
// This includes denominations unrelated to tokenfactory, such as the chain’s
// native staking denom.
//
// Unlike tokenfactory-specific creation flow, this method:
//   - does *not* set Bank.Metadata
//   - does *not* perform tokenfactory-specific denom validation
// However, CosmosSDK canonical denom validation *is* performed.
func (k *unboundDenomCreator) CreateDenom(ctx sdk.Context, creatorAddr string, denom string) error {
	if err := sdk.ValidateDenom(denom); err != nil {
		return err
	}

	return k.createDenomAfterValidation(ctx, creatorAddr, denom)
}
