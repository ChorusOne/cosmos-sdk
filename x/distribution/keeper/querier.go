package keeper

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryParams:
			return queryParams(ctx, path[1:], req, k)

		case types.QueryValidatorCommission:
			return queryValidatorCommission(ctx, path[1:], req, k)

		case types.QueryValidatorSlashes:
			return queryValidatorSlashes(ctx, path[1:], req, k)

		case types.QueryDelegatorValidators:
			return queryDelegatorValidators(ctx, path[1:], req, k)

		case types.QueryWithdrawAddr:
			return queryDelegatorWithdrawAddress(ctx, path[1:], req, k)

		case types.QueryCommunityPool:
			return queryCommunityPool(ctx, path[1:], req, k)

		default:
			return nil, sdk.ErrUnknownRequest("unknown distr query endpoint")
		}
	}
}

func queryParams(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	switch path[0] {
	case types.ParamCommunityTax:
		bz, err := codec.MarshalJSONIndent(k.cdc, k.GetCommunityTax(ctx))
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
		}
		return bz, nil
	case types.ParamBaseProposerReward:
		bz, err := codec.MarshalJSONIndent(k.cdc, k.GetBaseProposerReward(ctx))
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
		}
		return bz, nil
	case types.ParamBonusProposerReward:
		bz, err := codec.MarshalJSONIndent(k.cdc, k.GetBonusProposerReward(ctx))
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
		}
		return bz, nil
	case types.ParamWithdrawAddrEnabled:
		bz, err := codec.MarshalJSONIndent(k.cdc, k.GetWithdrawAddrEnabled(ctx))
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
		}
		return bz, nil
	default:
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("%s is not a valid query request path", req.Path))
	}
}

func queryValidatorCommission(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryValidatorCommissionParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}
	commission := k.GetValidatorAccumulatedCommission(ctx, params.ValidatorAddress)
	bz, err := codec.MarshalJSONIndent(k.cdc, commission)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryValidatorSlashes(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryValidatorSlashesParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}
	events := make([]types.ValidatorSlashEvent, 0)
	k.IterateValidatorSlashEventsBetween(ctx, params.ValidatorAddress, params.StartingHeight, params.EndingHeight,
		func(height uint64, event types.ValidatorSlashEvent) (stop bool) {
			events = append(events, event)
			return false
		},
	)
	bz, err := codec.MarshalJSONIndent(k.cdc, events)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryDelegatorValidators(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryDelegatorParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	// cache-wrap context as to not persist state changes during querying
	ctx, _ = ctx.CacheContext()

	var validators []sdk.ValAddress

	k.stakingKeeper.IterateDelegations(
		ctx, params.DelegatorAddress,
		func(_ int64, del exported.DelegationI) (stop bool) {
			validators = append(validators[:], del.GetValidatorAddr())
			return false
		},
	)

	bz, err := codec.MarshalJSONIndent(k.cdc, validators)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryDelegatorWithdrawAddress(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryDelegatorWithdrawAddrParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	// cache-wrap context as to not persist state changes during querying
	ctx, _ = ctx.CacheContext()
	withdrawAddr := k.GetDelegatorWithdrawAddr(ctx, params.DelegatorAddress)

	bz, err := codec.MarshalJSONIndent(k.cdc, withdrawAddr)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryCommunityPool(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	bz, err := k.cdc.MarshalJSON(k.GetFeePoolCommunityCoins(ctx))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
