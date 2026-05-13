package keeper

import (
	"context"

	"github.com/strangelove-ventures/tokenfactory/x/tokenfactory/types"

	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (server msgServer) CreateDenom(goCtx context.Context, msg *types.MsgCreateDenom) (*types.MsgCreateDenomResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	denom, err := server.Keeper.CreateDenom(ctx, msg.Sender, msg.Subdenom)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateDenomResponse{
		NewTokenDenom: denom,
	}, nil
}

func (server msgServer) Mint(goCtx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	var err error
	ctx := sdk.UnwrapSDKContext(goCtx)

	// verify that denom is an x/tokenfactory denom, and if it is not, then sudo mint capability must be enabled
	if _, _, err = types.DeconstructDenom(msg.Amount.GetDenom()); err == nil {
		// Denomination *MUST* already exist:
		_, denomExists := server.bankKeeper.GetDenomMetaData(ctx, msg.Amount.Denom)
		if !denomExists {
			return nil, types.ErrDenomDoesNotExist.Wrapf("denom: %s", msg.Amount.Denom)
		}
	} else {
		sudoEnabled := server.IsCapabilityEnabled(types.EnableSudoMint)
		if !sudoEnabled {
			return nil, types.ErrCapabilityNotEnabled.Wrapf("the '%s' capability is NOT enabled", types.EnableSudoMint)
		}
	}

	authorityMetadata, err := server.Keeper.GetAuthorityMetadata(ctx, msg.Amount.GetDenom())
	if err != nil {
		return nil, err
	}

	if msg.Sender != authorityMetadata.GetAdmin() {
		return nil, types.ErrUnauthorized
	}

	if msg.MintToAddress == "" {
		msg.MintToAddress = msg.Sender
	}

	err = server.Keeper.mintTo(ctx, msg.Amount, msg.MintToAddress)
	if err != nil {
		return nil, err
	}

	return &types.MsgMintResponse{}, nil
}

func (server msgServer) Burn(goCtx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isBurningOwn := false
	if msg.BurnFromAddress == "" {
		msg.BurnFromAddress = msg.Sender
		isBurningOwn = true
	} else if msg.BurnFromAddress == msg.Sender {
		isBurningOwn = true
	}

	authorityMetadata, err := server.Keeper.GetAuthorityMetadata(ctx, msg.Amount.GetDenom())
	isRegistered := err == nil
	if !isRegistered && !server.IsCapabilityEnabled(types.EnableBurnOwnUnregistered) {
		return nil, err
	}

	// The following code section is exclusively for case when:
	//  * either burning someone's else's tokens
	//  * or burning own tokens, but EnableBurnOwn is disabled
	if !(isBurningOwn && server.IsCapabilityEnabled(types.EnableBurnOwn)) {
		// Denom admin *can* burn its own tokens even if the EnableBurnFrom is *disabled*.
		// This is sensical, as the admin burns it sown tokens and *not* tokens from another account.
		if !isBurningOwn && !server.IsCapabilityEnabled(types.EnableBurnFrom) {
			return nil, types.ErrCapabilityNotEnabled.Wrapf("the '%s' capability is NOT enabled", types.EnableBurnFrom)
		}

		// verify that denom is an x/tokenfactory denom, and if it is not, then sudo mint capability must be enabled
		if _, _, err := types.DeconstructDenom(msg.Amount.GetDenom()); err == nil {
			// Denomination *MUST* already exist:
			_, denomExists := server.bankKeeper.GetDenomMetaData(ctx, msg.Amount.Denom)
			if !denomExists {
				return nil, types.ErrDenomDoesNotExist.Wrapf("denom: %s", msg.Amount.Denom)
			}
		} else {
			sudoEnabled := server.IsCapabilityEnabled(types.EnableSudoMint)
			if !sudoEnabled {
				return nil, types.ErrCapabilityNotEnabled.Wrapf("the '%s' capability is NOT enabled", types.EnableSudoMint)
			}
		}

		if !isRegistered || msg.Sender != authorityMetadata.GetAdmin() {
			return nil, types.ErrUnauthorized.Wrapf("the '%s' sender is NOT '%s' admin of the '%s' denomination", msg.Sender, authorityMetadata.GetAdmin(), msg.Amount.GetDenom())
		}
	}

	err = server.Keeper.burnFrom(ctx, msg.Amount, msg.BurnFromAddress)
	if err != nil {
		return nil, err
	}

	return &types.MsgBurnResponse{}, nil
}

func (server msgServer) ForceTransfer(goCtx context.Context, msg *types.MsgForceTransfer) (*types.MsgForceTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !server.IsCapabilityEnabled(types.EnableForceTransfer) {
		return nil, types.ErrCapabilityNotEnabled
	}

	authorityMetadata, err := server.Keeper.GetAuthorityMetadata(ctx, msg.Amount.GetDenom())
	if err != nil {
		return nil, err
	}

	if msg.Sender != authorityMetadata.GetAdmin() {
		return nil, types.ErrUnauthorized
	}

	err = server.Keeper.forceTransfer(ctx, msg.Amount, msg.TransferFromAddress, msg.TransferToAddress)
	if err != nil {
		return nil, err
	}

	return &types.MsgForceTransferResponse{}, nil
}

func (server msgServer) ChangeAdmin(goCtx context.Context, msg *types.MsgChangeAdmin) (*types.MsgChangeAdminResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// verify that denom is an x/tokenfactory denom
	if _, _, err := types.DeconstructDenom(msg.Denom); err != nil {
		return nil, err
	}

	authorityMetadata, err := server.Keeper.GetAuthorityMetadata(ctx, msg.Denom)
	if err != nil {
		return nil, err
	}

	if msg.Sender != authorityMetadata.GetAdmin() {
		return nil, types.ErrUnauthorized
	}

	err = server.Keeper.setAdmin(ctx, authorityMetadata, msg.Denom, msg.NewAdmin)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeMsgChangeAdmin,
			sdk.NewAttribute(types.AttributeDenom, msg.GetDenom()),
			sdk.NewAttribute(types.AttributeNewAdmin, msg.NewAdmin),
		),
	})

	return &types.MsgChangeAdminResponse{}, nil
}

func (server msgServer) SetDenomMetadata(goCtx context.Context, msg *types.MsgSetDenomMetadata) (*types.MsgSetDenomMetadataResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !server.IsCapabilityEnabled(types.EnableSetMetadata) {
		return nil, types.ErrCapabilityNotEnabled
	}

	// Defense in depth validation of metadata
	err := msg.Metadata.Validate()
	if err != nil {
		return nil, err
	}

	// verify that denom is an x/tokenfactory denom
	if _, _, err := types.DeconstructDenom(msg.Metadata.Base); err != nil {
		return nil, err
	}

	authorityMetadata, err := server.Keeper.GetAuthorityMetadata(ctx, msg.Metadata.Base)
	if err != nil {
		return nil, err
	}

	if msg.Sender != authorityMetadata.GetAdmin() {
		return nil, types.ErrUnauthorized
	}

	server.Keeper.bankKeeper.SetDenomMetaData(ctx, msg.Metadata)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeMsgSetDenomMetadata,
			sdk.NewAttribute(types.AttributeDenom, msg.Metadata.Base),
			sdk.NewAttribute(types.AttributeDenomMetadata, msg.Metadata.String()),
		),
	})

	return &types.MsgSetDenomMetadataResponse{}, nil
}

func (server msgServer) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if server.authority != req.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", server.authority, req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := server.SetParams(ctx, req.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
