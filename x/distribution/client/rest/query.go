package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/x/distribution/client/common"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	// Get the rewards withdrawal address
	r.HandleFunc(
		"/distribution/delegators/{delegatorAddr}/withdraw_address",
		delegatorWithdrawalAddrHandlerFn(cliCtx, queryRoute),
	).Methods("GET")

	// Validator distribution information
	r.HandleFunc(
		"/distribution/validators/{validatorAddr}",
		validatorInfoHandlerFn(cliCtx, queryRoute),
	).Methods("GET")

	// Get the current distribution parameter values
	r.HandleFunc(
		"/distribution/parameters",
		paramsHandlerFn(cliCtx, queryRoute),
	).Methods("GET")

	// Get the amount held in the community pool
	r.HandleFunc(
		"/distribution/community_pool",
		communityPoolHandler(cliCtx, queryRoute),
	).Methods("GET")

}

// HTTP request handler to query a delegation rewards
func delegatorWithdrawalAddrHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		delegatorAddr, ok := checkDelegatorAddressVar(w, r)
		if !ok {
			return
		}

		cliCtx, ok = rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		bz := cliCtx.Codec.MustMarshalJSON(types.NewQueryDelegatorWithdrawAddrParams(delegatorAddr))
		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/withdraw_addr", queryRoute), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// ValidatorDistInfo defines the properties of
// validator distribution information response.
type ValidatorDistInfo struct {
	OperatorAddress     sdk.AccAddress                       `json:"operator_address" yaml:"operator_address"`
	SelfBondRewards     sdk.DecCoins                         `json:"self_bond_rewards" yaml:"self_bond_rewards"`
	ValidatorCommission types.ValidatorAccumulatedCommission `json:"val_commission" yaml:"val_commission"`
}

// NewValidatorDistInfo creates a new instance of ValidatorDistInfo.
func NewValidatorDistInfo(operatorAddr sdk.AccAddress,
	commission types.ValidatorAccumulatedCommission) ValidatorDistInfo {
	return ValidatorDistInfo{
		OperatorAddress:     operatorAddr,
		ValidatorCommission: commission,
	}
}

// HTTP request handler to query validator's distribution information
func validatorInfoHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validatorAddr, ok := checkValidatorAddressVar(w, r)
		if !ok {
			return
		}

		cliCtx, ok = rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// query commission
		commissionRes, err := common.QueryValidatorCommission(cliCtx, queryRoute, validatorAddr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var valCom types.ValidatorAccumulatedCommission
		cliCtx.Codec.MustUnmarshalJSON(commissionRes, &valCom)

		// self bond rewards
		delAddr := sdk.AccAddress(validatorAddr)

		// Prepare response
		res := cliCtx.Codec.MustMarshalJSON(NewValidatorDistInfo(delAddr, valCom))
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// HTTP request handler to query the distribution params values
func paramsHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params, err := common.QueryParams(cliCtx, queryRoute)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, params)
	}
}

func communityPoolHandler(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/community_pool", queryRoute), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var result sdk.DecCoins
		if err := cliCtx.Codec.UnmarshalJSON(res, &result); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, result)
	}
}
