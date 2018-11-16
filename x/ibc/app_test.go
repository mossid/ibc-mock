package ibc

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mock"
)

func getMockApp(t *testing.T) (*mock.App, Keeper) {
	mApp := mock.NewApp()

	RegisterCodec(mApp.Cdc)

	key := sdk.NewKVStoreKey("ibc")
	keeper := NewKeeper(mApp.Cdc, key)

	mApp.Router().AddRoute("ibc", NewHandler(keeper))
	// TODO implement valset checkpointing
	// mApp.SetEndBlocker(getEndBlocker(mApp, keeper))
	mApp.SetInitChainer(getInitChainer(mApp, keeper))

	require.NoError(t, mApp.CompleteSetup(key))

	return mApp, keeper
}

func getInitChainer(mapp *mock.App, k Keeper) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {

	}
}

func TestLifecycle(t *testing.T) {

}
