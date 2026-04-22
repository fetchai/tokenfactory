package keeper_test

import (
	"github.com/strangelove-ventures/tokenfactory/x/tokenfactory/keeper"
	"github.com/strangelove-ventures/tokenfactory/x/tokenfactory/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestSudoAdminsStore() {
	sa := keeper.SudoAdmins{Keeper: suite.App.TokenFactoryKeeper}
	sudoAdmin := suite.TestAccs[0].String()

	suite.Require().False(sa.IsSudoAdmin(suite.Ctx, sudoAdmin))
	suite.Require().NoError(sa.AddSudoAdmin(suite.Ctx, sudoAdmin))
	suite.Require().True(sa.IsSudoAdmin(suite.Ctx, sudoAdmin))
	suite.Require().Contains(sa.GetAllSudoAdmins(suite.Ctx), sudoAdmin)

	suite.Require().NoError(sa.RemoveSudoAdmin(suite.Ctx, sudoAdmin))
	suite.Require().False(sa.IsSudoAdmin(suite.Ctx, sudoAdmin))
	suite.Require().NotContains(sa.GetAllSudoAdmins(suite.Ctx), sudoAdmin)
}

func (suite *KeeperTestSuite) TestSudoAdminsStoreRejectsInvalidAddress() {
	sa := keeper.SudoAdmins{Keeper: suite.App.TokenFactoryKeeper}

	suite.Require().Error(sa.AddSudoAdmin(suite.Ctx, "not-a-bech32-address"))
	suite.Require().False(sa.IsSudoAdmin(suite.Ctx, "not-a-bech32-address"))
	suite.Require().Error(sa.RemoveSudoAdmin(suite.Ctx, "not-a-bech32-address"))
}

func (suite *KeeperTestSuite) TestSudoAdminMintAnyDenom() {
	sa := keeper.SudoAdmins{Keeper: suite.App.TokenFactoryKeeper}
	suite.Require().NoError(sa.AddSudoAdmin(suite.Ctx, suite.TestAccs[1].String()))

	before := suite.App.BankKeeper.GetBalance(suite.Ctx, suite.TestAccs[2], NativeDenom).Amount

	_, err := suite.msgServer.Mint(
		suite.Ctx,
		types.NewMsgMintTo(
			suite.TestAccs[1].String(),
			sdk.NewInt64Coin(NativeDenom, 77),
			suite.TestAccs[2].String(),
		),
	)
	suite.Require().NoError(err)

	after := suite.App.BankKeeper.GetBalance(suite.Ctx, suite.TestAccs[2], NativeDenom).Amount
	suite.Require().Equal(before.AddRaw(77), after)
}

func (suite *KeeperTestSuite) TestSudoAdminMintFailsWhenCapabilityDisabled() {
	sa := keeper.SudoAdmins{Keeper: suite.App.TokenFactoryKeeper}
	suite.Require().NoError(sa.AddSudoAdmin(suite.Ctx, suite.TestAccs[1].String()))

	suite.App.TokenFactoryKeeper.SetEnabledCapabilities(suite.Ctx, []string{})
	suite.OverrideMsgServer(suite.App.TokenFactoryKeeper)

	before := suite.App.BankKeeper.GetBalance(suite.Ctx, suite.TestAccs[2], NativeDenom).Amount

	_, err := suite.msgServer.Mint(
		suite.Ctx,
		types.NewMsgMintTo(
			suite.TestAccs[1].String(),
			sdk.NewInt64Coin(NativeDenom, 77),
			suite.TestAccs[2].String(),
		),
	)
	suite.Require().Error(err)

	after := suite.App.BankKeeper.GetBalance(suite.Ctx, suite.TestAccs[2], NativeDenom).Amount
	suite.Require().Equal(before, after)
}

func (suite *KeeperTestSuite) TestNonSudoCannotMintAnyDenom() {
	before := suite.App.BankKeeper.GetBalance(suite.Ctx, suite.TestAccs[2], NativeDenom).Amount

	_, err := suite.msgServer.Mint(
		suite.Ctx,
		types.NewMsgMintTo(
			suite.TestAccs[1].String(),
			sdk.NewInt64Coin(NativeDenom, 77),
			suite.TestAccs[2].String(),
		),
	)
	suite.Require().Error(err)

	after := suite.App.BankKeeper.GetBalance(suite.Ctx, suite.TestAccs[2], NativeDenom).Amount
	suite.Require().Equal(before, after)
}
