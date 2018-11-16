package connection

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mossid/ibc-mock/x/ibc/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgOpenConnection:
			return handleMsgOpenConnection(ctx, k, msg)
		default:
			panic("aaa") // TODO
		}
	}
}

func handleMsgOpenConnection(ctx sdk.Context, k Keeper, msg MsgOpenConnection) sdk.Result {
	id := msg.Config.UniqueID(msg.UserChainID)
	if k.existsChain(ctx, id) {
		return types.ErrConnectionAlreadyOpened(types.DefaultCodespace).Result()
	}
	k.registerChain(ctx, id, msg.Config)

	return sdk.Result{} // TODO: add tags
}

var (
	// TODO: move it to parameter store
	DIFFREQ = sdk.NewDec(5).Quo(sdk.NewDec(6)) // 2/3 + (1/3)/2
)

func BeginBlocker(ctx sdk.Context, k Keeper, header abci.Header) {
	if header.Height == 1 || !k.checkpointer.lastvalset.IsStillValid(ctx, DIFFREQ) {
		k.checkpointer.lastvalset.Snapshot(ctx)
		k.checkpointer.checkpoints.Set(ctx, uint64(header.Height), header)
	}
}
