package pow

import (
    "github.com/cosmos/cosmos-sdk/types"
)

type AppModuleBasic struct{}

func (AppModuleBasic) Name() string {
    return "pow"
}

func (AppModuleBasic) RegisterInterfaces(registry types.InterfaceRegistry) {}