package types

import (
	"github.com/CosmWasm/wasmvm/api"
	"github.com/CosmWasm/wasmvm/types"
	ics23 "github.com/confio/ics23/go"
	sdk "github.com/cosmos/cosmos-sdk/types"
	types2 "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	types3 "github.com/cosmos/cosmos-sdk/x/ibc/core/23-commitment/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/core/28-wasm/keeper"
	"github.com/cosmos/cosmos-sdk/x/ibc/core/exported"
)

const GasMultiplier uint64 = 100

var _ exported.ClientState = (*ClientState)(nil)

type queryResponse struct {
	ProofSpecs []*ics23.ProofSpec `json:"proof_specs,omitempty"`
	Height types2.Height `json:"height,omitempty"`
	IsFrozen bool `json:"is_frozen,omitempty"`
	FrozenHeight types2.Height `json:"frozen_height,omitempty"`
	GenesisMetadata []types2.GenesisMetadata `json:"genesis_metadata,omitempty"`
	Result contractResult `json:"result,omitempty"`
	ClientType string `json:"client_type,omitempty"`
	Root types3.MerkleRoot `json:"root,omitempty"`
	Timestamp uint64 `json:"timestamp,omitempty"`
}

type contractResult struct {
	IsValid bool `json:"is_valid,omitempty"`
	ErrorMsg string `json:"err_msg,omitempty"`
}

type clientStateCallResponse struct {
	Self              *ClientState    `json:"self,omitempty"`
	NewConsensusState *ConsensusState `json:"consensus_state,omitempty"`
	NewClientState    *ClientState    `json:"client_state,omitempty"`
	Result contractResult `json:"result,omitempty"`
}

func callContract(codeId []byte, ctx sdk.Context, store sdk.KVStore, msg []byte) (*types.HandleResponse, error) {
	gasMeter := ctx.BlockGasMeter()
	chainID := ctx.BlockHeader().ChainID
	height := ctx.BlockHeader().Height
	// safety checks before casting below
	if height < 0 {
		panic("Block height must never be negative")
	}
	sec := ctx.BlockTime().Unix()
	if sec < 0 {
		panic("Block (unix) time must never be negative ")
	}
	nano := ctx.BlockTime().Nanosecond()
	env := types.Env{
		Block: types.BlockInfo{
			Height:    uint64(height),
			Time:      uint64(sec),
			TimeNanos: uint64(nano),
			ChainID:   chainID,
		},
		Contract: types.ContractInfo{
			Address: "",
		},
	}

	return callContractWithEnvAndMeter(codeId, &ctx, store, env, gasMeter, msg)
}

func callContractWithEnvAndMeter(codeId []byte, ctx *sdk.Context, store sdk.KVStore, env types.Env, gasMeter sdk.GasMeter, msg []byte) (*types.HandleResponse, error) {
	msgInfo := types.MessageInfo{
		Sender:  "",
		SentFunds: nil,
	}
	mockFailureAPI := *api.NewMockFailureAPI()
	mockQuerier := api.MockQuerier{}

	resp, gasUsed, err := keeper.WasmVM.Execute(codeId, env, msgInfo, msg, store, mockFailureAPI, mockQuerier, gasMeter, gasMeter.Limit())
	if ctx != nil {
		consumeGas(*ctx, gasUsed)
	}
	return resp, err
}

func queryContract(codeId []byte, msg []byte) ([]byte, error) {
	return queryContractWithStore(codeId, &FailKVStore{}, msg)
}

func queryContractWithStore(codeId []byte, store sdk.KVStore, msg []byte) ([]byte, error) {
	mockEnv := api.MockEnv()
	mockGasMeter := api.NewMockGasMeter(1)
	mockFailureAPI := *api.NewMockFailureAPI()
	mockQuerier := api.MockQuerier{}

	resp, _, err := keeper.WasmVM.Query(codeId, mockEnv, msg, store, mockFailureAPI, mockQuerier, mockGasMeter, 0)
	return resp, err
}

func consumeGas(ctx sdk.Context, gas uint64) {
	consumed := gas / GasMultiplier
	ctx.GasMeter().ConsumeGas(consumed, "wasm contract")
	// throw OutOfGas error if we ran out (got exactly to zero due to better limit enforcing)
	if ctx.GasMeter().IsOutOfGas() {
		panic(sdk.ErrorOutOfGas{Descriptor: "Wasmer function execution"})
	}
}

