package keeper_test

import (
	"github.com/strangelove-ventures/tokenfactory/x/tokenfactory/keeper"
	"github.com/strangelove-ventures/tokenfactory/x/tokenfactory/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func (suite *KeeperTestSuite) TestGenesis() {
	genesisState := types.GenesisState{
		FactoryDenoms: []types.GenesisDenom{
			{
				Denom: "factory/cosmos1t7egva48prqmzl59x5ngv4zx0dtrwewcdqdjr8/bitcoin",
				AuthorityMetadata: types.DenomAuthorityMetadata{
					Admin: "cosmos1t7egva48prqmzl59x5ngv4zx0dtrwewcdqdjr8",
				},
			},
			{
				Denom: "factory/cosmos1t7egva48prqmzl59x5ngv4zx0dtrwewcdqdjr8/diff-admin",
				AuthorityMetadata: types.DenomAuthorityMetadata{
					Admin: "cosmos15czt5nhlnvayqq37xun9s9yus0d6y26dx74r5p",
				},
			},
			{
				Denom: "factory/cosmos1t7egva48prqmzl59x5ngv4zx0dtrwewcdqdjr8/litecoin",
				AuthorityMetadata: types.DenomAuthorityMetadata{
					Admin: "cosmos1t7egva48prqmzl59x5ngv4zx0dtrwewcdqdjr8",
				},
			},
		},
		SudoAdmins: []string{
			"cosmos1t7egva48prqmzl59x5ngv4zx0dtrwewcdqdjr8",
			"cosmos15czt5nhlnvayqq37xun9s9yus0d6y26dx74r5p",
		},
	}

	suite.SetupTestForInitGenesis()
	app := suite.App

	// Test both with bank denom metadata set, and not set.
	for i, denom := range genesisState.FactoryDenoms {
		// hacky, sets bank metadata to exist if i != 0, to cover both cases.
		if i != 0 {
			app.BankKeeper.SetDenomMetaData(suite.Ctx, banktypes.Metadata{Base: denom.GetDenom()})
		}
	}

	if err := app.TokenFactoryKeeper.SetParams(suite.Ctx, types.Params{DenomCreationFee: sdk.Coins{sdk.NewInt64Coin("stake", 100)}}); err != nil {
		panic(err)
	}
	app.TokenFactoryKeeper.InitGenesis(suite.Ctx, genesisState)

	exportedGenesis := app.TokenFactoryKeeper.ExportGenesis(suite.Ctx)
	suite.Require().NotNil(exportedGenesis)
	suite.Require().Equal(genesisState, *exportedGenesis)
}

func (suite *KeeperTestSuite) TestGenesisStoresSudoAdmins() {
	genesisState := types.GenesisState{
		SudoAdmins: []string{
			suite.TestAccs[0].String(),
			suite.TestAccs[1].String(),
		},
	}

	sa := keeper.SudoAdmins{Keeper: suite.App.TokenFactoryKeeper}
	suite.App.TokenFactoryKeeper.InitGenesis(suite.Ctx, genesisState)

	suite.Require().True(sa.IsSudoAdmin(suite.Ctx, suite.TestAccs[0].String()))
	suite.Require().True(sa.IsSudoAdmin(suite.Ctx, suite.TestAccs[1].String()))
}
