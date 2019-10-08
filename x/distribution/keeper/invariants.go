package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"
)

// register all distribution invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "can-withdraw",
		CanWithdrawInvariant(k))
	ir.RegisterRoute(types.ModuleName, "reference-count",
		ReferenceCountInvariant(k))
}

// AllInvariants runs all invariants of the distribution module
func AllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		res, stop := CanWithdrawInvariant(k)(ctx)
		if stop {
			return res, stop
		}
		return ReferenceCountInvariant(k)(ctx)
	}
}

// CanWithdrawInvariant checks that current rewards can be completely withdrawn
func CanWithdrawInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {

		// cache, we don't want to write changes
		ctx, _ = ctx.CacheContext()

		var remaining sdk.DecCoins

		valDelegationAddrs := make(map[string][]sdk.AccAddress)
		for _, del := range k.stakingKeeper.GetAllSDKDelegations(ctx) {
			valAddr := del.GetValidatorAddr().String()
			valDelegationAddrs[valAddr] = append(valDelegationAddrs[valAddr], del.GetDelegatorAddr())
		}

		// iterate over all validators
		k.stakingKeeper.IterateValidators(ctx, func(_ int64, val exported.ValidatorI) (stop bool) {
			_, _ = k.WithdrawValidatorCommission(ctx, val.GetOperator())

			//delegationAddrs, ok := valDelegationAddrs[val.GetOperator().String()]
			// if ok {
			// 	for _, delAddr := range delegationAddrs {
			// 		if _, err := k.WithdrawDelegationRewards(ctx, delAddr, val.GetOperator()); err != nil {
			// 			panic(err)
			// 		}
			// 	}
			// }

			return false
		})

		broken := len(remaining) > 0 && remaining[0].Amount.LT(sdk.ZeroDec())
		return sdk.FormatInvariant(types.ModuleName, "can withdraw",
			fmt.Sprintf("remaining coins: %v\n", remaining)), broken
	}
}

// ReferenceCountInvariant checks that the number of historical rewards records is correct
func ReferenceCountInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {

		valCount := uint64(0)
		k.stakingKeeper.IterateValidators(ctx, func(_ int64, val exported.ValidatorI) (stop bool) {
			valCount++
			return false
		})
		dels := k.stakingKeeper.GetAllSDKDelegations(ctx)
		slashCount := uint64(0)
		k.IterateValidatorSlashEvents(ctx,
			func(_ sdk.ValAddress, _ uint64, _ types.ValidatorSlashEvent) (stop bool) {
				slashCount++
				return false
			})

		// one record per validator (last tracked period), one record per
		// delegation (previous period), one record per slash (previous period)
		expected := valCount + uint64(len(dels)) + slashCount
		count := k.GetValidatorHistoricalReferenceCount(ctx)
		broken := count != expected

		return sdk.FormatInvariant(types.ModuleName, "reference count",
			fmt.Sprintf("expected historical reference count: %d = %v validators + %v delegations + %v slashes\n"+
				"total validator historical reference count: %d\n",
				expected, valCount, len(dels), slashCount, count)), broken
	}
}
