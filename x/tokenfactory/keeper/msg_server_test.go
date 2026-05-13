package keeper_test

import (
	"fmt"

	"github.com/strangelove-ventures/tokenfactory/app"
	"github.com/strangelove-ventures/tokenfactory/x/tokenfactory/keeper"
	"github.com/strangelove-ventures/tokenfactory/x/tokenfactory/types"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// TestMintDenomMsg tests TypeMsgMint message is emitted on a successful mint
func (suite *KeeperTestSuite) TestMintDenomMsg() {
	// Create a denom
	suite.CreateDefaultDenom()
	ctx := suite.Ctx.WithEventManager(sdk.NewEventManager())
	nonFactoryDenom := "unique"
	udc := keeper.NewUnboundDenomCreator(suite.App.TokenFactoryKeeper)
	udc.CreateDenom(ctx, suite.TestAccs[0].String(), nonFactoryDenom)

	for _, tc := range []struct {
		desc                  string
		amount                int64
		mintDenom             string
		sender                string
		sudoer                string
		expectedMessageEvents int // the valid case should emit >= 1
	}{
		{
			desc:      "denom does not exist",
			amount:    10,
			mintDenom: "factory/osmo1t7egva48prqmzl59x5ngv4zx0dtrwewc9m7z44/evmos",
			sender:    suite.TestAccs[0].String(),
		},
		{
			desc:                  "success case tokenfactory",
			amount:                10,
			mintDenom:             suite.defaultDenom,
			sender:                suite.TestAccs[0].String(),
			expectedMessageEvents: 1,
		},
		{
			desc:                  "invalid mint from non admin for factory-own denom",
			amount:                10,
			mintDenom:             suite.defaultDenom,
			sender:                suite.TestAccs[1].String(),
			expectedMessageEvents: 0,
		},
		// Sudo Mints
		{
			desc:                  "successful mint of *NON* factory-own denom by non admin",
			amount:                10,
			mintDenom:             nonFactoryDenom,
			sender:                suite.TestAccs[0].String(),
			expectedMessageEvents: 1,
		},
		{
			desc:                  "invalid mint of *NON* factory-own denom by non admin",
			amount:                10,
			mintDenom:             nonFactoryDenom,
			sender:                suite.TestAccs[1].String(),
			expectedMessageEvents: 0,
		},
	} {
		suite.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			ctx := suite.Ctx.WithEventManager(sdk.NewEventManager())
			suite.Require().Equal(0, len(ctx.EventManager().Events()))

			suite.OverrideMsgServer(suite.App.TokenFactoryKeeper)

			// Test mint message
			suite.msgServer.Mint(ctx, types.NewMsgMint(tc.sender, sdk.NewInt64Coin(tc.mintDenom, tc.amount))) //nolint:errcheck

			// Ensure current number and type of event is emitted
			suite.AssertEventEmitted(ctx, types.TypeMsgMint, tc.expectedMessageEvents)
		})
	}
}

// TestBurnDenomMsg tests TypeMsgBurn message is emitted on a successful burn
func (suite *KeeperTestSuite) TestBurnDenomMsg() {
	// Create a denom.
	suite.CreateDefaultDenom()
	nonFactoryDenom := "unique"
	unregisteredNonFactoryDenom := "unregistered"
	amount := int64(10)

	unboundCoins := sdk.NewCoins(sdk.NewInt64Coin(nonFactoryDenom, amount))
	suite.Assert().NoError(suite.App.BankKeeper.MintCoins(suite.Ctx, types.ModuleName, unboundCoins))
	suite.Assert().NoError(suite.App.BankKeeper.SendCoinsFromModuleToAccount(suite.Ctx, types.ModuleName, suite.TestAccs[0], unboundCoins))

	unregisteredUnboundCoins := sdk.NewCoins(sdk.NewInt64Coin(unregisteredNonFactoryDenom, amount))
	suite.Assert().NoError(suite.App.BankKeeper.MintCoins(suite.Ctx, types.ModuleName, unregisteredUnboundCoins))
	suite.Assert().NoError(suite.App.BankKeeper.SendCoinsFromModuleToAccount(suite.Ctx, types.ModuleName, suite.TestAccs[0], unregisteredUnboundCoins))

	ctx := suite.Ctx.WithEventManager(sdk.NewEventManager())

	admin2 := suite.TestAccs[1].String()

	udc := keeper.NewUnboundDenomCreator(suite.App.TokenFactoryKeeper)
	suite.Assert().NoError(udc.CreateDenom(ctx, admin2, nonFactoryDenom))

	//factoryDenom, err := suite.App.TokenFactoryKeeper.CreateDenom(ctx, admin2, nonFactoryDenom)
	//suite.Assert().NoError(err)

	////expectedFactoryDenom, err := types.GetTokenDenom(admin2, nonFactoryDenom)
	////suite.Assert().NoError(err)
	////suite.Assert().Equal(factoryDenom, expectedFactoryDenom)

	// Proof, that both denoms have been correctly created and the `admin2` account is their admin:
	res, _ := suite.App.TokenFactoryKeeper.DenomsFromAdmin(ctx, &types.QueryDenomsFromAdminRequest{Admin: admin2})
	denoms := types.NewSet[string](res.GetDenoms()...)
	suite.Assert().True(denoms.Contains(nonFactoryDenom))
	//suite.Assert().True(denoms.Contains(factoryDenom))

	// mint 10 default token for testAcc[2]
	suite.msgServer.Mint(suite.Ctx, types.NewMsgMintTo(suite.TestAccs[0].String(), sdk.NewInt64Coin(suite.defaultDenom, amount), suite.TestAccs[2].String())) //nolint:errcheck
	// mint 10 default token for admin testAcc[0]
	suite.msgServer.Mint(suite.Ctx, types.NewMsgMintTo(suite.TestAccs[0].String(), sdk.NewInt64Coin(suite.defaultDenom, amount), suite.TestAccs[0].String())) //nolint:errcheck

	capabilities_EnableBurnUnregistered_DISABLED := []string{
		types.EnableBurnOwn,
		//types.EnableBurnOwnUnregistered,
		types.EnableBurnFrom,
		types.EnableForceTransfer,
		types.EnableSetMetadata,
		types.EnableSudoMint,
		types.EnableCommunityPoolFeeFunding,
	}

	capabilities_EnableBurnOwn_DISABLED := []string{
		//types.EnableBurnOwn,
		types.EnableBurnOwnUnregistered,
		types.EnableBurnFrom,
		types.EnableForceTransfer,
		types.EnableSetMetadata,
		types.EnableSudoMint,
		types.EnableCommunityPoolFeeFunding,
	}

	capabilities_EnableBurnFrom_DISABLED := []string{
		types.EnableBurnOwn,
		types.EnableBurnOwnUnregistered,
		//types.EnableBurnFrom,
		types.EnableForceTransfer,
		types.EnableSetMetadata,
		types.EnableSudoMint,
		types.EnableCommunityPoolFeeFunding,
	}

	for _, tc := range []struct {
		desc                  string
		amount                int64
		burnDenom             string
		admin                 string
		burnFrom              string
		valid                 bool
		expectedMessageEvents int
		capabilities          []string
	}{
		{
			desc:                  "denom does not exist",
			burnDenom:             "factory/osmo1t7egva48prqmzl59x5ngv4zx0dtrwewc9m7z44/evmos",
			admin:                 suite.TestAccs[0].String(),
			burnFrom:              suite.TestAccs[2].String(),
			valid:                 false,
			amount:                1,
			expectedMessageEvents: 0,
		},
		{
			desc:                  "success case",
			burnDenom:             suite.defaultDenom,
			admin:                 suite.TestAccs[0].String(),
			burnFrom:              suite.TestAccs[2].String(),
			valid:                 true,
			amount:                1,
			expectedMessageEvents: 1,
		},
		{
			desc:                  "EnableBurnOwnUnregistered ENABLED: successful burn of OWN Registerd non-factory denom coins",
			burnDenom:             nonFactoryDenom,
			admin:                 suite.TestAccs[0].String(),
			valid:                 true,
			amount:                1,
			expectedMessageEvents: 1,
		},
		{
			desc:                  "EnableBurnOwnUnregistered DISABLED: successful burn of OWN Registerd non-factory denom coins",
			burnDenom:             nonFactoryDenom,
			admin:                 suite.TestAccs[0].String(),
			valid:                 true,
			amount:                1,
			expectedMessageEvents: 1,
			capabilities:          capabilities_EnableBurnUnregistered_DISABLED,
		},
		{
			desc:                  "EnableBurnOwnUnregistered ENABLED: successful burn of OWN UNregisterd non-factory denom coins",
			burnDenom:             unregisteredNonFactoryDenom,
			admin:                 suite.TestAccs[0].String(),
			valid:                 true,
			amount:                1,
			expectedMessageEvents: 1,
		},
		{
			desc:                  "EnableBurnOwnUnregistered DISABLED: failed burn of OWN UNregisterd non-factory denom coins",
			burnDenom:             unregisteredNonFactoryDenom,
			admin:                 suite.TestAccs[0].String(),
			valid:                 false,
			amount:                1,
			expectedMessageEvents: 0,
			capabilities:          capabilities_EnableBurnUnregistered_DISABLED,
		},
		{
			desc:                  "EnableBurnOwn DISABLED: failed burn of OWN non-factory denom coins",
			burnDenom:             nonFactoryDenom,
			admin:                 suite.TestAccs[0].String(),
			valid:                 false,
			amount:                1,
			expectedMessageEvents: 0,
			capabilities:          capabilities_EnableBurnOwn_DISABLED,
		},
		{
			desc:                  "EnableBurnOwn DISABLED: successful burn of non-factory denom coins as admin",
			burnDenom:             nonFactoryDenom,
			admin:                 admin2,
			burnFrom:              suite.TestAccs[0].String(),
			valid:                 true,
			amount:                1,
			expectedMessageEvents: 1,
			capabilities:          capabilities_EnableBurnOwn_DISABLED,
		},
		{
			desc:                  "EnableBurnFrom DISABLED: successful burn of OWN non-factory denom coins",
			burnDenom:             nonFactoryDenom,
			admin:                 suite.TestAccs[0].String(),
			valid:                 true,
			amount:                1,
			expectedMessageEvents: 1,
			capabilities:          capabilities_EnableBurnFrom_DISABLED,
		},
		{
			desc:                  "EnableBurnFrom DISABLED: failed burn of non-factory denom coins as admin",
			burnDenom:             nonFactoryDenom,
			admin:                 admin2,
			burnFrom:              suite.TestAccs[0].String(),
			valid:                 false,
			amount:                1,
			expectedMessageEvents: 0,
			capabilities:          capabilities_EnableBurnFrom_DISABLED,
		},
		{
			desc:                  "EnableBurnFrom DISABLED: failed burn of factory denom coins as admin",
			burnDenom:             suite.defaultDenom,
			admin:                 suite.TestAccs[0].String(),
			burnFrom:              suite.TestAccs[2].String(),
			valid:                 false,
			amount:                1,
			expectedMessageEvents: 0,
			capabilities:          capabilities_EnableBurnFrom_DISABLED,
		},
		{
			desc:                  "EnableBurnFrom DISABLED: successful burn of admin's own factory denom coins",
			burnDenom:             suite.defaultDenom,
			admin:                 suite.TestAccs[0].String(),
			valid:                 true,
			amount:                1,
			expectedMessageEvents: 1,
			capabilities:          capabilities_EnableBurnFrom_DISABLED,
		},
	} {
		suite.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			if tc.capabilities == nil {
				tc.capabilities = app.TokenFactoryAllCapabilities
			}

			suite.App.TokenFactoryKeeper.SetEnabledCapabilities(suite.Ctx, tc.capabilities)
			suite.msgServer = keeper.NewMsgServerImpl(suite.App.TokenFactoryKeeper)

			ctx := suite.Ctx.WithEventManager(sdk.NewEventManager())
			suite.Require().Equal(0, len(ctx.EventManager().Events()))

			// Test burn message
			var msgBurn *types.MsgBurn
			if tc.burnFrom == "" {
				msgBurn = types.NewMsgBurn(tc.admin, sdk.NewInt64Coin(tc.burnDenom, tc.amount))
			} else {
				msgBurn = types.NewMsgBurnFrom(tc.admin, sdk.NewInt64Coin(tc.burnDenom, tc.amount), tc.burnFrom)
			}
			suite.msgServer.Burn(ctx, msgBurn) //nolint:errcheck

			// Ensure current number and type of event is emitted
			suite.AssertEventEmitted(ctx, types.TypeMsgBurn, tc.expectedMessageEvents)
		})
	}
}

// TestCreateDenomMsg tests TypeMsgCreateDenom message is emitted on a successful denom creation
func (suite *KeeperTestSuite) TestCreateDenomMsg() {
	defaultDenomCreationFee := types.Params{DenomCreationFee: sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(50000000)))}
	for _, tc := range []struct {
		desc                  string
		denomCreationFee      types.Params
		subdenom              string
		valid                 bool
		expectedMessageEvents int
	}{
		{
			desc:             "subdenom too long",
			denomCreationFee: defaultDenomCreationFee,
			subdenom:         "assadsadsadasdasdsadsadsadsadsadsadsklkadaskkkdasdasedskhanhassyeunganassfnlksdflksafjlkasd",
			valid:            false,
		},
		{
			desc:                  "success case: defaultDenomCreationFee",
			denomCreationFee:      defaultDenomCreationFee,
			subdenom:              "evmos",
			valid:                 true,
			expectedMessageEvents: 1,
		},
	} {
		suite.SetupTest()
		suite.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			tokenFactoryKeeper := suite.App.TokenFactoryKeeper
			ctx := suite.Ctx.WithEventManager(sdk.NewEventManager())
			suite.Require().Equal(0, len(ctx.EventManager().Events()))
			// Set denom creation fee in params
			if err := tokenFactoryKeeper.SetParams(suite.Ctx, tc.denomCreationFee); err != nil {
				suite.Require().NoError(err)
			}
			// Test create denom message
			suite.msgServer.CreateDenom(ctx, types.NewMsgCreateDenom(suite.TestAccs[0].String(), tc.subdenom)) //nolint:errcheck
			// Ensure current number and type of event is emitted
			suite.AssertEventEmitted(ctx, types.TypeMsgCreateDenom, tc.expectedMessageEvents)
		})
	}
}

// TestChangeAdminDenomMsg tests TypeMsgChangeAdmin message is emitted on a successful admin change
func (suite *KeeperTestSuite) TestChangeAdminDenomMsg() {
	for _, tc := range []struct {
		desc                    string
		msgChangeAdmin          func(denom string) *types.MsgChangeAdmin
		expectedChangeAdminPass bool
		expectedAdminIndex      int
		msgMint                 func(denom string) *types.MsgMint
		expectedMintPass        bool
		expectedMessageEvents   int
	}{
		{
			desc: "non-admins can't change the existing admin",
			msgChangeAdmin: func(denom string) *types.MsgChangeAdmin {
				return types.NewMsgChangeAdmin(suite.TestAccs[1].String(), denom, suite.TestAccs[2].String())
			},
			expectedChangeAdminPass: false,
			expectedAdminIndex:      0,
		},
		{
			desc: "success change admin",
			msgChangeAdmin: func(denom string) *types.MsgChangeAdmin {
				return types.NewMsgChangeAdmin(suite.TestAccs[0].String(), denom, suite.TestAccs[1].String())
			},
			expectedAdminIndex:      1,
			expectedChangeAdminPass: true,
			expectedMessageEvents:   1,
			msgMint: func(denom string) *types.MsgMint {
				return types.NewMsgMint(suite.TestAccs[1].String(), sdk.NewInt64Coin(denom, 5))
			},
			expectedMintPass: true,
		},
	} {
		suite.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			// setup test
			suite.SetupTest()
			ctx := suite.Ctx.WithEventManager(sdk.NewEventManager())
			suite.Require().Equal(0, len(ctx.EventManager().Events()))
			// Create a denom and mint
			res, err := suite.msgServer.CreateDenom(ctx, types.NewMsgCreateDenom(suite.TestAccs[0].String(), "bitcoin"))
			suite.Require().NoError(err)
			suite.AssertEventEmitted(ctx, types.SetAdminEvent, 1)
			testDenom := res.GetNewTokenDenom()
			suite.msgServer.Mint(ctx, types.NewMsgMint(suite.TestAccs[0].String(), sdk.NewInt64Coin(testDenom, 10))) //nolint:errcheck
			// Test change admin message
			suite.msgServer.ChangeAdmin(ctx, tc.msgChangeAdmin(testDenom)) //nolint:errcheck
			// Ensure current number and type of event is emitted
			suite.AssertEventEmitted(ctx, types.TypeMsgChangeAdmin, tc.expectedMessageEvents)
		})
	}
}

// TestSetDenomMetaDataMsg tests TypeMsgSetDenomMetadata message is emitted on a successful denom metadata change
func (suite *KeeperTestSuite) TestSetDenomMetaDataMsg() {
	// setup test
	suite.SetupTest()
	suite.CreateDefaultDenom()

	ctx := suite.Ctx.WithEventManager(sdk.NewEventManager())

	admin2 := suite.TestAccs[1].String()

	nonFactoryDenom := "unique"
	udc := keeper.NewUnboundDenomCreator(suite.App.TokenFactoryKeeper)
	suite.Assert().NoError(udc.CreateDenom(ctx, admin2, nonFactoryDenom))

	factoryDenom, err := suite.App.TokenFactoryKeeper.CreateDenom(ctx, admin2, nonFactoryDenom)
	suite.Assert().NoError(err)

	//expectedFactoryDenom, err := types.GetTokenDenom(admin2, nonFactoryDenom)
	//suite.Assert().NoError(err)
	//suite.Assert().Equal(factoryDenom, expectedFactoryDenom)

	// Proof, that both denoms have been correctly created and the `admin2` account is their admin:
	res, _ := suite.App.TokenFactoryKeeper.DenomsFromAdmin(ctx, &types.QueryDenomsFromAdminRequest{Admin: admin2})
	denoms := types.NewSet[string](res.GetDenoms()...)
	suite.Assert().True(denoms.Contains(nonFactoryDenom))
	suite.Assert().True(denoms.Contains(factoryDenom))

	crateDenomMetadataMsg := func(admin string, fullDenom string) *types.MsgSetDenomMetadata {
		return types.NewMsgSetDenomMetadata(admin2, banktypes.Metadata{
			Description: "yeehaw",
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    fullDenom,
					Exponent: 0,
				},
				{
					Denom:    "u" + fullDenom,
					Exponent: 6,
				},
			},
			Base:    fullDenom,
			Display: fullDenom,
			Name:    "UNIQUE",
			Symbol:  "UNIQUE",
		})
	}

	for _, tc := range []struct {
		desc                  string
		msgSetDenomMetadata   types.MsgSetDenomMetadata
		expectedPass          bool
		expectedMessageEvents int
	}{
		{
			desc: "successful set denom metadata",
			msgSetDenomMetadata: *types.NewMsgSetDenomMetadata(suite.TestAccs[0].String(), banktypes.Metadata{
				Description: "yeehaw",
				DenomUnits: []*banktypes.DenomUnit{
					{
						Denom:    suite.defaultDenom,
						Exponent: 0,
					},
					{
						Denom:    "uosmo",
						Exponent: 6,
					},
				},
				Base:    suite.defaultDenom,
				Display: "uosmo",
				Name:    "OSMO",
				Symbol:  "OSMO",
			}),
			expectedPass:          true,
			expectedMessageEvents: 1,
		},
		{
			desc: "non existent factory denom name",
			msgSetDenomMetadata: *types.NewMsgSetDenomMetadata(suite.TestAccs[0].String(), banktypes.Metadata{
				Description: "yeehaw",
				DenomUnits: []*banktypes.DenomUnit{
					{
						Denom:    fmt.Sprintf("factory/%s/litecoin", suite.TestAccs[0].String()),
						Exponent: 0,
					},
					{
						Denom:    "uosmo",
						Exponent: 6,
					},
				},
				Base:    fmt.Sprintf("factory/%s/litecoin", suite.TestAccs[0].String()),
				Display: "uosmo",
				Name:    "OSMO",
				Symbol:  "OSMO",
			}),
			expectedPass: false,
		},
		{
			desc:                  "proving opposite scenario for the next test: setting denom metadata for factory-own denom passes",
			msgSetDenomMetadata:   *crateDenomMetadataMsg(admin2, factoryDenom),
			expectedPass:          true,
			expectedMessageEvents: 1,
		},
		{
			desc:                  "setting denom metadata for non-factory denom MUST fail",
			msgSetDenomMetadata:   *crateDenomMetadataMsg(admin2, nonFactoryDenom),
			expectedPass:          false,
			expectedMessageEvents: 0,
		},
	} {
		suite.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			tc := tc
			ctx := suite.Ctx.WithEventManager(sdk.NewEventManager())
			suite.Require().Equal(0, len(ctx.EventManager().Events()))

			// Test set denom metadata message
			suite.msgServer.SetDenomMetadata(ctx, &tc.msgSetDenomMetadata) //nolint:errcheck
			// Ensure current number and type of event is emitted
			suite.AssertEventEmitted(ctx, types.TypeMsgSetDenomMetadata, tc.expectedMessageEvents)
		})
	}
}
