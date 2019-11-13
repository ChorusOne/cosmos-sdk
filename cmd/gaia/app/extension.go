package app

import (
	"database/sql"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/abci/types"
)

func (app *GaiaApp) BeginBlockHook(database *sql.DB, blockerFunctions []func(*GaiaApp, *sql.DB, sdk.Context, []sdk.ValAddress, []sdk.AccAddress, string, string, types.RequestBeginBlock, *sync.WaitGroup), vals []sdk.ValAddress, accs []sdk.AccAddress, network string, chainid string) sdk.BeginBlocker {
	return func(ctx sdk.Context, req types.RequestBeginBlock) types.ResponseBeginBlock {
		res := app.BeginBlocker(ctx, req)
		// fucntions
		wg := sync.WaitGroup{}
		for _, fn := range blockerFunctions {
			wg.Add(1)
			go fn(app, database, ctx, vals, accs, network, chainid, req, &wg)

		}
		wg.Wait()
		return res
	}
}
