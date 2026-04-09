package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	store "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/strangelove-ventures/tokenfactory/x/tokenfactory/types"
)

type SudoAdmins struct {
	Keeper
}

// GetSudoAdminsStore returns the substore for sudoers
func (k SudoAdmins) GetSudoAdminsStore(ctx sdk.Context) store.KVStore {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.GetSudoAdmins())
}

func (k SudoAdmins) AddSudoAdmin(ctx context.Context, admin string) error {
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

func (k SudoAdmins) RemoveSudoAdmin(ctx context.Context, admin string) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := k.GetSudoAdminsStore(sdkCtx)

	addr, err := sdk.AccAddressFromBech32(admin)
	if err != nil {
		return err
	}

	store.Delete(addr.Bytes())
	return nil
}

func (k SudoAdmins) IsSudoAdmin(ctx context.Context, admin string) bool {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := k.GetSudoAdminsStore(sdkCtx)

	addr, err := sdk.AccAddressFromBech32(admin)
	if err != nil {
		return false
	}

	return store.Has(addr.Bytes())
}
