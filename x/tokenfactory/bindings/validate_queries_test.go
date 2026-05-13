package bindings_test

import (
	"fmt"
	"testing"

	wasmbinding "github.com/strangelove-ventures/tokenfactory/x/tokenfactory/bindings"
	"github.com/strangelove-ventures/tokenfactory/x/tokenfactory/keeper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestFullDenom(t *testing.T) {
	actor := RandomAccountAddress()

	specs := map[string]struct {
		addr         string
		subdenom     string
		expFullDenom string
		expErr       bool
	}{
		"valid address": {
			addr:         actor.String(),
			subdenom:     "subDenom1",
			expFullDenom: fmt.Sprintf("factory/%s/subDenom1", actor.String()),
		},
		"empty address": {
			addr:     "",
			subdenom: "subDenom1",
			expErr:   true,
		},
		"invalid address": {
			addr:     "invalid",
			subdenom: "subDenom1",
			expErr:   true,
		},
		"empty sub-denom": {
			addr:         actor.String(),
			subdenom:     "",
			expFullDenom: fmt.Sprintf("factory/%s/", actor.String()),
		},
		"valid sub-denom (contains underscore)": {
			addr:         actor.String(),
			subdenom:     "sub_denom",
			expFullDenom: fmt.Sprintf("factory/%s/sub_denom", actor.String()),
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			// when
			gotFullDenom, gotErr := wasmbinding.GetFullDenom(spec.addr, spec.subdenom)
			// then
			if spec.expErr {
				require.Error(t, gotErr)
				return
			}
			require.NoError(t, gotErr)
			assert.Equal(t, spec.expFullDenom, gotFullDenom, "exp %s but got %s", spec.expFullDenom, gotFullDenom)
		})
	}
}

func TestDenomAdmin(t *testing.T) {
	addr := RandomAccountAddress()
	app, ctx := SetupCustomApp(t, addr)

	// set token creation fee to zero to make testing easier
	tfParams := app.TokenFactoryKeeper.GetParams(ctx)
	tfParams.DenomCreationFee = sdk.NewCoins()
	if err := app.TokenFactoryKeeper.SetParams(ctx, tfParams); err != nil {
		t.Fatal(err)
	}

	// create a subdenom via the token factory
	admin := sdk.AccAddress([]byte("addr1_______________"))
	validSubDenom := "validdenom"
	tfDenom, err := app.TokenFactoryKeeper.CreateDenom(ctx, admin.String(), validSubDenom)
	require.NoError(t, err)
	require.NotEmpty(t, tfDenom)

	registeredUnboundDenom := "registered"
	udc := keeper.NewUnboundDenomCreator(app.TokenFactoryKeeper)
	assert.NoError(t, udc.CreateDenom(ctx, admin.String(), registeredUnboundDenom))

	queryPlugin := wasmbinding.NewQueryPlugin(app.BankKeeper, &app.TokenFactoryKeeper)

	testCases := []struct {
		name        string
		denom       string
		expectErr   bool
		expectAdmin string
	}{
		{
			name:        "valid token factory denom",
			denom:       tfDenom,
			expectAdmin: admin.String(),
		},
		{
			name:        "unregistered valid token factory denom",
			denom:       fmt.Sprintf("factory/%s/%s", RandomBech32AccountAddress(), validSubDenom),
			expectErr:   true,
			expectAdmin: "",
		},
		{
			name:        "unregistered unbound denom",
			denom:       "uosmo",
			expectErr:   true,
			expectAdmin: "",
		},
		{
			name:        "registered unbound denom",
			denom:       registeredUnboundDenom,
			expectErr:   false,
			expectAdmin: admin.String(),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			resp, err := queryPlugin.GetDenomAdmin(ctx, tc.denom)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, tc.expectAdmin, resp.Admin)
			}
		})
	}
}
