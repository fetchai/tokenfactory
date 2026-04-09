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
	return prefix.NewStore(store, types.GetSudoAdminsPrefix())
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

func (k SudoAdmins) GetAllSudoAdmins(ctx context.Context) []string {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := k.GetSudoAdminsStore(sdkCtx)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	prefix := sdk.GetConfig().GetBech32AccountAddrPrefix()
	admins := []string{}
	for ; iterator.Valid(); iterator.Next() {
		admin, err := sdk.Bech32ifyAddressBytes(prefix, iterator.Key())
		if err != nil {
			panic(err)
		}
		admins = append(admins, admin)
	}

	return admins
}
