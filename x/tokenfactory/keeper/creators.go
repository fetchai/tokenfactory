package keeper

import (
	"context"

	"cosmossdk.io/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/strangelove-ventures/tokenfactory/x/tokenfactory/types"
)

func (k Keeper) addDenomFromCreator(ctx sdk.Context, creator, denom string) {
	store := k.GetCreatorPrefixStore(sdk.UnwrapSDKContext(ctx), creator)
	store.Set([]byte(denom), []byte(denom))

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeMsgCreateDenom,
			sdk.NewAttribute(types.AttributeCreator, creator),
			sdk.NewAttribute(types.AttributeNewTokenDenom, denom),
		),
	})
}

func (k Keeper) GetDenomsFromCreator(ctx context.Context, creator string) []string {
	store := k.GetCreatorPrefixStore(sdk.UnwrapSDKContext(ctx), creator)

	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	var denoms []string
	for ; iterator.Valid(); iterator.Next() {
		denoms = append(denoms, string(iterator.Key()))
	}
	return denoms
}

func (k Keeper) GetAllDenomsIterator(ctx context.Context) store.Iterator {
	return k.GetCreatorsPrefixStore(sdk.UnwrapSDKContext(ctx)).Iterator(nil, nil)
}
