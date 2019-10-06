package common

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// QueryParams actually queries distribution params.
func QueryParams(cliCtx context.CLIContext, queryRoute string) (PrettyParams, error) {
	route := fmt.Sprintf("custom/%s/params/%s", queryRoute, types.ParamCommunityTax)

	retCommunityTax, _, err := cliCtx.QueryWithData(route, []byte{})
	if err != nil {
		return PrettyParams{}, err
	}

	route = fmt.Sprintf("custom/%s/params/%s", queryRoute, types.ParamBaseProposerReward)
	retBaseProposerReward, _, err := cliCtx.QueryWithData(route, []byte{})
	if err != nil {
		return PrettyParams{}, err
	}

	route = fmt.Sprintf("custom/%s/params/%s", queryRoute, types.ParamBonusProposerReward)
	retBonusProposerReward, _, err := cliCtx.QueryWithData(route, []byte{})
	if err != nil {
		return PrettyParams{}, err
	}

	route = fmt.Sprintf("custom/%s/params/%s", queryRoute, types.ParamWithdrawAddrEnabled)
	retWithdrawAddrEnabled, _, err := cliCtx.QueryWithData(route, []byte{})
	if err != nil {
		return PrettyParams{}, err
	}

	return NewPrettyParams(
		retCommunityTax, retBaseProposerReward, retBonusProposerReward, retWithdrawAddrEnabled,
	), nil
}

// QueryDelegatorValidators returns delegator's list of validators
// it submitted delegations to.
func QueryDelegatorValidators(cliCtx context.CLIContext, queryRoute string, delegatorAddr sdk.AccAddress) ([]byte, error) {
	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryDelegatorValidators),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryDelegatorParams(delegatorAddr)),
	)
	return res, err
}

// QueryValidatorCommission returns a validator's commission.
func QueryValidatorCommission(cliCtx context.CLIContext, queryRoute string, validatorAddr sdk.ValAddress) ([]byte, error) {
	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryValidatorCommission),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryValidatorCommissionParams(validatorAddr)),
	)
	return res, err
}

// WithdrawValidatorRewardsAndCommission builds a two-message message slice to be
// used to withdraw both validation's commission and self-delegation reward.
func WithdrawValidatorCommission(validatorAddr sdk.ValAddress) ([]sdk.Msg, error) {
	commissionMsg := types.NewMsgWithdrawValidatorCommission(validatorAddr)
	if err := commissionMsg.ValidateBasic(); err != nil {
		return nil, err
	}

	return []sdk.Msg{commissionMsg}, nil
}
