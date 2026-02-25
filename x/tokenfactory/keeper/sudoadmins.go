package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) AddSudoAdmin(ctx context.Context, admin string) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := k.GetSudoAdminsStore(sdkCtx)

	addr, err := sdk.AccAddressFromBech32(admin)
	if err != nil {
		return err
	}

	// key = canonical address bytes
	// value = dummy marker
	store.Set(addr.Bytes(), []byte{1})

	return nil
}

func (k Keeper) RemoveSudoAdmin(ctx context.Context, admin string) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := k.GetSudoAdminsStore(sdkCtx)

	addr, err := sdk.AccAddressFromBech32(admin)
	if err != nil {
		return err
	}

	store.Delete(addr.Bytes())
	return nil
}

func (k Keeper) IsSudoAdmin(ctx context.Context, admin string) bool {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := k.GetSudoAdminsStore(sdkCtx)

	addr, err := sdk.AccAddressFromBech32(admin)
	if err != nil {
		return false
	}

	return store.Has(addr.Bytes())
}
