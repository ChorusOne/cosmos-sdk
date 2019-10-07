package keeper

import (
	"fmt"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Return all validators that a delegator is bonded to. If maxRetrieve is supplied, the respective amount will be returned.
func (k Keeper) GetDelegatorValidators(ctx sdk.Context, delegatorAddr sdk.AccAddress,
	maxRetrieve uint16) (validators []types.Validator) {
	validators = []types.Validator{}

	// iterate over validators
	k.IterateValidators(ctx, func(index int64, validator exported.ValidatorI) (stop bool) {
		denom := fmt.Sprintf("%s%s", validator.GetSharesDenomPrefix(), k.BondDenom(ctx))
		coins := k.bankKeeper.GetCoins(ctx, delegatorAddr)
		if coins.AmountOf(denom).GT(sdk.ZeroInt()) {
			validators = append(validators, validator.(types.Validator))
		}
		return false
	})
	return validators[:int(math.Min(float64(len(validators)), float64(maxRetrieve)))] // trim if the array length < maxRetrieve
}

// return a validator that a delegator is bonded to
func (k Keeper) GetDelegatorValidator(ctx sdk.Context, delegatorAddr sdk.AccAddress,
	validatorAddr sdk.ValAddress) (validator types.Validator, err sdk.Error) {

	delegation, found := k.GetDelegation(ctx, delegatorAddr, validatorAddr)
	if !found {
		return validator, types.ErrNoDelegation(types.DefaultCodespace)
	}

	validator, found = k.GetValidator(ctx, delegation.ValidatorAddress)
	if !found {
		panic(types.ErrNoValidatorFound(types.DefaultCodespace))
	}
	return
}

//_____________________________________________________________________________________

// return all delegations for a delegator
func (k Keeper) GetAllDelegatorDelegations(ctx sdk.Context, delegator sdk.AccAddress) []types.Delegation {
	return k.GetDelegatorDelegations(ctx, delegator, math.MaxInt16)
}

// return all unbonding-delegations for a delegator
func (k Keeper) GetAllUnbondingDelegations(ctx sdk.Context, delegator sdk.AccAddress) []types.UnbondingDelegation {
	unbondingDelegations := make([]types.UnbondingDelegation, 0)

	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetUBDsKey(delegator)
	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey) // smallest to largest
	defer iterator.Close()

	for i := 0; iterator.Valid(); iterator.Next() {
		unbondingDelegation := types.MustUnmarshalUBD(k.cdc, iterator.Value())
		unbondingDelegations = append(unbondingDelegations, unbondingDelegation)
		i++
	}

	return unbondingDelegations
}

// return all redelegations for a delegator
func (k Keeper) GetAllRedelegations(ctx sdk.Context, delegator sdk.AccAddress,
	srcValAddress, dstValAddress sdk.ValAddress) (
	redelegations []types.Redelegation) {

	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetREDsKey(delegator)
	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey) // smallest to largest
	defer iterator.Close()

	srcValFilter := !(srcValAddress.Empty())
	dstValFilter := !(dstValAddress.Empty())

	for ; iterator.Valid(); iterator.Next() {
		redelegation := types.MustUnmarshalRED(k.cdc, iterator.Value())
		if srcValFilter && !(srcValAddress.Equals(redelegation.ValidatorSrcAddress)) {
			continue
		}
		if dstValFilter && !(dstValAddress.Equals(redelegation.ValidatorDstAddress)) {
			continue
		}
		redelegations = append(redelegations, redelegation)
	}
	return redelegations
}
