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
	keeper := NewKeeper(mApp.Cdc, key, NewMockValidatorSet())

	mApp.Router().AddRoute("ibc", NewHandler(keeper))
	// TODO implement valset checkpointing
	// mApp.SetEndBlocker(getEndBlocker(mApp, keeper))
	mApp.SetInitChainer(getInitChainer(keeper))

	require.NoError(t, mApp.CompleteSetup(key))

	return mApp, keeper
}

func getInitChainer(k Keeper) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) (res abci.ResponseInitChain) {
		return
	}
}

type runtime struct {
	app1, app2 *mock.App
	k1, k2     Keeper
}

func newRuntime(t *testing.T) runtime {
	app1, k1 := getMockApp(t)
	app2, k2 := getMockApp(t)
	return runtime{
		app1: app1,
		app2: app2,
		k1:   k1,
		k2:   k2,
	}
}

func (r runtime) commit() {
	r.app1.Commit()
	r.app2.Commit()
}

func (r runtime) beginBlock() {
	r.app1.BeginBlock()
	r.app2.BeginBlock()
}

func (r runtime) openConnection(t *testing.T) {
	r.commit()
	r.beginBlock()
}

func (r runtime) establishConnection(t *testing.T) {
	// relay PayloadConnectionListening
}

func (r runtime) openChannel(t *testing.T) {
	// call MyChannel.Open()
}

func (r runtime) establishChannel(t *testing.T) {
	// relay PayloadChannelListening
}

func (r runtime) sendPackets(t *testing.T) {
	// send []MsgReceive, relay []MyPayload
}

func (r runtime) updateConnection(t *testing.T) {
	// send MsgUpdateConnection
}

func TestLifecycle(t *testing.T) {
	r := newRuntime(t)
	r.openConnection(t)
	r.establishConnection(t)
	r.openChannel(t)
	r.establishChannel(t)
	r.sendPackets(t)
	r.updateConnection(t)
	r.sendPackets(t)
}

func TestLongDuration(t *testing.T) {
	r := newRuntime(t)
	r.openConnection(t)
	r.establishConnection(t)
	r.openChannel(t)
	r.establishChannel(t)

	for i := 0; i < 1000; i++ {
		r.sendPackets(t)
		r.updateConnection(t)
	}
}
