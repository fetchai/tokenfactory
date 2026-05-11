package types

import (
	"sync"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

var (
	ModuleAddress = sync.OnceValue(func() string {
		return authtypes.NewModuleAddress(ModuleName).String()
	})
)
