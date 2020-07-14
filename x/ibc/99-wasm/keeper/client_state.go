package keeper

import (
	"errors"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ibcwasmtypes "github.com/cosmos/cosmos-sdk/x/ibc/99-wasm/types"
)

// InitializeFromMsg creates a tendermint client state from a CreateClientMsg
func (k *Keeper) InitializeFromMsg(ctx sdk.Context,
	msg ibcwasmtypes.MsgCreateWasmClient,
) (ibcwasmtypes.ClientState, error) {
	return k.Initialize(ctx, msg.GetClientID(), msg.TrustingPeriod, msg.UnbondingPeriod, msg.MaxClockDrift, msg, msg.WasmId)
}

// Initialize creates a client state and validates its contents, checking that
// the provided consensus state is from the same client type.
func (k *Keeper) Initialize(
	ctx sdk.Context,
	id string, trustingPeriod, ubdPeriod, maxClockDrift time.Duration,
	initMsg ibcwasmtypes.MsgCreateWasmClient, wasmId int,
) (ibcwasmtypes.ClientState, error) {
	if trustingPeriod >= ubdPeriod {
		return ibcwasmtypes.ClientState{}, errors.New("trusting period should be < unbonding period")
	}
	contractAddress, err := k.Instantiate(ctx,  initMsg.ClientID, uint64(wasmId), ibcwasmtypes.ModuleAccount, initMsg.Header, fmt.Sprintf("wasm-client-%s-%d", id, wasmId))
	if err != nil {
		return ibcwasmtypes.ClientState{}, err
	}
	clientState := ibcwasmtypes.NewClientState(
		id, trustingPeriod, ubdPeriod, maxClockDrift, ibcwasmtypes.Header{initMsg.Header}, contractAddress,
	)
	return clientState, nil
}
