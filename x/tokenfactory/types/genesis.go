package types

import (
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// this line is used by starport scaffolding # genesis/types/import

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:        DefaultParams(),
		FactoryDenoms: []GenesisDenom{},
		SudoAdmins:    []string{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	err := gs.Params.Validate()
	if err != nil {
		return err
	}

	seenDenoms := NewSet[string]()
	seenSudoAdmins := NewSet[string]()

	for _, denom := range gs.GetFactoryDenoms() {
		if seenDenoms.Contains(denom.GetDenom()) {
			return errorsmod.Wrapf(ErrInvalidGenesis, "duplicate denom: %s", denom.GetDenom())
		}
		seenDenoms.Add(denom.GetDenom())

		_, _, err := DeconstructDenom(denom.GetDenom())
		if err != nil {
			return err
		}

		if denom.AuthorityMetadata.Admin != "" {
			_, err = sdk.AccAddressFromBech32(denom.AuthorityMetadata.Admin)
			if err != nil {
				return errorsmod.Wrapf(ErrInvalidAuthorityMetadata, "Invalid admin address (%s)", err)
			}
		}
	}

	for _, sudoAdmin := range gs.GetSudoAdmins() {
		if seenSudoAdmins.Contains(sudoAdmin) {
			return errorsmod.Wrapf(ErrInvalidGenesis, "duplicate sudo admin: %s", sudoAdmin)
		}
		seenSudoAdmins.Add(sudoAdmin)

		if _, err := sdk.AccAddressFromBech32(sudoAdmin); err != nil {
			return errorsmod.Wrapf(ErrInvalidGenesis, "invalid sudo admin address (%s)", err)
		}
	}

	return nil
}
