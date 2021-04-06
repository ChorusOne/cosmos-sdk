package types

import (
	"fmt"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	// DefaultSendEnabled enabled
	DefaultWASMClientEnabled = false
)

var (
	KeyWasmClientEnabled = []byte("WasmClientEnabled")
)

// ParamKeyTable type declaration for parameters
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new parameter configuration for the ibc transfer module
func NewParams(enableWASMClient bool) Params {
	return Params{
		WasmClientEnabled: enableWASMClient,
	}
}

// DefaultParams is the default parameter configuration for the ibc-transfer module
func DefaultParams() Params {
	return NewParams(DefaultWASMClientEnabled)
}

// Validate all ibc-transfer module parameters
func (p Params) Validate() error {
	return validateWASMClientEnabled(p.WasmClientEnabled)
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyWasmClientEnabled, p.WasmClientEnabled, validateWASMClientEnabled),
	}
}

func validateWASMClientEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

