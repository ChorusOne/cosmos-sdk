package rest

import (
	"time"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdkrest "github.com/cosmos/cosmos-sdk/types/rest"
	evidenceexported "github.com/cosmos/cosmos-sdk/x/evidence/exported"
	ibctmtypes "github.com/cosmos/cosmos-sdk/x/ibc/07-tendermint/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/99-wasm/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"net/http"
		"github.com/cosmos/cosmos-sdk/client/flags"
		abci "github.com/tendermint/tendermint/abci/types"
		ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
)

// REST client flags
const (
	RestClientID   = "client-id"
	RestRootHeight = "height"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	registerTxRoutes(cliCtx, r)
	r.HandleFunc(fmt.Sprintf("/ibc/wasm/client/{%s}", RestClientID), queryWasmClientStatesFn(cliCtx)).Methods("GET")
}

func queryWasmClientStatesFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		clientID := vars[RestClientID]
		prove := sdkrest.ParseQueryParamBool(r, flags.FlagProve)

		cliCtx, ok := sdkrest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.QueryClientStateParams{}
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			sdkrest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		req := abci.RequestQuery{
			Path:  "store/ibc/key",
			Data:  prefixClientKey(clientID, ibctypes.KeyClientState()),
			Prove: prove,
		}

		res, err := cliCtx.QueryABCI(req)
		if err != nil {
			sdkrest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var clientState types.ClientState
		if err := cliCtx.Codec.UnmarshalBinaryBare(res.Value, &clientState); err != nil {
			sdkrest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var address = clientState.ValidityPredicateAddress

		route := fmt.Sprintf("custom/%s/%s/%s/%s/smart", "ibc", types.QuerierRoute, types.QueryGetContractState, address)
		qres, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			sdkrest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(int64(height))
		sdkrest.PostProcessResponse(w, cliCtx, string(qres))
	}
}

// CreateClientReq defines the properties of a create client request's body.
type CreateClientReq struct {
	BaseReq         rest.BaseReq      `json:"base_req" yaml:"base_req"`
	ClientID        string            `json:"client_id" yaml:"client_id"`
	ChainID         string            `json:"chain_id" yaml:"chain_id"`
	Header          ibctmtypes.Header `json:"consensus_state" yaml:"consensus_state"`
	TrustingPeriod  time.Duration     `json:"trusting_period" yaml:"trusting_period"`
	MaxClockDrift   time.Duration     `json:"max_clock_drift" yaml:"max_clock_drift"`
	UnbondingPeriod time.Duration     `json:"unbonding_period" yaml:"unbonding_period"`
	WasmId          int 			  `json:"wasm_id" yaml:"wasm_id"`
}

// UpdateClientReq defines the properties of a update client request's body.
type UpdateClientReq struct {
	BaseReq rest.BaseReq      `json:"base_req" yaml:"base_req"`
	Header  ibctmtypes.Header `json:"header" yaml:"header"`
}

// SubmitMisbehaviourReq defines the properties of a submit misbehaviour request's body.
type SubmitMisbehaviourReq struct {
	BaseReq  rest.BaseReq              `json:"base_req" yaml:"base_req"`
	Evidence evidenceexported.Evidence `json:"evidence" yaml:"evidence"`
}

func prefixClientKey(clientID string, key []byte) []byte {
	return append([]byte(fmt.Sprintf("clients/%s/", clientID)), key...)
}
